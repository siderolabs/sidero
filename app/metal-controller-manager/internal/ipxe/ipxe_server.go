// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ipxe

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"text/template"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/controller-runtime/pkg/client"

	metalv1alpha1 "github.com/talos-systems/sidero/app/metal-controller-manager/api/v1alpha1"
	"github.com/talos-systems/sidero/app/metal-controller-manager/internal/server"
	agentclient "github.com/talos-systems/sidero/app/metal-controller-manager/pkg/client"
	"github.com/talos-systems/sidero/app/metal-controller-manager/pkg/constants"
)

var (
	ErrNotInUse     = errors.New("server not in use")
	ErrBootFromDisk = errors.New("boot from disk")
)

const bootFile = `#!ipxe
chain ipxe?uuid=${uuid}&mac=${mac:hexhyp}&domain=${domain}&hostname=${hostname}&serial=${serial}
`

var ipxeTemplate = template.Must(template.New("iPXE config").Parse(`#!ipxe
kernel /env/{{ .Env.Name }}/{{ .KernelAsset }} {{range $arg := .Env.Spec.Kernel.Args}} {{$arg}}{{end}}
initrd /env/{{ .Env.Name }}/{{ .InitrdAsset }}
boot
`))

const ipxeBootFromDisk = `#!ipxe
exit
`

var (
	apiEndpoint          string
	extraAgentKernelArgs []string
)

func bootFileHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, bootFile)
}

