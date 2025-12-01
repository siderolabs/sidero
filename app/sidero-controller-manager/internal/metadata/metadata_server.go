// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package metadata

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/siderolabs/talos/pkg/machinery/config/configpatcher"
	"gopkg.in/yaml.v3"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/cluster-api/util"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "github.com/siderolabs/sidero/app/caps-controller-manager/api/v1alpha3"
	metalv1 "github.com/siderolabs/sidero/app/sidero-controller-manager/api/v1alpha2"
	"github.com/siderolabs/sidero/app/sidero-controller-manager/internal/siderolink"
)

type errorWithCode struct {
	errorCode int
	errorObj  error
}

type metadataConfigs struct {
	client      runtimeclient.Client
	apiEndpoint string
	apiPort     int
}

func throwError(w http.ResponseWriter, ewc errorWithCode) {
	http.Error(w, ewc.errorObj.Error(), ewc.errorCode)
	log.Println(ewc.errorObj)
}

func RegisterServer(mux *http.ServeMux, k8sClient runtimeclient.Client, apiEndpoint string, apiPort int) error {
	mm := metadataConfigs{
		client:      k8sClient,
		apiEndpoint: apiEndpoint,
		apiPort:     apiPort,
	}

	mux.HandleFunc("/configdata", mm.FetchConfig)

	return nil
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

	// Find serverBinding and metalMachine by server UUID.
	metalMachine, serverBinding, ewc := m.findMetalMachineServerBinding(ctx, uuid)
	if ewc.errorObj != nil {
		throwError(
			w,
			ewc,
		)

		return
	}

	// Given the MetalMachine, find the Machine resource that owns it
	ownerMachine, err := util.GetOwnerMachine(ctx, m.client, metalMachine.ObjectMeta)
	if err != nil || ownerMachine == nil {
		throwError(
			w,
			errorWithCode{
				http.StatusInternalServerError,
				fmt.Errorf(
					"failure fetching owner machine from metal machine %s/%s: %s",
					metalMachine.GetNamespace(),
					metalMachine.GetName(),
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
	serverObj := &metalv1.Server{}

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
	serverClassObj := &metalv1.ServerClass{}

	if serverBinding.Spec.ServerClassRef != nil {
		err = m.client.Get(
			ctx,
			types.NamespacedName{
				Namespace: serverBinding.Spec.ServerClassRef.Namespace,
				Name:      serverBinding.Spec.ServerClassRef.Name,
			},
			serverClassObj,
		)
		if err != nil {
			throwError(
				w,
				errorWithCode{
					http.StatusInternalServerError,
					fmt.Errorf(
						"failure fetching serverclass %s: %s",
						serverBinding.Spec.ServerClassRef.Name,
						err,
					),
				},
			)

			return
		}
	}

	decodedData, ewc = handlePatches(decodedData, serverClassObj.Spec.ConfigPatches, serverClassObj.Spec.StrategicPatches)
	if ewc.errorObj != nil {
		throwError(
			w,
			ewc,
		)

		return
	}

	decodedData, ewc = handlePatches(decodedData, serverObj.Spec.ConfigPatches, serverObj.Spec.StrategicPatches)
	if ewc.errorObj != nil {
		throwError(
			w,
			ewc,
		)

		return
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

	// Patch machine configuration with SideroLink config so it survives reboots
	// TODO(laurazard): only do this if Talos v1.10+
	decodedData, ewc = m.patchSideroLinkConfig(decodedData)
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

// this function is responsible for applying rfc6902 and a strategic merge patch to bootstrap data.
func handlePatches(decodedData []byte, patches []metalv1.ConfigPatches, strategicPatches []string) ([]byte, errorWithCode) {
	var ewc errorWithCode

	// Handle rfc6902 patches
	if len(patches) > 0 {
		decodedData, ewc = patchRFC6902Configs(decodedData, patches)
		if ewc.errorObj != nil {
			return decodedData, ewc
		}
	}

	// Handle strategic merge patch
	for _, strategicPatch := range strategicPatches {
		decodedData, ewc = patchConfig(decodedData, []byte(strategicPatch))
		if ewc.errorObj != nil {
			return decodedData, ewc
		}
	}

	return decodedData, errorWithCode{}
}

// patchRFC6902Configs is responsible for applying rfc6902 configPatches to the bootstrap data.
func patchRFC6902Configs(decodedData []byte, patches []metalv1.ConfigPatches) ([]byte, errorWithCode) {
	marshalledPatches, err := json.Marshal(patches)
	if err != nil {
		return nil, errorWithCode{http.StatusInternalServerError, fmt.Errorf("failure marshaling config patches from server: %s", err)}
	}

	return patchConfig(decodedData, marshalledPatches)
}

// patchConfig is responsible for applying marshaled rfc6902 configPatches or a strategic merge patch to the bootstrap data.
func patchConfig(decodedData []byte, patches []byte) ([]byte, errorWithCode) {
	patch, err := configpatcher.LoadPatch(patches)
	if err != nil {
		return nil, errorWithCode{http.StatusInternalServerError, fmt.Errorf("failure loading patches from server: %s", err)}
	}

	patched, err := configpatcher.Apply(configpatcher.WithBytes(decodedData), []configpatcher.Patch{patch})
	if err != nil {
		return nil, errorWithCode{http.StatusInternalServerError, fmt.Errorf("failure applying patches %s to machine config: %s", string(patches), err)}
	}

	result, err := patched.Bytes()
	if err != nil {
		return nil, errorWithCode{http.StatusInternalServerError, fmt.Errorf("failure converting patched config to bytes: %s", err)}
	}

	return result, errorWithCode{}
}

// labelNodes is responsible for editing the kubelet extra args such that a given
// server gets registered with a label containing the UUID of the server resource it's actually running on.
func labelNodes(decodedData []byte, serverName string) ([]byte, errorWithCode) {
	isMultiDoc := bytes.Contains(decodedData, []byte("---\n"))
	if isMultiDoc {
		return labelNodesStrategic(decodedData, serverName)
	}

	return labelNodesLegacy(decodedData, serverName)
}

// labelNodes is responsible for editing the kubelet extra args such that a given
// server gets registered with a label containing the UUID of the server resource it's actually running on.
func labelNodesLegacy(decodedData []byte, serverName string) ([]byte, errorWithCode) {
	// avoid using the `configloader` from Talos machinery here, as it will fail on "unknown" fields
	// causing a dependency on Talos version that Sidero was built with
	var cfg struct {
		Version string `yaml:"version"`
		Machine *struct {
			Kubelet *struct {
				ExtraArgs map[string]string `yaml:"extraArgs"`
			} `yaml:"kubelet"`
		} `yaml:"machine"`
	}

	if err := yaml.Unmarshal(decodedData, &cfg); err != nil {
		return nil, errorWithCode{http.StatusInternalServerError, fmt.Errorf("failure creating config struct: %s", err)}
	}

	switch cfg.Version {
	case "v1alpha1":
		var (
			patch      metalv1.ConfigPatches
			patchValue any
		)

		label := fmt.Sprintf("metal.sidero.dev/uuid=%s", serverName)

		switch {
		case cfg.Machine == nil:
			patch = metalv1.ConfigPatches{
				Path: "/machine",
				Op:   "add",
			}

			patchValue = map[string]any{
				"kubelet": map[string]any{
					"extraArgs": map[string]string{
						"node-labels": label,
					},
				},
			}
		case cfg.Machine.Kubelet == nil:
			patch = metalv1.ConfigPatches{
				Path: "/machine/kubelet",
				Op:   "add",
			}

			patchValue = map[string]any{
				"extraArgs": map[string]string{
					"node-labels": label,
				},
			}
		case cfg.Machine.Kubelet.ExtraArgs == nil:
			patch = metalv1.ConfigPatches{
				Path: "/machine/kubelet/extraArgs",
				Op:   "add",
			}

			patchValue = map[string]any{
				"node-labels": label,
			}
		default:
			kubeletExtraArgs := cfg.Machine.Kubelet.ExtraArgs
			if _, ok := kubeletExtraArgs["node-labels"]; ok {
				kubeletExtraArgs["node-labels"] += "," + label
			} else {
				kubeletExtraArgs["node-labels"] = label
			}

			patch = metalv1.ConfigPatches{
				Path: "/machine/kubelet/extraArgs",
				Op:   "replace",
			}

			patchValue = kubeletExtraArgs
		}

		value, err := json.Marshal(patchValue)
		if err != nil {
			return nil, errorWithCode{http.StatusInternalServerError, fmt.Errorf("failure marshaling kubelet.extraArgs: %s", err)}
		}

		patch.Value.Raw = value

		return patchRFC6902Configs(decodedData, []metalv1.ConfigPatches{patch})
	default:
		return nil, errorWithCode{http.StatusInternalServerError, fmt.Errorf("unknown config type")}
	}
}

func labelNodesStrategic(decodedData []byte, serverName string) ([]byte, errorWithCode) {
	// avoid using the `configloader` from Talos machinery here, as it will fail on "unknown" fields
	// causing a dependency on Talos version that Sidero was built with
	var cfg struct {
		Version string `yaml:"version"`
		Machine *struct {
			Kubelet *struct {
				ExtraArgs map[string]string `yaml:"extraArgs"`
			} `yaml:"kubelet"`
		} `yaml:"machine"`
	}

	if err := yaml.Unmarshal(decodedData, &cfg); err != nil {
		return nil, errorWithCode{http.StatusInternalServerError, fmt.Errorf("failure creating config struct: %s", err)}
	}

	label := fmt.Sprintf("metal.sidero.dev/uuid=%s", serverName)

	patch := map[string]any{
		"machine": map[string]any{
			"kubelet": map[string]any{
				"extraArgs": map[string]*string{
					"node-labels": &label,
				},
			},
		},
	}

	switch cfg.Version {
	case "v1alpha1":
		switch {
		case cfg.Machine == nil:
		case cfg.Machine.Kubelet == nil:
		case cfg.Machine.Kubelet.ExtraArgs == nil:
		default:
			kubeletExtraArgs := cfg.Machine.Kubelet.ExtraArgs
			if _, ok := kubeletExtraArgs["node-labels"]; ok {
				label = kubeletExtraArgs["node-labels"] + "," + label
			}
		}

		patchMarshaled, err := yaml.Marshal(patch)
		if err != nil {
			return nil, errorWithCode{http.StatusInternalServerError, fmt.Errorf("failure marshaling kubelet.extraArgs: %s", err)}
		}

		return patchConfig(decodedData, patchMarshaled)
	default:
		return nil, errorWithCode{http.StatusInternalServerError, fmt.Errorf("unknown config type")}
	}
}

// findMetalMachineServerBinding is responsible for looking up ServerBinding and MetalMachine.
func (m *metadataConfigs) findMetalMachineServerBinding(ctx context.Context, serverName string) (infrav1.MetalMachine, infrav1.ServerBinding, errorWithCode) {
	var serverBinding infrav1.ServerBinding

	err := m.client.Get(ctx, types.NamespacedName{Name: serverName}, &serverBinding)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return infrav1.MetalMachine{}, infrav1.ServerBinding{}, errorWithCode{http.StatusNotFound, fmt.Errorf("server is not allocated (missing serverbinding): %w", err)}
		}

		return infrav1.MetalMachine{}, infrav1.ServerBinding{}, errorWithCode{http.StatusInternalServerError, fmt.Errorf("failure getting server binding: %w", err)}
	}

	var metalMachine infrav1.MetalMachine

	if err = m.client.Get(ctx, types.NamespacedName{
		// XXX: where is the namespace in owner refs?
		Namespace: serverBinding.Spec.MetalMachineRef.Namespace,
		Name:      serverBinding.Spec.MetalMachineRef.Name,
	}, &metalMachine); err != nil {
		return infrav1.MetalMachine{}, infrav1.ServerBinding{}, errorWithCode{http.StatusInternalServerError, fmt.Errorf("failure getting metalmachine: %w", err)}
	}

	return metalMachine, serverBinding, errorWithCode{}
}

func (m *metadataConfigs) patchSideroLinkConfig(decodedData []byte) ([]byte, errorWithCode) {
	var ewc errorWithCode

	var cfg struct {
		Version string `yaml:"version"`
	}

	if err := yaml.Unmarshal(decodedData, &cfg); err != nil {
		return nil, errorWithCode{http.StatusInternalServerError, fmt.Errorf("failure decoding config struct: %s", err)}
	}

	if cfg.Version != "v1alpha1" {
		return nil, errorWithCode{http.StatusInternalServerError, fmt.Errorf("unknown config type")}
	}

	sideroLinkPatch := map[string]string{
		"apiVersion": "v1alpha1",
		"kind":       "SideroLinkConfig",
		"apiUrl":     fmt.Sprintf("%s:%d", m.apiEndpoint, m.apiPort),
	}

	sideroLinkPatchB, err := yaml.Marshal(sideroLinkPatch)
	if err != nil {
		return decodedData, errorWithCode{
			errorCode: http.StatusInternalServerError,
			errorObj:  err,
		}
	}

	decodedData, ewc = patchConfig(decodedData, sideroLinkPatchB)
	if ewc.errorObj != nil {
		return decodedData, ewc
	}

	eventsSinkPatch := map[string]string{
		"apiVersion": "v1alpha1",
		"kind":       "EventSinkConfig",
		"endpoint":   fmt.Sprintf("[%s]:%d", siderolink.Cfg.ServerAddress.Addr(), siderolink.EventsSinkPort),
	}

	eventsSinkPatchB, err := yaml.Marshal(eventsSinkPatch)
	if err != nil {
		return decodedData, errorWithCode{
			errorCode: http.StatusInternalServerError,
			errorObj:  err,
		}
	}

	decodedData, ewc = patchConfig(decodedData, eventsSinkPatchB)
	if ewc.errorObj != nil {
		return decodedData, ewc
	}

	kmsLogPatch := map[string]string{
		"apiVersion": "v1alpha1",
		"kind":       "KmsgLogConfig",
		"name":       "remote-log",
		"url":        fmt.Sprintf("tcp://[%s]:%d", siderolink.Cfg.ServerAddress.Addr(), siderolink.LogReceiverPort),
	}

	kmsLogPatchB, err := yaml.Marshal(kmsLogPatch)
	if err != nil {
		return decodedData, errorWithCode{
			errorCode: http.StatusInternalServerError,
			errorObj:  err,
		}
	}

	decodedData, ewc = patchConfig(decodedData, kmsLogPatchB)
	if ewc.errorObj != nil {
		return decodedData, ewc
	}

	return decodedData, errorWithCode{}
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
