module github.com/talos-systems/sidero/sfyra

go 1.14

replace github.com/talos-systems/sidero => ../

require (
	github.com/go-logr/logr v0.2.1-0.20200730175230-ee2de8da5be6 // indirect
	github.com/spf13/cobra v1.0.0
	github.com/stretchr/testify v1.6.1
	github.com/talos-systems/cluster-api-bootstrap-provider-talos v0.2.0-alpha.3
	github.com/talos-systems/cluster-api-control-plane-provider-talos v0.1.0-alpha.4
	github.com/talos-systems/go-loadbalancer v0.1.1-0.20200918172657-3c8f3471d14e
	github.com/talos-systems/go-procfs v0.0.0-20200219015357-57c7311fdd45
	github.com/talos-systems/go-retry v0.1.0
	github.com/talos-systems/net v0.2.0
	github.com/talos-systems/sidero v0.1.0-alpha.1.0.20200915181156-11a0a80e3d8b
	github.com/talos-systems/talos v0.7.0-alpha.1.0.20200916165852-41ecb826469a
	github.com/talos-systems/talos/pkg/machinery v0.0.0-20200916165852-41ecb826469a
	google.golang.org/grpc v1.29.1
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776
	k8s.io/api v0.19.1
	k8s.io/apiextensions-apiserver v0.19.1
	k8s.io/apimachinery v0.19.1
	k8s.io/client-go v0.19.1
	sigs.k8s.io/cluster-api v0.3.9
	sigs.k8s.io/controller-runtime v0.6.3
)
