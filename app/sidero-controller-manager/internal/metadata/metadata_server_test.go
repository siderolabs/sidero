// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package metadata_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	capiv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	infrav1 "github.com/siderolabs/sidero/app/caps-controller-manager/api/v1alpha3"
	metalv1 "github.com/siderolabs/sidero/app/sidero-controller-manager/api/v1alpha2"
	"github.com/siderolabs/sidero/app/sidero-controller-manager/internal/metadata"
	"github.com/siderolabs/sidero/app/sidero-controller-manager/internal/siderolink"
)

var sideroLinkCfgs = []map[string]any{
	{
		"apiVersion": "v1alpha1",
		"kind":       "SideroLinkConfig",
		"apiUrl":     "grpc://192.168.1.1:8081",
	},
	{
		"apiVersion": "v1alpha1",
		"kind":       "EventSinkConfig",
		"endpoint":   "192.168.1.1:4002",
	},
	{
		"apiVersion": "v1alpha1",
		"kind":       "KmsgLogConfig",
		"name":       "remote-log",
		"url":        "tcp://192.168.1.1:4001",
	},
}

var extensionServiceCfg = map[string]any{
	"apiVersion": "v1alpha1",
	"kind":       "ExtensionServiceConfig",
	"name":       "frr",
	"environment": []any{
		"TESTKEY=TESTVALUE",
	},
}

