module github.com/talos-systems/sidero/sfyra

go 1.16

replace (
	github.com/talos-systems/sidero => ../

	// See https://github.com/talos-systems/go-loadbalancer/pull/4
	// `go get github.com/smira/tcpproxy@combined-fixes`, then copy pseudo-version there
	inet.af/tcpproxy => github.com/smira/tcpproxy v0.0.0-20201015133617-de5f7797b95b
)

require (
	github.com/spf13/cobra v1.2.1
	github.com/stretchr/testify v1.7.0
	github.com/talos-systems/cluster-api-bootstrap-provider-talos v0.4.0
	github.com/talos-systems/cluster-api-control-plane-provider-talos v0.3.0
	github.com/talos-systems/go-debug v0.2.1
	github.com/talos-systems/go-loadbalancer v0.1.1
	github.com/talos-systems/go-procfs v0.1.0
	github.com/talos-systems/go-retry v0.3.1
	github.com/talos-systems/net v0.3.0
	github.com/talos-systems/sidero v0.0.0-00010101000000-000000000000
	github.com/talos-systems/talos v0.13.0
	github.com/talos-systems/talos/pkg/machinery v0.13.1
	google.golang.org/grpc v1.41.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	k8s.io/api v0.22.2
	k8s.io/apiextensions-apiserver v0.22.2
	k8s.io/apimachinery v0.22.2
	k8s.io/client-go v0.22.2
	sigs.k8s.io/cluster-api v0.4.3
	sigs.k8s.io/controller-runtime v0.9.7
)
