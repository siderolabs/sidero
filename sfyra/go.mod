module github.com/talos-systems/sidero/sfyra

go 1.16

replace github.com/talos-systems/sidero => ../

// See https://github.com/talos-systems/go-loadbalancer/pull/4
replace inet.af/tcpproxy => github.com/smira/tcpproxy v0.0.0-20201015133617-de5f7797b95b

require (
	github.com/spf13/cobra v1.1.3
	github.com/stretchr/testify v1.7.0
	github.com/talos-systems/cluster-api-bootstrap-provider-talos v0.2.0-alpha.11
	github.com/talos-systems/cluster-api-control-plane-provider-talos v0.1.0-alpha.11
	github.com/talos-systems/go-loadbalancer v0.1.1-0.20201015151439-a4457024d518
	github.com/talos-systems/go-procfs v0.0.0-20210108152626-8cbc42d3dc24
	github.com/talos-systems/go-retry v0.2.1-0.20210119124456-b9dc1a990133
	github.com/talos-systems/net v0.2.1-0.20210212213224-05190541b0fa
	github.com/talos-systems/sidero v0.1.0-alpha.1.0.20200915181156-11a0a80e3d8b
	github.com/talos-systems/talos v0.9.1
	github.com/talos-systems/talos/pkg/machinery v0.0.0-20210401163915-1d8e9674a91b
	google.golang.org/grpc v1.36.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	k8s.io/api v0.20.5
	k8s.io/apiextensions-apiserver v0.19.1
	k8s.io/apimachinery v0.20.5
	k8s.io/client-go v0.20.5
	sigs.k8s.io/cluster-api v0.3.12
	sigs.k8s.io/controller-runtime v0.6.3
)

replace sigs.k8s.io/cluster-api v0.3.12 => sigs.k8s.io/cluster-api v0.3.9
