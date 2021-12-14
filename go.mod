module github.com/talos-systems/sidero

go 1.16

replace github.com/pensando/goipmi v0.0.0-20200303170213-e858ec1cf0b5 => github.com/talos-systems/goipmi v0.0.0-20211214143420-35f956689e67

require (
	github.com/evanphx/json-patch v4.11.0+incompatible
	github.com/ghodss/yaml v1.0.0
	github.com/go-logr/logr v0.4.0
	github.com/hashicorp/go-multierror v1.1.1
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.16.0
	github.com/pensando/goipmi v0.0.0-20200303170213-e858ec1cf0b5
	github.com/pin/tftp v2.1.1-0.20200117065540-2f79be2dba4e+incompatible
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.7.0
	github.com/talos-systems/cluster-api-bootstrap-provider-talos v0.4.0
	github.com/talos-systems/cluster-api-control-plane-provider-talos v0.3.0
	github.com/talos-systems/go-blockdevice v0.2.4
	github.com/talos-systems/go-debug v0.2.1
	github.com/talos-systems/go-kmsg v0.1.1
	github.com/talos-systems/go-procfs v0.1.0
	github.com/talos-systems/go-retry v0.3.1
	github.com/talos-systems/go-smbios v0.1.1
	github.com/talos-systems/net v0.3.1
	github.com/talos-systems/talos/pkg/machinery v0.13.0
	golang.org/x/net v0.0.0-20210716203947-853a461950ff
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/sys v0.0.0-20210927094055-39ccf1dd6fa6
	google.golang.org/grpc v1.41.0
	google.golang.org/protobuf v1.27.1
	k8s.io/api v0.22.2
	k8s.io/apiextensions-apiserver v0.22.2
	k8s.io/apimachinery v0.22.2
	k8s.io/client-go v0.22.2
	k8s.io/utils v0.0.0-20210819203725-bdf08cb9a70a
	sigs.k8s.io/cluster-api v0.4.3
	sigs.k8s.io/controller-runtime v0.9.7
)
