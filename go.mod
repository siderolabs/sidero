module github.com/talos-systems/sidero

go 1.16

replace (
	github.com/pensando/goipmi v0.0.0-20200303170213-e858ec1cf0b5 => github.com/talos-systems/goipmi v0.0.0-20210504182258-b54796c8d678
	sigs.k8s.io/cluster-api v0.3.12 => sigs.k8s.io/cluster-api v0.3.9
)

require (
	github.com/evanphx/json-patch v4.9.0+incompatible
	github.com/ghodss/yaml v1.0.0
	github.com/go-logr/logr v0.4.0
	github.com/go-logr/zapr v0.2.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1
	github.com/onsi/ginkgo v1.16.2
	github.com/onsi/gomega v1.12.0
	github.com/opencontainers/runtime-spec v1.0.3-0.20200728170252-4d89ac9fbff6 // indirect
	github.com/pensando/goipmi v0.0.0-20200303170213-e858ec1cf0b5
	github.com/pin/tftp v2.1.1-0.20200117065540-2f79be2dba4e+incompatible
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.7.0
	github.com/talos-systems/cluster-api-bootstrap-provider-talos v0.2.0-alpha.11
	github.com/talos-systems/cluster-api-control-plane-provider-talos v0.1.0-alpha.11
	github.com/talos-systems/go-blockdevice v0.1.1-0.20201218174450-f2728a581972
	github.com/talos-systems/go-debug v0.1.0
	github.com/talos-systems/go-kmsg v0.1.0
	github.com/talos-systems/go-procfs v0.0.0-20210108152626-8cbc42d3dc24
	github.com/talos-systems/go-retry v0.2.0
	github.com/talos-systems/go-smbios v0.0.0-20210422124317-d3a32bea731a
	github.com/talos-systems/net v0.2.1-0.20210212213224-05190541b0fa
	github.com/talos-systems/talos/pkg/machinery v0.0.0-20210416105550-2b83440d6f7f // v0.9.3
	go.uber.org/zap v1.16.0 // indirect
	golang.org/x/crypto v0.0.0-20201124201722-c8d3bf9c5392 // indirect
	golang.org/x/mod v0.4.0 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/sys v0.0.0-20210503173754-0981d6026fa6
	golang.org/x/tools v0.0.0-20210101214203-2dba1e4ea05c // indirect
	google.golang.org/grpc v1.37.0
	google.golang.org/protobuf v1.26.0
	k8s.io/api v0.19.3
	k8s.io/apiextensions-apiserver v0.19.1
	k8s.io/apimachinery v0.19.3
	k8s.io/apiserver v0.19.3 // indirect
	k8s.io/client-go v0.19.3
	k8s.io/utils v0.0.0-20201110183641-67b214c5f920
	sigs.k8s.io/cluster-api v0.3.12
	sigs.k8s.io/controller-runtime v0.6.3
)
