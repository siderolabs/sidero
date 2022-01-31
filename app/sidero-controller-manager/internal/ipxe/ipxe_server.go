// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ipxe

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/healthz"

	"github.com/talos-systems/go-procfs/procfs"
	talosconstants "github.com/talos-systems/talos/pkg/machinery/constants"
	"github.com/talos-systems/talos/pkg/machinery/kernel"

	infrav1 "github.com/talos-systems/sidero/app/caps-controller-manager/api/v1alpha3"
	metalv1 "github.com/talos-systems/sidero/app/sidero-controller-manager/api/v1alpha2"
	"github.com/talos-systems/sidero/app/sidero-controller-manager/internal/siderolink"
	"github.com/talos-systems/sidero/app/sidero-controller-manager/pkg/constants"
	siderotypes "github.com/talos-systems/sidero/app/sidero-controller-manager/pkg/types"
)

var ErrBootFromDisk = errors.New("boot from disk")

// BootTemplate is embedded into iPXE binary when that binary is sent to the node.
var BootTemplate = template.Must(template.New("iPXE embedded").Parse(`#!ipxe
prompt --key 0x02 --timeout 2000 Press Ctrl-B for the iPXE command line... && shell ||

# print interfaces
ifstat

# retry 10 times overall
set attempts:int32 10
set x:int32 0

:retry_loop

	set idx:int32 0

	:loop
		# try DHCP on each interface
		isset ${net${idx}/mac} || goto exhausted

		ifclose
		iflinkwait --timeout 5000 net${idx} || goto next_iface
		dhcp net${idx} && goto boot

	:next_iface
		inc idx && goto loop

	:boot
		# attempt boot, if fails try next iface
		route

		chain --replace http://{{ .Endpoint }}:{{ .Port }}/ipxe?uuid=${uuid}&mac=${net${idx}/mac:hexhyp}&domain=${domain}&hostname=${hostname}&serial=${serial}&arch=${buildarch} || goto next_iface

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

var (
	apiEndpoint               string
	apiPort                   int
	extraAgentKernelArgs      string
	defaultBootFromDiskMethod siderotypes.BootFromDisk
	c                         client.Client
)

func bootFileHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, embeddedScriptBuf.Bytes())
}

//nolint:unparam
func bootFromDiskHandler(method siderotypes.BootFromDisk, w http.ResponseWriter, r *http.Request) {
	switch method {
	case siderotypes.Boot404:
		w.WriteHeader(http.StatusNotFound)
	case siderotypes.BootSANDisk:
		fmt.Fprint(w, ipxeBootFromDiskSanboot)
	case siderotypes.BootIPXEExit:
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

			method, err := getBootFromDiskMethod(server, serverBinding)
			if err != nil {
				log.Printf("%v", err)
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

			bootFromDiskHandler(method, w, r)

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

	if !env.IsReady() {
		log.Printf("Environment not ready: %q", env.Name)

		w.WriteHeader(http.StatusPreconditionFailed)
		fmt.Fprintf(w, "environment %q is not ready", env.Name)

		return
	}

	if server != nil {
		log.Printf("Using %q environment for %q", env.Name, server.Name)
	} else {
		log.Printf("Using %q environment", env.Name)
	}

	args := struct {
		Env         *metalv1.Environment
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

	// This code is left here only for backward compatibility with Talos <= v0.13.
	if !strings.HasPrefix(env.ObjectMeta.Name, "agent") {
		// Do not mark as PXE booted here if SideroLink events are available and Talos installation is in progress.
		// SideroLink events handler will mark the machine with TalosInstalledCondition condition,
		// then server controller will reconcile this status and mark server as PXEBooted.
		if conditions.Has(serverBinding, infrav1.TalosInstalledCondition) {
			return
		}

		if err = markAsPXEBooted(server); err != nil {
			log.Printf("error marking server as PXE booted: %s", err)
		}
	}
}

func getBootFromDiskMethod(server *metalv1.Server, serverBinding *infrav1.ServerBinding) (siderotypes.BootFromDisk, error) {
	method := defaultBootFromDiskMethod

	if server.Spec.BootFromDiskMethod != "" {
		method = server.Spec.BootFromDiskMethod
	} else if serverBinding.Spec.ServerClassRef != nil {
		var serverClass metalv1.ServerClass

		if err := c.Get(
			context.TODO(),
			types.NamespacedName{
				Namespace: serverBinding.Spec.ServerClassRef.Namespace,
				Name:      serverBinding.Spec.ServerClassRef.Name,
			},
			&serverClass,
		); err != nil {
			return "", err
		}

		if serverClass.Spec.BootFromDiskMethod != "" {
			method = serverClass.Spec.BootFromDiskMethod
		}
	}

	return method, nil
}

var embeddedScriptBuf bytes.Buffer

func RegisterIPXE(mux *http.ServeMux, endpoint string, port int, args string, bootMethod siderotypes.BootFromDisk, iPXEPort int, mgrClient client.Client) error {
	apiEndpoint = endpoint
	apiPort = port
	extraAgentKernelArgs = args
	defaultBootFromDiskMethod = bootMethod
	c = mgrClient

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

func lookupServer(uuid string) (*metalv1.Server, *infrav1.ServerBinding, error) {
	key := client.ObjectKey{
		Name: uuid,
	}

	s := &metalv1.Server{}

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
func newEnvironment(server *metalv1.Server, serverBinding *infrav1.ServerBinding, arch string) (env *metalv1.Environment, err error) {
	// NB: The order of this switch statement is important. It defines the
	// precedence of which environment to boot.
	switch {
	case server == nil:
		return newAgentEnvironment(arch), nil
	case serverBinding == nil:
		return newAgentEnvironment(arch), nil
	case conditions.Has(server, metalv1.ConditionPXEBooted) && !server.Spec.PXEBootAlways:
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

func newAgentEnvironment(arch string) *metalv1.Environment {
	args := append([]string(nil), kernel.DefaultArgs...)
	args = append(args,
		"console=tty0",
		"console=ttyS0",
		"initrd=initramfs.xz",
		"ip=dhcp",
		"panic=30",
		fmt.Sprintf("%s=%s:%d", constants.AgentEndpointArg, apiEndpoint, apiPort),
	)

	cmdline := procfs.NewCmdline(strings.Join(args, " "))
	extra := procfs.NewCmdline(extraAgentKernelArgs)

	// override defaults with extra kernel agent params
	for _, p := range extra.Parameters {
		cmdline.Set(p.Key(), p)
	}

	env := &metalv1.Environment{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("agent-%s", arch),
		},
		Spec: metalv1.EnvironmentSpec{
			Kernel: metalv1.Kernel{
				Args: cmdline.Strings(),
			},
		},
	}

	return env
}

func newDefaultEnvironment() (env *metalv1.Environment, err error) {
	env = &metalv1.Environment{}

	if err := c.Get(context.Background(), types.NamespacedName{Namespace: "", Name: metalv1.EnvironmentDefault}, env); err != nil {
		return nil, err
	}

	appendTalosArguments(env)

	return env, nil
}

func newEnvironmentFromServer(server *metalv1.Server) (env *metalv1.Environment, err error) {
	env = &metalv1.Environment{}

	if err := c.Get(context.Background(), types.NamespacedName{Namespace: "", Name: server.Spec.EnvironmentRef.Name}, env); err != nil {
		return nil, err
	}

	appendTalosArguments(env)

	return env, nil
}

func newEnvironmentFromServerClass(serverBinding *infrav1.ServerBinding) (env *metalv1.Environment, err error) {
	serverClassResource := &metalv1.ServerClass{}

	if err := c.Get(context.Background(), types.NamespacedName{Namespace: serverBinding.Spec.ServerClassRef.Namespace, Name: serverBinding.Spec.ServerClassRef.Name}, serverClassResource); err != nil {
		return nil, err
	}

	if serverClassResource.Spec.EnvironmentRef == nil {
		return env, nil
	}

	env = &metalv1.Environment{}

	if err := c.Get(context.Background(), types.NamespacedName{Namespace: "", Name: serverClassResource.Spec.EnvironmentRef.Name}, env); err != nil {
		return nil, err
	}

	appendTalosArguments(env)

	return env, nil
}

func appendTalosArguments(env *metalv1.Environment) {
	args := env.Spec.Kernel.Args

	talosConfigPrefix := talosconstants.KernelParamConfig + "="
	sideroLinkPrefix := talosconstants.KernelParamSideroLink + "="
	logDeliveryPrefix := talosconstants.KernelParamLoggingKernel + "="
	eventsSinkPrefix := talosconstants.KernelParamEventsSink + "="

	for _, prefix := range []string{
		talosConfigPrefix,
		sideroLinkPrefix,
		logDeliveryPrefix,
		eventsSinkPrefix,
	} {
		for _, arg := range args {
			if strings.HasPrefix(arg, prefix) {
				// Environment already has variable, skip it
				return
			}
		}

		switch prefix {
		case talosConfigPrefix:
			// patch environment with the link to the metadata server
			env.Spec.Kernel.Args = append(env.Spec.Kernel.Args,
				fmt.Sprintf("%s=http://%s:%d/configdata?uuid=", talosconstants.KernelParamConfig, apiEndpoint, apiPort),
			)
		case sideroLinkPrefix:
			// patch environment with the SideroLink API
			env.Spec.Kernel.Args = append(env.Spec.Kernel.Args,
				fmt.Sprintf("%s=%s:%d", talosconstants.KernelParamSideroLink, apiEndpoint, apiPort),
			)
		case logDeliveryPrefix:
			// patch environment with the log receiver endpoint
			env.Spec.Kernel.Args = append(env.Spec.Kernel.Args,
				fmt.Sprintf("%s=tcp://[%s]:%d", talosconstants.KernelParamLoggingKernel, siderolink.Cfg.ServerAddress.IP(), siderolink.LogReceiverPort),
			)
		case eventsSinkPrefix:
			// patch environment with the events sink endpoint
			env.Spec.Kernel.Args = append(env.Spec.Kernel.Args,
				fmt.Sprintf("%s=[%s]:%d", talosconstants.KernelParamEventsSink, siderolink.Cfg.ServerAddress.IP(), siderolink.EventsSinkPort),
			)
		}
	}
}

func markAsPXEBooted(server *metalv1.Server) error {
	patchHelper, err := patch.NewHelper(server, c)
	if err != nil {
		return err
	}

	conditions.MarkTrue(server, metalv1.ConditionPXEBooted)

	return patchHelper.Patch(context.Background(), server, patch.WithOwnedConditions{
		Conditions: []clusterv1.ConditionType{metalv1.ConditionPXEBooted},
	})
}

func Check(addr string) healthz.Checker {
	return func(_ *http.Request) error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("http://%s/boot.ipxe", addr), nil)
		if err != nil {
			return err
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}

		defer func() {
			if resp.Body != nil {
				io.Copy(io.Discard, resp.Body) //nolint:errcheck
				resp.Body.Close()              //nolint:errcheck
			}
		}()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("unexpected code %d", resp.StatusCode)
		}

		return nil
	}
}
