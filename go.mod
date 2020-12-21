module github.com/talos-systems/sidero

go 1.14

require (
	github.com/containerd/containerd v1.4.1 // indirect
	github.com/containerd/go-cni v1.0.1 // indirect
	github.com/evanphx/json-patch v4.9.0+incompatible
	github.com/ghodss/yaml v1.0.0
	github.com/go-logr/logr v0.2.1-0.20200730175230-ee2de8da5be6
	github.com/go-logr/zapr v0.2.0 // indirect
	github.com/golang/protobuf v1.4.3
	github.com/google/go-cmp v0.5.2 // indirect
	github.com/hashicorp/go-multierror v1.1.0
	github.com/onsi/ginkgo v1.12.1
	github.com/onsi/gomega v1.10.1
	github.com/opencontainers/runtime-spec v1.0.3-0.20200728170252-4d89ac9fbff6 // indirect
	github.com/pensando/goipmi v0.0.0-20200303170213-e858ec1cf0b5
	github.com/pin/tftp v2.1.1-0.20200117065540-2f79be2dba4e+incompatible
	github.com/pkg/errors v0.9.1
	github.com/prometheus/procfs v0.2.0 // indirect
	github.com/talos-systems/cluster-api-bootstrap-provider-talos v0.2.0-alpha.6
	github.com/talos-systems/cluster-api-control-plane-provider-talos v0.1.0-alpha.8
	github.com/talos-systems/go-blockdevice v0.1.1-0.20201111103554-874213371a3f
	github.com/talos-systems/go-procfs v0.0.0-20201215195843-16ce2ef52acd
	github.com/talos-systems/go-retry v0.1.1-0.20200922131245-752f081252cf
	github.com/talos-systems/go-smbios v0.0.0-20200807005123-80196199691e
	github.com/talos-systems/net v0.2.0
	github.com/talos-systems/talos/pkg/machinery v0.0.0-20201218234035-9aa0d3414f7a
	go.uber.org/zap v1.14.1 // indirect
	golang.org/x/crypto v0.0.0-20201016220609-9e8e0b390897 // indirect
	golang.org/x/net v0.0.0-20201110031124-69a78807bb2b // indirect
	golang.org/x/sync v0.0.0-20201008141435-b3e1573b7520
	golang.org/x/sys v0.0.0-20201018230417-eeed37f84f13
	golang.org/x/time v0.0.0-20200630173020-3af7569d3a1e // indirect
	golang.org/x/tools v0.0.0-20201019175715-b894a3290fff // indirect
	google.golang.org/grpc v1.29.1
	k8s.io/api v0.19.3
	k8s.io/apiextensions-apiserver v0.19.1
	k8s.io/apimachinery v0.19.3
	k8s.io/apiserver v0.19.3 // indirect
	k8s.io/client-go v0.19.3
	k8s.io/utils v0.0.0-20200729134348-d5654de09c73
	sigs.k8s.io/cluster-api v0.3.9
	sigs.k8s.io/controller-runtime v0.6.3
)