func TestMetadataService(t *testing.T) {
	oldSideroCfg := siderolink.Cfg
	siderolink.Cfg = siderolink.Config{
		ServerAddress: netip.PrefixFrom(netip.AddrFrom4([4]byte{192, 168, 1, 1}), 32),
	}
	// restore config after test
	t.Cleanup(func() {
		siderolink.Cfg = oldSideroCfg
	})

	t.Parallel()

	scheme := runtime.NewScheme()
	require.NoError(t, infrav1.AddToScheme(scheme))
	require.NoError(t, metalv1.AddToScheme(scheme))
	require.NoError(t, capiv1.AddToScheme(scheme))
	require.NoError(t, corev1.AddToScheme(scheme))

	fakeClient := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(
			fixture()...,
		).
		Build()

	mux := http.NewServeMux()

	metadata.RegisterServer(mux, fakeClient, "192.168.1.1", 8081)

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	tests := []struct {
		name            string
		path            string
		expectedCode    int
		expectedBody    string
		expectedConfigs []map[string]any
	}{
		{
			name:         "invalid",
			path:         "/configdata",
			expectedCode: http.StatusInternalServerError,
			expectedBody: "received metadata request with empty uuid\n",
		},
		{
			name:         "not found",
			path:         "/configdata?uuid=xxx-yyy",
			expectedCode: http.StatusNotFound,
			expectedBody: "server is not allocated (missing serverbinding): serverbindings.infrastructure.cluster.x-k8s.io \"xxx-yyy\" not found\n",
		},
		{
			name:         "no patches",
			path:         "/configdata?uuid=0000-1111-2222",
			expectedCode: http.StatusOK,
			expectedConfigs: append([]map[string]any{
				{
					"version": "v1alpha1",
					"cluster": nil,
					"machine": map[string]any{
						"certSANs": []any{},
						"kubelet": map[string]any{
							"extraArgs": map[string]any{
								"node-labels": "metal.sidero.dev/uuid=0000-1111-2222",
							},
						},
						"token": "",
						"type":  "",
					},
				},
			}, sideroLinkCfgs...),
		},
		{
			name:         "server patch",
			path:         "/configdata?uuid=1111-2222-3333",
			expectedCode: http.StatusOK,
			expectedConfigs: append([]map[string]any{
				{
					"version": "v1alpha1",
					"cluster": nil,
					"machine": map[string]any{
						"certSANs": []any{},
						"kubelet": map[string]any{
							"extraArgs": map[string]any{
								"foo":         "bar",
								"node-labels": "metal.sidero.dev/uuid=1111-2222-3333",
							},
						},
						"network": map[string]any{
							"hostname": "example2",
						},
						"token": "",
						"type":  "",
					},
				},
			}, sideroLinkCfgs...),
		},
		{
			name:         "server and server class patch",
			path:         "/configdata?uuid=2222-3333-4444",
			expectedCode: http.StatusOK,
			expectedConfigs: append([]map[string]any{
				{
					"version": "v1alpha1",
					"cluster": nil,
					"machine": map[string]any{
						"certSANs": []any{},
						"kubelet": map[string]any{
							"extraArgs": map[string]any{
								"node-labels": "foo=bar,metal.sidero.dev/uuid=2222-3333-4444",
							},
						},
						"network": map[string]any{
							"hostname": "example3",
						},
						"token": "",
						"type":  "",
					},
				},
			}, sideroLinkCfgs...),
		},
		{
			name:         "machine config without kubelet",
			path:         "/configdata?uuid=4444-5555-6666",
			expectedCode: http.StatusOK,
			expectedConfigs: append([]map[string]any{
				{
					"version": "v1alpha1",
					"cluster": nil,
					"machine": map[string]any{
						"certSANs": []any{},
						"kubelet": map[string]any{
							"extraArgs": map[string]any{
								"node-labels": "metal.sidero.dev/uuid=4444-5555-6666",
							},
						},
						"token": "",
						"type":  "",
					},
				},
			}, sideroLinkCfgs...),
		},
		{
			name:         "machine config without machine",
			path:         "/configdata?uuid=5555-6666-7777",
			expectedCode: http.StatusOK,
			expectedConfigs: append([]map[string]any{{
				"version": "v1alpha1",
				"cluster": map[string]any{
					"controlPlane": nil,
				},
				"machine": map[string]any{
					"certSANs": []any{},
					"kubelet": map[string]any{
						"extraArgs": map[string]any{
							"node-labels": "metal.sidero.dev/uuid=5555-6666-7777",
						},
					},
					"token": "",
					"type":  "",
				},
			}}, sideroLinkCfgs...),
		},
		{
			name:         "server and server class as strategic merge patch",
			path:         "/configdata?uuid=6666-7777-8888",
			expectedCode: http.StatusOK,
			expectedConfigs: append([]map[string]any{
				{
					"version": "v1alpha1",
					"cluster": nil,

					"machine": map[string]any{
						"token":    "",
						"type":     "",
						"certSANs": []any{},
						"kubelet": map[string]any{
							"extraArgs": map[string]any{
								"node-labels": "foo=bar,metal.sidero.dev/uuid=6666-7777-8888",
							},
						},
						"network": map[string]any{
							"hostname": "example6",
						},
					},
				},
				extensionServiceCfg,
			}, sideroLinkCfgs...),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			resp, err := http.Get(srv.URL + test.path) //nolint:noctx
			require.NoError(t, err)

			t.Cleanup(func() { resp.Body.Close() })

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.Equal(t, test.expectedCode, resp.StatusCode)

			if test.expectedBody != "" {
				assert.Equal(t, test.expectedBody, string(body))
			} else if len(test.expectedConfigs) > 0 {
				docs := parseYAMLDocs(t, body)
				for i := 0; i < len(test.expectedConfigs); i++ {
					require.EqualValues(t, test.expectedConfigs[i], docs[i], fmt.Sprintf("actual:\n%s\n", string(body)))
				}
			}
		})
	}
}

func parseYAMLDocs(t *testing.T, data []byte) []map[string]any {
	t.Helper()

	var (
		decoder = yaml.NewDecoder(strings.NewReader(string(data)))
		docs    []map[string]any
	)

	for {
		var doc map[string]any

		err := decoder.Decode(&doc)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}

			require.NoError(t, err)
		}

		docs = append(docs, doc)
	}

	return docs
}
