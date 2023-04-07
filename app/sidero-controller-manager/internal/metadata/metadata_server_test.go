// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package metadata_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	capiv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	infrav1 "github.com/siderolabs/sidero/app/caps-controller-manager/api/v1alpha3"
	metalv1 "github.com/siderolabs/sidero/app/sidero-controller-manager/api/v1alpha2"
	"github.com/siderolabs/sidero/app/sidero-controller-manager/internal/metadata"
)

func TestMetadataService(t *testing.T) {
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

	metadata.RegisterServer(mux, fakeClient)

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	for _, test := range []struct {
		name string
		path string

		expectedCode int
		expectedBody string
	}{
		{
			name: "invalid",
			path: "/configdata",

			expectedCode: http.StatusInternalServerError,
			expectedBody: "received metadata request with empty uuid\n",
		},
		{
			name: "not found",
			path: "/configdata?uuid=xxx-yyy",

			expectedCode: http.StatusNotFound,
			expectedBody: "server is not allocated (missing serverbinding): serverbindings.infrastructure.cluster.x-k8s.io \"xxx-yyy\" not found\n",
		},
		{
			name: "no patches",
			path: "/configdata?uuid=0000-1111-2222",

			expectedCode: http.StatusOK,
			expectedBody: "machine:\n  kubelet:\n    extraArgs:\n      node-labels: metal.sidero.dev/uuid=0000-1111-2222\nversion: v1alpha1\n",
		},
		{
			name: "server patch",
			path: "/configdata?uuid=1111-2222-3333",

			expectedCode: http.StatusOK,
			expectedBody: "machine:\n  kubelet:\n    extraArgs:\n      foo: bar\n      node-labels: metal.sidero.dev/uuid=1111-2222-3333\n  network:\n    hostname: example2\nversion: v1alpha1\n",
		},
		{
			name: "server and server class patch",
			path: "/configdata?uuid=2222-3333-4444",

			expectedCode: http.StatusOK,
			expectedBody: "machine:\n  kubelet:\n    extraArgs:\n      node-labels: foo=bar,metal.sidero.dev/uuid=2222-3333-4444\n  network:\n    hostname: example3\nversion: v1alpha1\n",
		},
		{
			name: "machine config without kubelet",
			path: "/configdata?uuid=4444-5555-6666",

			expectedCode: http.StatusOK,
			expectedBody: "machine:\n  kubelet:\n    extraArgs:\n      node-labels: metal.sidero.dev/uuid=4444-5555-6666\n  unsupported: {}\nversion: v1alpha1\n",
		},
		{
			name: "machine config without machine",
			path: "/configdata?uuid=5555-6666-7777",

			expectedCode: http.StatusOK,
			expectedBody: "cluster: {}\nmachine:\n  kubelet:\n    extraArgs:\n      node-labels: metal.sidero.dev/uuid=5555-6666-7777\nversion: v1alpha1\n",
		},
	} {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			resp, err := http.Get(srv.URL + test.path) //nolint:noctx
			require.NoError(t, err)

			t.Cleanup(func() { resp.Body.Close() })

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.Equal(t, test.expectedCode, resp.StatusCode)
			assert.Equal(t, test.expectedBody, string(body))
		})
	}
}
