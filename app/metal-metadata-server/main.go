// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/ghodss/yaml"
	"github.com/talos-systems/talos/pkg/machinery/config/configloader"
	"github.com/talos-systems/talos/pkg/machinery/config/types/v1alpha1"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/cluster-api/util"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/talos-systems/sidero/app/cluster-api-provider-sidero/api/v1alpha3"
	metalv1alpha1 "github.com/talos-systems/sidero/app/metal-controller-manager/api/v1alpha1"
	"github.com/talos-systems/sidero/app/metal-metadata-server/pkg/client"
)

var (
	kubeconfigPath *string
	port           *string
)

type errorWithCode struct {
	errorCode int
	errorObj  error
}

type metadataConfigs struct {
	client runtimeclient.Client
}

func throwError(w http.ResponseWriter, ewc errorWithCode) {
	http.Error(w, ewc.errorObj.Error(), ewc.errorCode)
	log.Println(ewc.errorObj)
}

func main() {
	kubeconfigPath = flag.String("kubeconfig-path", "", "absolute path to the kubeconfig file")
	port = flag.String("port", "8080", "port to use for serving metadata")
	flag.Parse()

	k8sClient, err := client.NewClient(kubeconfigPath)
	if err != nil {
		log.Fatal(fmt.Errorf("failure talking to kubernetes: %s", err))
	}

	mm := metadataConfigs{
		client: k8sClient,
	}

	http.HandleFunc("/configdata", mm.FetchConfig)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

func (m *metadataConfigs) FetchConfig(w http.ResponseWriter, r *http.Request) {
	// Parse info out of incoming request
	ctx := r.Context()

	vals := r.URL.Query()

	uuid := vals.Get("uuid")

	// Throw out requests with no uuid param
	if len(uuid) == 0 {
		throwError(
			w,
			errorWithCode{
				http.StatusInternalServerError,
				fmt.Errorf(
					"received metadata request with empty uuid",
				),
			},
		)

		return
	}

	log.Printf("received metadata request for uuid: %s", uuid)

	// Find all MetalMachine resources that have this UUID as a ServerRef in their spec.
	// It's an error if 0 or more than 1 result is found.
	matchingMetalMachine, ewc := m.matchingMetalMachine(ctx, uuid)
	if ewc.errorObj != nil {
		throwError(
			w,
			ewc,
		)

		return
	}

	// Given the MetalMachine, find the Machine resource that owns it
	ownerMachine, err := util.GetOwnerMachine(ctx, m.client, matchingMetalMachine.ObjectMeta)
	if err != nil {
		throwError(
			w,
			errorWithCode{
				http.StatusInternalServerError,
				fmt.Errorf(
					"failure fetching owner machine from metal machine %s/%s: %s",
					matchingMetalMachine.GetNamespace(),
					matchingMetalMachine.GetName(),
					err,
				),
			},
		)

		return
	}

	// Dig bootstrap secret name out of owner Machine resource and fetch secret data
	bootstrapSecretName := ownerMachine.Spec.Bootstrap.DataSecretName

	if bootstrapSecretName == nil {
		throwError(
			w,
			errorWithCode{
				http.StatusNotFound,
				fmt.Errorf(
					"no dataSecretName present for machine %s/%s",
					ownerMachine.Namespace,
					ownerMachine.Name,
				),
			},
		)

		return
	}

	decodedData, ewc := m.fetchBootstrapSecret(
		ctx,
		types.NamespacedName{
			Name:      *bootstrapSecretName,
			Namespace: ownerMachine.Namespace,
		},
	)
	if ewc.errorObj != nil {
		throwError(
			w,
			ewc,
		)

		return
	}

	// Get the server resource by the UUID that was passed in.
	// We do this to fetch serverclass and any configPatches in the server resource that we need to handle.
	serverObj := &metalv1alpha1.Server{}

	err = m.client.Get(
		ctx,
		types.NamespacedName{
			Namespace: "",
			Name:      uuid,
		},
		serverObj,
	)
	if err != nil {
		throwError(
			w,
			errorWithCode{
				http.StatusInternalServerError,
				fmt.Errorf(
					"failure fetching server %s: %s",
					uuid,
					err,
				),
			},
		)

		return
	}

	// Given a server object, see if it came from a serverclass (it will have an ownerref)
	// If so, fetch the serverclass so we can use configPatches from it.
	serverClassObj := &metalv1alpha1.ServerClass{}

	if len(serverObj.OwnerReferences) > 0 {
		for _, ownerRef := range serverObj.OwnerReferences {
			if ownerRef.Kind == "ServerClass" {
				err = m.client.Get(
					ctx,
					types.NamespacedName{
						Namespace: "",
						Name:      ownerRef.Name,
					},
					serverClassObj,
				)
				if err != nil {
					throwError(
						w,
						errorWithCode{
							http.StatusInternalServerError,
							fmt.Errorf(
								"failure fetching serverclass  %s: %s",
								ownerRef.Name,
								err,
							),
						},
					)

					return
				}
			}
		}
	}

	// Handle patches added to serverclass object
	if serverClassObj != nil && len(serverClassObj.Spec.ConfigPatches) > 0 {
		decodedData, ewc = patchConfigs(decodedData, serverClassObj.Spec.ConfigPatches)
		if ewc.errorObj != nil {
			throwError(
				w,
				ewc,
			)

			return
		}
	}

	// Handle patches added to server object
	if len(serverObj.Spec.ConfigPatches) > 0 {
		decodedData, ewc = patchConfigs(decodedData, serverObj.Spec.ConfigPatches)
		if ewc.errorObj != nil {
			throwError(
				w,
				ewc,
			)

			return
		}
	}

	// Append or add a node label to kubelet extra args.
	// We must do this so that we can map a given server resource to a k8s node in the workload cluster.
	decodedData, ewc = labelNodes(decodedData, serverObj.Name)
	if ewc.errorObj != nil {
		throwError(
			w,
			ewc,
		)

		return
	}

	// Finally return config data
	if _, err = w.Write(decodedData); err != nil {
		log.Printf("failed to write data: %v", err)
		return
	}

	log.Printf("successfully returned metadata for %q", uuid)
}

// patchConfigs is responsible for applying a set of configPatches to the bootstrap data.
func patchConfigs(decodedData []byte, patches []metalv1alpha1.ConfigPatches) ([]byte, errorWithCode) {
	marshalledPatches, err := json.Marshal(patches)
	if err != nil {
		return nil, errorWithCode{http.StatusInternalServerError, fmt.Errorf("failure marshalling config patches from server: %s", err)}
	}

	jsonDecodedData, err := yaml.YAMLToJSON(decodedData)
	if err != nil {
		return nil, errorWithCode{http.StatusInternalServerError, fmt.Errorf("failure converting bootstrap data to json: %s", err)}
	}

	patch, err := jsonpatch.DecodePatch(marshalledPatches)
	if err != nil {
		return nil, errorWithCode{http.StatusInternalServerError, fmt.Errorf("failure decoding config patches from server to rfc6902 patch: %s", err)}
	}

	jsonDecodedData, err = patch.Apply(jsonDecodedData)
	if err != nil {
		return nil, errorWithCode{http.StatusInternalServerError, fmt.Errorf("failure applying rfc6902 patches to machine config: %s", err)}
	}

	decodedData, err = yaml.JSONToYAML(jsonDecodedData)
	if err != nil {
		return nil, errorWithCode{http.StatusInternalServerError, fmt.Errorf("failure converting bootstrap data from json to yaml: %s", err)}
	}

	return decodedData, errorWithCode{}
}

// labelNodes is responsible for editing the kubelet extra args such that a given
// server gets registered with a label containing the UUID of the server resource it's actually running on.
func labelNodes(decodedData []byte, serverName string) ([]byte, errorWithCode) {
	configProvider, err := configloader.NewFromBytes(decodedData)
	if err != nil {
		return nil, errorWithCode{http.StatusInternalServerError, fmt.Errorf("failure creating config struct: %s", err)}
	}

	switch configProvider.Version() {
	case "v1alpha1":
		config, ok := configProvider.(*v1alpha1.Config)
		if !ok {
			return nil, errorWithCode{http.StatusInternalServerError, fmt.Errorf("unable to case config")}
		}

		if _, ok = config.MachineConfig.MachineKubelet.KubeletExtraArgs["node-labels"]; ok {
			config.MachineConfig.MachineKubelet.KubeletExtraArgs["node-labels"] += fmt.Sprintf(",metal.sidero.dev/uuid=%s", serverName)
		} else {
			if config.MachineConfig.MachineKubelet.KubeletExtraArgs == nil {
				config.MachineConfig.MachineKubelet.KubeletExtraArgs = make(map[string]string)
			}
			config.MachineConfig.MachineKubelet.KubeletExtraArgs["node-labels"] = fmt.Sprintf("metal.sidero.dev/uuid=%s", serverName)
		}

		decodedData, err = config.Bytes()
		if err != nil {
			return nil, errorWithCode{http.StatusInternalServerError, fmt.Errorf("failure converting config to bytes: %s", err)}
		}
	default:
		return nil, errorWithCode{http.StatusInternalServerError, fmt.Errorf("unknown config type")}
	}

	return decodedData, errorWithCode{}
}

// matchingMetalMachine is responsible for looking up MetalMachines that contain a ServerRef with the UUID that requested metadata.
func (m *metadataConfigs) matchingMetalMachine(ctx context.Context, serverName string) (v1alpha3.MetalMachine, errorWithCode) {
	metalMachineList := &v1alpha3.MetalMachineList{}

	err := m.client.List(ctx, metalMachineList)
	if err != nil {
		return v1alpha3.MetalMachine{}, errorWithCode{http.StatusInternalServerError, fmt.Errorf("failure listing metal machines: %w", err)}
	}

	numMatchingMachines := 0

	var matchingMetalMachine v1alpha3.MetalMachine

	// range through all metalmachines and find all that match our server string
	for _, metalMachine := range metalMachineList.Items {
		// handle metalmachines where serverref hasn't been updated yet
		if metalMachine.Spec.ServerRef == nil {
			continue
		}

		if metalMachine.Spec.ServerRef.Name == serverName {
			numMatchingMachines++

			matchingMetalMachine = metalMachine
		}
	}

	// bail early if we have multiple matches or no matches
	switch {
	case numMatchingMachines > 1:
		return v1alpha3.MetalMachine{}, errorWithCode{http.StatusInternalServerError, fmt.Errorf("multiple matching metal machines for uuid %q, possible orphaned cluster", serverName)}

	case numMatchingMachines == 0:
		return v1alpha3.MetalMachine{}, errorWithCode{http.StatusInternalServerError, fmt.Errorf("failure finding matching metal machine for uuid: %q", serverName)}
	}

	return matchingMetalMachine, errorWithCode{}
}

// fetchBootstrapSecret is responsible for fetching a secret that contains the bootstrap data created by our bootstrap provider.
func (m *metadataConfigs) fetchBootstrapSecret(ctx context.Context, secretNSN types.NamespacedName) ([]byte, errorWithCode) {
	bootstrapSecretData := &v1.Secret{}

	err := m.client.Get(
		ctx,
		secretNSN,
		bootstrapSecretData,
	)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil, errorWithCode{http.StatusNotFound, fmt.Errorf("bootstrap secret %s/%s not found", secretNSN.Namespace, secretNSN.Name)}
		}

		return nil, errorWithCode{http.StatusInternalServerError, fmt.Errorf("failure fetching bootstrap secret data from secret %s/%s: %s", secretNSN.Namespace, secretNSN.Name, err)}
	}

	if _, ok := bootstrapSecretData.Data["value"]; !ok {
		return nil, errorWithCode{http.StatusNotFound, fmt.Errorf("value key not found in bootstrap data: %s/%s", secretNSN.Namespace, secretNSN.Name)}
	}

	return bootstrapSecretData.Data["value"], errorWithCode{}
}
