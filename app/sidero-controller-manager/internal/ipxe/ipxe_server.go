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
	"strconv"
	"strings"
	"text/template"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha4"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/talos-systems/go-procfs/procfs"
	talosconstants "github.com/talos-systems/talos/pkg/machinery/constants"

	infrav1 "github.com/talos-systems/sidero/app/caps-controller-manager/api/v1alpha3"
	metalv1alpha1 "github.com/talos-systems/sidero/app/sidero-controller-manager/api/v1alpha1"
	"github.com/talos-systems/sidero/app/sidero-controller-manager/pkg/constants"
)

var ErrBootFromDisk = errors.New("boot from disk")

// bootFile is used when iPXE is booted without embedded script via iPXE request http://endpoint:8081/boot.ipxe.
const bootFile = `#!ipxe
chain ipxe?uuid=${uuid}&mac=${mac:hexhyp}&domain=${domain}&hostname=${hostname}&serial=${serial}&arch=${buildarch}
`

// BootTemplate is embedded into iPXE binary when that binary is sent to the node.
var BootTemplate = template.Must(template.New("iPXE embedded").Parse(`
prompt --key 0x02 --timeout 2000 Press Ctrl-B for the iPXE command line... && shell ||

# print interfaces
ifstat

# retry 10 times overall
set attempts:int32 10
set x:int32 0

:retry_loop

	set idx:int32 0

	:loop
		# try DHCP on each available interface
		isset ${net${idx}/mac} || goto exhausted

		ifclose
		dhcp net${idx} && goto boot

	:next_iface
		inc idx && goto loop

	:boot
		# attempt boot, if fails try next iface
		route

		chain http://{{ .Endpoint }}:{{ .Port }}/ipxe?uuid=${uuid}&mac=${net${idx}/mac:hexhyp}&domain=${domain}&hostname=${hostname}&serial=${serial}&arch=${buildarch} || goto next_iface

:exhausted
	echo
	echo Failed to iPXE boot successfully via all interfaces

	iseq ${x} ${attempts} && goto fail ||

	echo Retrying...
	echo

	inc x
	goto retry_loop

:fail
	echo
	echo Failed to get a valid response after ${attempts} attempts
	echo

	echo Rebooting in 5 seconds...
	sleep 5
	reboot
`))

// ipxeTemplate is returned as response to `chain` request from the bootFile/bootTemplate to boot actual OS (or Sidero agent).
var ipxeTemplate = template.Must(template.New("iPXE config").Parse(`#!ipxe
kernel /env/{{ .Env.Name }}/{{ .KernelAsset }} {{range $arg := .Env.Spec.Kernel.Args}} {{$arg}}{{end}}
initrd /env/{{ .Env.Name }}/{{ .InitrdAsset }}
boot
`))

// ipxeBootFromDiskExit script is used to skip PXE booting and boot from disk via exit.
const ipxeBootFromDiskExit = `#!ipxe
exit
`

// ipxeBootFromDiskSanboot script is used to skip PXE booting and boot from disk via sanboot.
const ipxeBootFromDiskSanboot = `#!ipxe
sanboot --no-describe --drive 0x80
`

// BootFromDisk defines a way to boot from disk.
type BootFromDisk string

const (
	BootIPXEExit BootFromDisk = "ipxe-exit"    // Use iPXE script with `exit` command.
	Boot404      BootFromDisk = "http-404"     // Return HTTP 404 response to iPXE.
	BootSANDisk  BootFromDisk = "ipxe-sanboot" // Use iPXE script with `sanboot` command.
)

var (
	apiEndpoint               string
	apiPort                   int
	extraAgentKernelArgs      string
	defaultBootFromDiskMethod BootFromDisk
	c                         client.Client
)

func bootFileHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, bootFile)
}

//nolint:unparam
func bootFromDiskHandler(method BootFromDisk, w http.ResponseWriter, r *http.Request) {
	switch method { //nolint:exhaustive
	case Boot404:
		w.WriteHeader(http.StatusNotFound)
	case BootSANDisk:
		fmt.Fprint(w, ipxeBootFromDiskSanboot)
	case BootIPXEExit:
		fallthrough
	default:
		fmt.Fprint(w, ipxeBootFromDiskExit)
	}
}

func ipxeHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	labels := labelsFromRequest(r)

	uuid := labels["uuid"]

	var arch string

	switch labels["arch"] {
	case "arm64":
		arch = "arm64"
	default:
		arch = "amd64"
	}

	server, serverBinding, err := lookupServer(uuid)
	if err != nil {
		log.Printf("Error looking up server: %v", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	env, err := newEnvironment(server, serverBinding, arch)
	if err != nil {
		if errors.Is(err, ErrBootFromDisk) {
			log.Printf("Server %q booting from disk", uuid)
			bootFromDiskHandler(defaultBootFromDiskMethod, w, r)

			return
		}

		if apierrors.IsNotFound(err) {
			log.Printf("Environment not found: %v", err)
			w.WriteHeader(http.StatusNotFound)

			return
		}

		log.Printf("%v", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	if server != nil {
		log.Printf("Using %q environment for %q", env.Name, server.Name)
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

	if !strings.HasPrefix(env.ObjectMeta.Name, "agent") {
		if err = markAsPXEBooted(server); err != nil {
			log.Printf("error marking server as PXE booted: %s", err)
		}
	}
}

func RegisterIPXE(mux *http.ServeMux, endpoint string, port int, args string, bootMethod BootFromDisk, iPXEPort int, mgrClient client.Client) error {
	apiEndpoint = endpoint
	apiPort = port
	extraAgentKernelArgs = args
	defaultBootFromDiskMethod = bootMethod
	c = mgrClient

	var embeddedScriptBuf bytes.Buffer

	if err := BootTemplate.Execute(&embeddedScriptBuf, map[string]string{
		"Endpoint": apiEndpoint,
		"Port":     strconv.Itoa(iPXEPort),
	}); err != nil {
		return err
	}

	if err := PatchBinaries(embeddedScriptBuf.Bytes()); err != nil {
		return err
	}

	mux.Handle("/boot.ipxe", logRequest(http.HandlerFunc(bootFileHandler)))
	mux.Handle("/ipxe", logRequest(http.HandlerFunc(ipxeHandler)))
	mux.Handle("/env/", logRequest(http.StripPrefix("/env/", http.FileServer(http.Dir("/var/lib/sidero/env")))))
	mux.Handle("/tftp/", logRequest(http.StripPrefix("/tftp/", http.FileServer(http.Dir("/var/lib/sidero/tftp")))))

	return nil
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

func lookupServer(uuid string) (*metalv1alpha1.Server, *infrav1.ServerBinding, error) {
	key := client.ObjectKey{
		Name: uuid,
	}

	s := &metalv1alpha1.Server{}

	if err := c.Get(context.Background(), key, s); err != nil {
		if apierrors.IsNotFound(err) {
			return nil, nil, nil
		}

		return nil, nil, err
	}

	b := &infrav1.ServerBinding{}

	if err := c.Get(context.Background(), key, b); err != nil {
		if apierrors.IsNotFound(err) {
			return s, nil, nil
		}

		return nil, nil, err
	}

	return s, b, nil
}

// newEnvironment handles which env CRD we'll respect for a given server.
// specied in the server spec overrides everything, specified in the server class overrides default, default is default :).
func newEnvironment(server *metalv1alpha1.Server, serverBinding *infrav1.ServerBinding, arch string) (env *metalv1alpha1.Environment, err error) {
	// NB: The order of this switch statement is important. It defines the
	// precedence of which environment to boot.
	switch {
	case server == nil:
		return newAgentEnvironment(arch), nil
	case serverBinding == nil:
		return newAgentEnvironment(arch), nil
	case conditions.Has(server, metalv1alpha1.ConditionPXEBooted) && !server.Spec.PXEBootAlways:
		return nil, ErrBootFromDisk
	case server.Spec.EnvironmentRef != nil:
		env, err = newEnvironmentFromServer(server)
		if err != nil {
			return nil, err
		}
	case serverBinding.Spec.ServerClassRef != nil:
		env, err = newEnvironmentFromServerClass(serverBinding)
		if err != nil {
			return nil, err
		}
	}

	if env == nil {
		env, err = newDefaultEnvironment()
		if err != nil {
			return nil, err
		}
	}

	if env == nil {
		return nil, fmt.Errorf("could not find environment for %q", server.Name)
	}

	return env, nil
}

func newAgentEnvironment(arch string) *metalv1alpha1.Environment {
	args := []string{
		"console=tty0",
		"console=ttyS0",
		"ima_appraise=fix",
		"ima_hash=sha512",
		"ima_template=ima-ng",
		"initrd=initramfs.xz",
		"ip=dhcp",
		"page_poison=1",
		"panic=30",
		"printk.devkmsg=on",
		"pti=on",
		"random.trust_cpu=on",
		"slab_nomerge=",
		"slub_debug=P",
		fmt.Sprintf("%s=%s:%d", constants.AgentEndpointArg, apiEndpoint, apiPort),
	}

	cmdline := procfs.NewCmdline(strings.Join(args, " "))
	extra := procfs.NewCmdline(extraAgentKernelArgs)

	// override defaults with extra kernel agent params
	for _, p := range extra.Parameters {
		cmdline.Set(p.Key(), p)
	}

	env := &metalv1alpha1.Environment{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("agent-%s", arch),
		},
		Spec: metalv1alpha1.EnvironmentSpec{
			Kernel: metalv1alpha1.Kernel{
				Args: cmdline.Strings(),
			},
		},
	}

	return env
}

func newDefaultEnvironment() (env *metalv1alpha1.Environment, err error) {
	env = &metalv1alpha1.Environment{}

	if err := c.Get(context.Background(), types.NamespacedName{Namespace: "", Name: metalv1alpha1.EnvironmentDefault}, env); err != nil {
		return nil, err
	}

	appendTalosConfigArgument(env)

	return env, nil
}

func newEnvironmentFromServer(server *metalv1alpha1.Server) (env *metalv1alpha1.Environment, err error) {
	env = &metalv1alpha1.Environment{}

	if err := c.Get(context.Background(), types.NamespacedName{Namespace: "", Name: server.Spec.EnvironmentRef.Name}, env); err != nil {
		return nil, err
	}

	appendTalosConfigArgument(env)

	return env, nil
}

func newEnvironmentFromServerClass(serverBinding *infrav1.ServerBinding) (env *metalv1alpha1.Environment, err error) {
	serverClassResource := &metalv1alpha1.ServerClass{}

	if err := c.Get(context.Background(), types.NamespacedName{Namespace: serverBinding.Spec.ServerClassRef.Namespace, Name: serverBinding.Spec.ServerClassRef.Name}, serverClassResource); err != nil {
		return nil, err
	}

	if serverClassResource.Spec.EnvironmentRef == nil {
		return env, nil
	}

	env = &metalv1alpha1.Environment{}

	if err := c.Get(context.Background(), types.NamespacedName{Namespace: "", Name: serverClassResource.Spec.EnvironmentRef.Name}, env); err != nil {
		return nil, err
	}

	appendTalosConfigArgument(env)

	return env, nil
}

func appendTalosConfigArgument(env *metalv1alpha1.Environment) {
	args := env.Spec.Kernel.Args

	talosConfigPrefix := talosconstants.KernelParamConfig + "="

	for _, arg := range args {
		if strings.HasPrefix(arg, talosConfigPrefix) {
			// Environment already has talos.config
			return
		}
	}

	// patch environment with the link to the metadata server
	env.Spec.Kernel.Args = append(env.Spec.Kernel.Args,
		fmt.Sprintf("%s=http://%s:%d/configdata?uuid=", talosconstants.KernelParamConfig, apiEndpoint, apiPort),
	)
}

func markAsPXEBooted(server *metalv1alpha1.Server) error {
	patchHelper, err := patch.NewHelper(server, c)
	if err != nil {
		return err
	}

	conditions.MarkTrue(server, metalv1alpha1.ConditionPXEBooted)

	return patchHelper.Patch(context.Background(), server, patch.WithOwnedConditions{
		Conditions: []clusterv1.ConditionType{metalv1alpha1.ConditionPXEBooted},
	})
}
