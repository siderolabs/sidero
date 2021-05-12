module github.com/talos-systems/sidero/sfyra

go 1.16

replace (
	github.com/talos-systems/sidero => ../

	// See https://github.com/talos-systems/go-loadbalancer/pull/4
	// `go get github.com/smira/tcpproxy@combined-fixes`, then copy pseudo-version there
	inet.af/tcpproxy => github.com/smira/tcpproxy v0.0.0-20201015133617-de5f7797b95b

	sigs.k8s.io/cluster-api v0.3.12 => sigs.k8s.io/cluster-api v0.3.9
)

require (
	github.com/spf13/cobra v1.1.3
	github.com/stretchr/testify v1.7.0
	github.com/talos-systems/cluster-api-bootstrap-provider-talos v0.2.0-alpha.11
	github.com/talos-systems/cluster-api-control-plane-provider-talos v0.1.0-alpha.11
	github.com/talos-systems/go-debug v0.1.0
	github.com/talos-systems/go-loadbalancer v0.1.1
	github.com/talos-systems/go-procfs v0.0.0-20210108152626-8cbc42d3dc24
	github.com/talos-systems/go-retry v0.2.1-0.20210119124456-b9dc1a990133
	github.com/talos-systems/net v0.2.1-0.20210212213224-05190541b0fa
	github.com/talos-systems/sidero v0.0.0-00010101000000-000000000000
	github.com/talos-systems/talos v0.9.3
	github.com/talos-systems/talos/pkg/machinery v0.0.0-20210416105550-2b83440d6f7f // v0.9.3
	google.golang.org/grpc v1.37.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	k8s.io/api v0.20.5
	k8s.io/apiextensions-apiserver v0.19.1
	k8s.io/apimachinery v0.20.5
	k8s.io/client-go v0.20.5
	sigs.k8s.io/cluster-api v0.3.12
	sigs.k8s.io/controller-runtime v0.6.3
)