//nolint: unparam
func bootFromDiskHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, ipxeBootFromDisk)
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
			log.Printf("Error creating config: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		config, err = rest.InClusterConfig()
		if err != nil {
			log.Printf("Error creating config: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	c, err := agentclient.NewClient(config)
	if err != nil {
		log.Printf("Error creating client: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	labels := labelsFromRequest(r)

	uuid := labels["uuid"]

	obj, err := lookupServer(c, uuid)
	if err != nil {
		log.Printf("Error looking up server: %v", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	env, err := newEnvironment(c, obj)
	if err != nil {
		if errors.Is(err, ErrBootFromDisk) {
			log.Printf("Server %q booting from disk", uuid)
			bootFromDiskHandler(w, r)

			return
		}

		if apierrors.IsNotFound(err) {
			log.Printf("Environment not found: %v", err)
			w.WriteHeader(http.StatusNotFound)

			return
		}

		if errors.Is(err, ErrNotInUse) {
			log.Printf("Server %q not in use, skipping", uuid)
			w.WriteHeader(http.StatusNotFound)

			return
		}

		log.Printf("%v", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	if obj != nil {
		log.Printf("Using %q environment for %q", env.Name, obj.Name)
	} else {
		log.Printf("Using %q environment", env.Name)
	}

	args := struct {
		Env         *metalv1alpha1.Environment
		KernelAsset string
		InitrdAsset string
	}{
		Env:         env,
		KernelAsset: constants.KernelAsset,
		InitrdAsset: constants.InitrdAsset,
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

		return
	}

	if env.ObjectMeta.Name != "agent" {
		if err = markAsPXEBooted(c, obj); err != nil {
			log.Printf("error marking server as PXE booted: %s", err)
		}
	}
}

func ServeIPXE(endpoint string, args []string) error {
	apiEndpoint = endpoint
	extraAgentKernelArgs = args

	mux := http.NewServeMux()

	mux.Handle("/boot.ipxe", logRequest(http.HandlerFunc(bootFileHandler)))
	mux.Handle("/ipxe", logRequest(http.HandlerFunc(ipxeHandler)))
	mux.Handle("/env/", logRequest(http.StripPrefix("/env/", http.FileServer(http.Dir("/var/lib/sidero/env")))))

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

func lookupServer(c client.Client, uuid string) (*metalv1alpha1.Server, error) {
	key := client.ObjectKey{
		Name: uuid,
	}

	obj := &metalv1alpha1.Server{}

	if err := c.Get(context.Background(), key, obj); err != nil {
		if apierrors.IsNotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	return obj, nil
}

// newEnvironment handles which env CRD we'll respect for a given server.
// specied in the server spec overrides everything, specified in the server class overrides default, default is default :).
func newEnvironment(c client.Client, server *metalv1alpha1.Server) (env *metalv1alpha1.Environment, err error) {
	// NB: The order of this switch statement is important. It defines the
	// precedence of which environment to boot.
	switch {
	case server == nil:
		return newAgentEnvironment(), nil
	case !server.Status.InUse && !server.Status.IsClean:
		return newAgentEnvironment(), nil
	case !server.Status.InUse:
		return nil, ErrNotInUse
	case conditions.Has(server, metalv1alpha1.ConditionPXEBooted) && !server.Spec.PXEBootAlways:
		return nil, ErrBootFromDisk
	case server.Spec.EnvironmentRef != nil:
		env, err = newEnvironmentFromServer(c, server)
		if err != nil {
			return nil, err
		}
	case server.OwnerReferences != nil:
		env, err = newEnvironmentFromServerClass(c, server)
		if err != nil {
			return nil, err
		}
	}

	if env == nil {
		env, err = newDefaultEnvironment(c)
		if err != nil {
			return nil, err
		}
	}

	if env == nil {
		return nil, fmt.Errorf("could not find environment for %q", server.Name)
	}

	return env, nil
}

func newAgentEnvironment() *metalv1alpha1.Environment {
	args := []string{
		"initrd=initramfs.xz",
		"page_poison=1",
		"slab_nomerge",
		"slub_debug=P",
		"pti=on",
		"panic=30",
		"random.trust_cpu=on",
		"ima_template=ima-ng",
		"ima_appraise=fix",
		"ima_hash=sha512",
		"ip=dhcp",
		"console=tty0",
		"console=ttyS0",
		"printk.devkmsg=on",
		fmt.Sprintf("%s=%s:%s", constants.AgentEndpointArg, apiEndpoint, server.Port),
	}

	args = append(args, extraAgentKernelArgs...)

	env := &metalv1alpha1.Environment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "agent",
		},
		Spec: metalv1alpha1.EnvironmentSpec{
			Kernel: metalv1alpha1.Kernel{
				Args: args,
			},
		},
	}

	return env
}

func newDefaultEnvironment(c client.Client) (env *metalv1alpha1.Environment, err error) {
	env = &metalv1alpha1.Environment{}

	if err := c.Get(context.Background(), types.NamespacedName{Namespace: "", Name: "default"}, env); err != nil {
		return nil, err
	}

	return env, nil
}

func newEnvironmentFromServer(c client.Client, server *metalv1alpha1.Server) (env *metalv1alpha1.Environment, err error) {
	env = &metalv1alpha1.Environment{}

	if err := c.Get(context.Background(), types.NamespacedName{Namespace: "", Name: server.Spec.EnvironmentRef.Name}, env); err != nil {
		return nil, err
	}

	return env, nil
}

func newEnvironmentFromServerClass(c client.Client, server *metalv1alpha1.Server) (env *metalv1alpha1.Environment, err error) {
	for _, owner := range server.OwnerReferences {
		if owner.Kind == "ServerClass" {
			serverClassResource := &metalv1alpha1.ServerClass{}

			if err := c.Get(context.Background(), types.NamespacedName{Namespace: "", Name: owner.Name}, serverClassResource); err != nil {
				return nil, err
			}

			if serverClassResource.Spec.EnvironmentRef != nil {
				env = &metalv1alpha1.Environment{}

				if err := c.Get(context.Background(), types.NamespacedName{Namespace: "", Name: serverClassResource.Spec.EnvironmentRef.Name}, env); err != nil {
					return nil, err
				}

				break
			}
		}
	}

	return env, nil
}

func markAsPXEBooted(c client.Client, server *metalv1alpha1.Server) error {
	patchHelper, err := patch.NewHelper(server, c)
	if err != nil {
		return err
	}

	conditions.MarkTrue(server, metalv1alpha1.ConditionPXEBooted)

	return patchHelper.Patch(context.Background(), server, patch.WithOwnedConditions{
		Conditions: []clusterv1.ConditionType{metalv1alpha1.ConditionPXEBooted},
	})
}
