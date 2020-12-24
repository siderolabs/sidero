module github.com/talos-systems/sidero/sfyra

go 1.14

replace github.com/talos-systems/sidero => ../

// See https://github.com/talos-systems/go-loadbalancer/pull/4
replace inet.af/tcpproxy => github.com/smira/tcpproxy v0.0.0-20201015133617-de5f7797b95b

require (
	github.com/spf13/cobra v1.1.1
	github.com/stretchr/testify v1.6.1
	github.com/talos-systems/cluster-api-bootstrap-provider-talos v0.2.0-alpha.6
	github.com/talos-systems/cluster-api-control-plane-provider-talos v0.1.0-alpha.8
	github.com/talos-systems/go-loadbalancer v0.1.1-0.20201015151439-a4457024d518
	github.com/talos-systems/go-procfs v0.0.0-20201223150035-a82654edcec1
	github.com/talos-systems/go-retry v0.1.1-0.20201113203059-8c63d290a688
	github.com/talos-systems/net v0.2.0
	github.com/talos-systems/sidero v0.1.0-alpha.1.0.20200915181156-11a0a80e3d8b
	github.com/talos-systems/talos v0.8.0
	github.com/talos-systems/talos/pkg/machinery v0.0.0-20201223224110-cf5226a5ab55
	google.golang.org/grpc v1.29.1
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776
	k8s.io/api v0.20.1
	k8s.io/apiextensions-apiserver v0.19.1
	k8s.io/apimachinery v0.20.1
	k8s.io/client-go v0.20.1
	sigs.k8s.io/cluster-api v0.3.9
	sigs.k8s.io/controller-runtime v0.6.3
)
