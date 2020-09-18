// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/ghodss/yaml"
	"github.com/talos-systems/talos/pkg/config"
	"github.com/talos-systems/talos/pkg/config/types/v1alpha1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	metalv1alpha1 "github.com/talos-systems/sidero/app/metal-controller-manager/api/v1alpha1"
	"github.com/talos-systems/sidero/app/metal-metadata-server/pkg/client"
)

var (
	kubeconfig *string
	port       *string
)

const (
	capiVersion  = "v1alpha3"
	metalVersion = "v1alpha1"
)

func throwError(w http.ResponseWriter, code int, err error) {
	http.Error(w, err.Error(), code)
	log.Println(err)
}

func main() {
	kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	port = flag.String("port", "8080", "port to use for serving metadata")
	flag.Parse()

	http.HandleFunc("/configdata", FetchConfig)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

func FetchConfig(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vals := r.URL.Query()

	uuid := vals.Get("uuid")

	if len(uuid) == 0 {
		throwError(
			w,
			http.StatusInternalServerError,
			fmt.Errorf(
				"received metadata request with empty uuid",
			),
		)
	}

	log.Printf("received metadata request for uuid: %s", uuid)

	k8sClient, err := client.NewClient(kubeconfig)
	if err != nil {
		throwError(
			w,
			http.StatusInternalServerError,
			fmt.Errorf(
				"failure talking to kubernetes: %s",
				err,
			),
		)

		return
	}

	metalMachineGVR := schema.GroupVersionResource{
		Group:    "infrastructure.cluster.x-k8s.io",
		Version:  capiVersion,
		Resource: "metalmachines",
	}

	capiMachineGVR := schema.GroupVersionResource{
		Group:    "cluster.x-k8s.io",
		Version:  capiVersion,
		Resource: "machines",
	}

	secretGVR := schema.GroupVersionResource{
		Version:  "v1",
		Resource: "secrets",
	}

	serverGVR := schema.GroupVersionResource{
		Group:    "metal.sidero.dev",
		Version:  metalVersion,
		Resource: "servers",
	}

	metalMachineList, err := k8sClient.Resource(metalMachineGVR).Namespace("").List(ctx, metav1.ListOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			throwError(
				w,
				http.StatusNotFound,
				fmt.Errorf(
					"failure listing metal machines",
				),
			)

			return
		}

		throwError(w, http.StatusInternalServerError, fmt.Errorf("failure listing metal machines: %w", err))

		return
	}

	// Range through all metalMachines, seeing if we can match inventory by UUID
	for _, metalMachine := range metalMachineList.Items {
		serverRefString, _, err := unstructured.NestedString(metalMachine.Object, "spec", "serverRef", "name")
		if err != nil {
			throwError(w,
				http.StatusInternalServerError,
				fmt.Errorf(
					"failure finding serverRef name from metal machine %s/%s: %s",
					metalMachine.GetNamespace(),
					metalMachine.GetName(),
					err,
				))

			return
		}

		// Check if server ref isn't set. Assuming this is an unstructured thing where it's not an error, just empty.
		if serverRefString == "" {
			continue
		}

		// If ref matches, fetch the bootstrap data from machine resource that owns this metal machine
		if serverRefString == uuid {
			ownerList, present, err := unstructured.NestedSlice(metalMachine.Object, "metadata", "ownerReferences")
			if err != nil {
				throwError(
					w,
					http.StatusInternalServerError,
					fmt.Errorf(
						"failure fetching ownerRefs from metal machine %s/%s: %s",
						metalMachine.GetNamespace(),
						metalMachine.GetName(),
						err,
					),
				)

				return
			}

			if !present {
				throwError(
					w,
					http.StatusNotFound,
					fmt.Errorf(
						"ownerRefList not found for metalMachine",
					),
				)

				return
			}

			var ownerRef *metav1.OwnerReference

			for _, ownerItem := range ownerList {
				tempOwnerRef := &metav1.OwnerReference{}

				err = runtime.DefaultUnstructuredConverter.FromUnstructured(ownerItem.(map[string]interface{}), tempOwnerRef)
				if err != nil {
					throwError(
						w,
						http.StatusInternalServerError,
						fmt.Errorf("failure converting ownerItem to ownerRef: %s", err),
					)

					return
				}

				if tempOwnerRef.APIVersion == "cluster.x-k8s.io/"+capiVersion && tempOwnerRef.Kind == "Machine" {
					ownerRef = tempOwnerRef
					break
				}
			}

			if ownerRef == nil {
				throwError(
					w,
					http.StatusNotFound,
					fmt.Errorf(
						"no ownerrefs for metal machine %s/%s",
						metalMachine.GetNamespace(),
						metalMachine.GetName(),
					),
				)

				return
			}

			metalMachineNS, present, err := unstructured.NestedString(metalMachine.Object, "metadata", "namespace")
			if err != nil {
				throwError(
					w,
					http.StatusInternalServerError,
					fmt.Errorf(
						"failure fetching namespace for metal machine %s/%s: %s",
						metalMachine.GetNamespace(),
						metalMachine.GetName(),
						err,
					),
				)

				return
			}

			if !present {
				throwError(
					w,
					http.StatusNotFound,
					fmt.Errorf(
						"no namespace present for metal machine %s/%s",
						metalMachine.GetNamespace(),
						metalMachine.GetName(),
					),
				)

				return
			}

			machineData, err := k8sClient.Resource(capiMachineGVR).Namespace(metalMachineNS).Get(ctx, ownerRef.Name, metav1.GetOptions{})
			if err != nil {
				if apierrors.IsNotFound(err) {
					throwError(
						w,
						http.StatusNotFound,
						fmt.Errorf(
							"owner machine %s/%s not found",
							metalMachineNS,
							ownerRef.Name,
						),
					)

					return
				}

				throwError(
					w,
					http.StatusInternalServerError,
					fmt.Errorf(
						"failure fetching machine based on owner ref %s: %s",
						ownerRef.Name,
						err,
					),
				)

				return
			}

			bootstrapSecretName, present, err := unstructured.NestedString(machineData.Object, "spec", "bootstrap", "dataSecretName")
			if err != nil {
				throwError(
					w,
					http.StatusInternalServerError,
					fmt.Errorf(
						"failure fetching bootstrap dataSecretName from machine %s/%s: %s",
						machineData.GetNamespace(),
						machineData.GetName(),
						err,
					),
				)

				return
			}

			if !present {
				throwError(
					w,
					http.StatusNotFound,
					fmt.Errorf(
						"no dataSecretName present for machine %s/%s: %s",
						machineData.GetNamespace(),
						machineData.GetName(),
						err,
					),
				)

				return
			}

			bootstrapSecretData, err := k8sClient.Resource(secretGVR).Namespace(metalMachineNS).Get(
				ctx,
				bootstrapSecretName,
				metav1.GetOptions{},
			)
			if err != nil {
				if apierrors.IsNotFound(err) {
					throwError(
						w,
						http.StatusNotFound,
						fmt.Errorf(
							"bootstrap secret %s/%s not found",
							metalMachineNS,
							bootstrapSecretName,
						),
					)

					return
				}

				throwError(
					w,
					http.StatusInternalServerError,
					fmt.Errorf(
						"failure fetching bootstrap secret data from secret %s/%s: %s",
						metalMachineNS,
						bootstrapSecretName,
						err,
					),
				)

				return
			}

			bootstrapData, present, err := unstructured.NestedString(bootstrapSecretData.Object, "data", "value")
			if err != nil {
				throwError(
					w,
					http.StatusInternalServerError,
					fmt.Errorf(
						"failure fetching value key from bootstrap secret %s/%s: %s",
						bootstrapSecretData.GetName(),
						bootstrapSecretData.GetNamespace(),
						err,
					),
				)

				return
			}

			if !present {
				throwError(
					w,
					http.StatusNotFound,
					fmt.Errorf(
						"no bootstrap data found in value key of secret %s/%s",
						bootstrapSecretData.GetName(),
						bootstrapSecretData.GetNamespace(),
					),
				)

				return
			}

			decodedData, err := base64.StdEncoding.DecodeString(bootstrapData)
			if err != nil {
				throwError(
					w,
					http.StatusInternalServerError,
					fmt.Errorf(
						"failure decoding base64 from bootstrap data: %s",
						err,
					),
				)

				return
			}

			// Convert server uuid to unstructured obj and then to structured obj.
			serverRef, err := k8sClient.Resource(serverGVR).Get(ctx, serverRefString, metav1.GetOptions{})
			if err != nil {
				throwError(
					w,
					http.StatusInternalServerError,
					fmt.Errorf(
						"failure fetching server %s: %s",
						serverRefString,
						err,
					),
				)

				return
			}

			serverObj := &metalv1alpha1.Server{}

			err = runtime.DefaultUnstructuredConverter.FromUnstructured(serverRef.UnstructuredContent(), serverObj)
			if err != nil {
				throwError(
					w,
					http.StatusInternalServerError,
					fmt.Errorf(
						"failure converting server to metalv1alpha1.Server type %s: %s",
						serverRefString,
						err,
					),
				)

				return
			}

			// Handle patches added to server object
			if len(serverObj.Spec.ConfigPatches) > 0 {
				marshalledPatches, err := json.Marshal(serverObj.Spec.ConfigPatches)
				if err != nil {
					throwError(
						w,
						http.StatusInternalServerError,
						fmt.Errorf(
							"failure marshalling config patches from server %s: %s",
							serverObj.Name,
							err,
						),
					)

					return
				}

				jsonDecodedData, err := yaml.YAMLToJSON(decodedData)
				if err != nil {
					throwError(
						w,
						http.StatusInternalServerError,
						fmt.Errorf(
							"failure converting bootstrap data to json: %s",
							err,
						),
					)

					return
				}

				patch, err := jsonpatch.DecodePatch(marshalledPatches)
				if err != nil {
					throwError(
						w,
						http.StatusInternalServerError,
						fmt.Errorf(
							"failure decoding config patches from server %s to rfc6902 patch: %s",
							serverObj.Name,
							err,
						),
					)

					return
				}

				jsonDecodedData, err = patch.Apply(jsonDecodedData)
				if err != nil {
					throwError(
						w,
						http.StatusInternalServerError,
						fmt.Errorf(
							"failure applying rfc6902 patches to machine config: %s",
							err,
						),
					)

					return
				}

				decodedData, err = yaml.JSONToYAML(jsonDecodedData)
				if err != nil {
					throwError(
						w,
						http.StatusInternalServerError,
						fmt.Errorf(
							"failure converting bootstrap data from json to yaml: %s",
							err,
						),
					)

					return
				}
			}

			// Append or add a node label to kubelet extra args
			configStruct, err := config.NewFromBytes(decodedData)
			if err != nil {
				throwError(
					w,
					http.StatusInternalServerError,
					fmt.Errorf(
						"failure creating config struct: %s",
						err,
					),
				)

				return
			}

			// nolint: gocritic
			switch config := configStruct.(type) {
			case *v1alpha1.Config:
				if _, ok := config.MachineConfig.MachineKubelet.KubeletExtraArgs["node-labels"]; ok {
					config.MachineConfig.MachineKubelet.KubeletExtraArgs["node-labels"] += fmt.Sprintf(",metal.sidero.dev/uuid=%s", serverObj.Name)
				} else {
					if config.MachineConfig.MachineKubelet.KubeletExtraArgs == nil {
						config.MachineConfig.MachineKubelet.KubeletExtraArgs = make(map[string]string)
					}
					config.MachineConfig.MachineKubelet.KubeletExtraArgs["node-labels"] = fmt.Sprintf("metal.sidero.dev/uuid=%s", serverObj.Name)
				}

				decodedData, err = config.Bytes()
				if err != nil {
					throwError(
						w,
						http.StatusInternalServerError,
						fmt.Errorf(
							"failure converting config to bytes: %s",
							err,
						),
					)

					return
				}
			}

			// Finally return config data
			if _, err = w.Write(decodedData); err != nil {
				log.Printf("Failed to write data: %v", err)
			}

			return
		}
	}

	// Made it through all metal machines w/ no result
	throwError(
		w,
		http.StatusNotFound,
		fmt.Errorf(
			"failure finding matching metal machine for uuid: %q",
			uuid,
		),
	)
}
