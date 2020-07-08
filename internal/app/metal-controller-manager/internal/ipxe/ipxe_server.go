/*
Copyright 2020 Talos Systems, Inc.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ipxe

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"text/template"

	metalv1alpha1 "github.com/talos-systems/sidero/internal/app/metal-controller-manager/api/v1alpha1"
	"github.com/talos-systems/sidero/internal/app/metal-controller-manager/internal/server"
	agentclient "github.com/talos-systems/sidero/internal/app/metal-controller-manager/pkg/client"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const bootFile = `#!ipxe
chain ipxe?uuid=${uuid}&mac=${mac:hexhyp}&domain=${domain}&hostname=${hostname}&serial=${serial}
`

var ipxeTemplate = template.Must(template.New("iPXE config").Parse(`#!ipxe
kernel /env/{{ .Env.Name }}/vmlinuz {{range $arg := .Env.Spec.Kernel.Args}} {{$arg}}{{end}}
initrd /env/{{ .Env.Name }}/initramfs.xz
boot
`))

var (
	apiEndpoint string
)

func bootFileHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, bootFile)
}

func ipxeHandler(w http.ResponseWriter, r *http.Request) {
	var (
		config *rest.Config
		err    error
	)
	kubeconfig, ok := os.LookupEnv("KUBECONFIG")
	if ok {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			log.Printf("error creating config: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		config, err = rest.InClusterConfig()
		if err != nil {
			log.Printf("error creating config: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	c, err := agentclient.NewClient(config)
	if err != nil {
		log.Printf("error creating client: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	labels := labelsFromRequest(r)

	log.Printf("UUID: %q", labels["uuid"])

	key := client.ObjectKey{
		Name: labels["uuid"],
	}

	obj := &metalv1alpha1.Server{}

	if err := c.Get(context.Background(), key, obj); err != nil {
		// If we can't find the server then we know that discovery has not been
		// performed yet.
		if apierrors.IsNotFound(err) {
			var args = struct {
				Env metalv1alpha1.Environment
			}{
				Env: metalv1alpha1.Environment{
					ObjectMeta: v1.ObjectMeta{
						Name: "discovery",
					},
					Spec: metalv1alpha1.EnvironmentSpec{
						Kernel: metalv1alpha1.Kernel{
							Args: []string{
								"initrd=initramfs.xz",
								"page_poison=1",
								"slab_nomerge",
								"slub_debug=P",
								"pti=on",
								"panic=0",
								"random.trust_cpu=on",
								"ima_template=ima-ng",
								"ima_appraise=fix",
								"ima_hash=sha512",
								"ip=dhcp",
								"console=tty0",
								"console=ttyS0",
								"arges.endpoint=" + fmt.Sprintf("%s:%s", apiEndpoint, server.Port),
							},
						},
					},
				},
			}

			var buf bytes.Buffer
			err = ipxeTemplate.Execute(&buf, args)
			if err != nil {
				log.Printf("error rendering template: %v", err)
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

			if _, err := buf.WriteTo(w); err != nil {
				log.Printf("error writing to response: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
			}

			return
		} else {
			log.Printf("error looking up server: %v", err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}
	}

	var env metalv1alpha1.Environment

	if err := c.Get(context.Background(), types.NamespacedName{Namespace: "", Name: "default"}, &env); err != nil {
		if apierrors.IsNotFound(err) {
			log.Printf("environment not found: %v", err)
			w.WriteHeader(http.StatusNotFound)

			return
		}
	}

	var args = struct {
		Env metalv1alpha1.Environment
	}{
		Env: env,
	}

	var buf bytes.Buffer
	err = ipxeTemplate.Execute(&buf, args)
	if err != nil {
		log.Printf("error rendering template: %v", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	if _, err := buf.WriteTo(w); err != nil {
		log.Printf("error writing to response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func ServeIPXE(endpoint string) error {
	apiEndpoint = endpoint

	mux := http.NewServeMux()

	mux.Handle("/boot.ipxe", logRequest(http.HandlerFunc(bootFileHandler)))
	mux.Handle("/ipxe", logRequest(http.HandlerFunc(ipxeHandler)))
	mux.Handle("/env/", logRequest(http.StripPrefix("/env/", http.FileServer(http.Dir("/var/lib/arges/env")))))

	log.Println("Listening...")
	return http.ListenAndServe(":8081", mux)
}

func logRequest(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		log.Printf("HTTP %s %v %s", r.Method, r.URL, r.RemoteAddr)
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func labelsFromRequest(req *http.Request) map[string]string {
	values := req.URL.Query()
	labels := map[string]string{}
	for key := range values {
		switch strings.ToLower(key) {
		case "mac":
			// set mac if and only if it parses
			if hw, err := parseMAC(values.Get(key)); err == nil {
				labels[key] = hw.String()
			}
		default:
			// matchers don't use multi-value keys, drop later values
			labels[key] = values.Get(key)
		}
	}
	return labels
}

func parseMAC(s string) (net.HardwareAddr, error) {
	macAddr, err := net.ParseMAC(s)
	if err != nil {
		return nil, err
	}
	return macAddr, err
}
