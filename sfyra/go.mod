module github.com/talos-systems/sidero/sfyra

go 1.17

replace (
	// local replace
	github.com/talos-systems/sidero => ../

	// See https://github.com/talos-systems/go-loadbalancer/pull/4
	// `go get github.com/smira/tcpproxy@combined-fixes`, then copy pseudo-version there
	inet.af/tcpproxy => github.com/smira/tcpproxy v0.0.0-20201015133617-de5f7797b95b
)

require (
	github.com/spf13/cobra v1.3.0
	github.com/stretchr/testify v1.7.0
	github.com/talos-systems/cluster-api-bootstrap-provider-talos v0.5.1
	github.com/talos-systems/cluster-api-control-plane-provider-talos v0.4.3
	github.com/talos-systems/go-debug v0.2.1
	github.com/talos-systems/go-loadbalancer v0.1.1
	github.com/talos-systems/go-procfs v0.1.0
	github.com/talos-systems/go-retry v0.3.1
	github.com/talos-systems/net v0.3.1
	github.com/talos-systems/sidero v0.0.0-00010101000000-000000000000
	github.com/talos-systems/talos v0.14.1
	github.com/talos-systems/talos/pkg/machinery v0.14.1
	google.golang.org/grpc v1.44.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	k8s.io/api v0.23.1
	k8s.io/apiextensions-apiserver v0.23.1
	k8s.io/apimachinery v0.23.1
	k8s.io/client-go v0.23.1
	sigs.k8s.io/cluster-api v1.1.0
	sigs.k8s.io/controller-runtime v0.11.0
)

require (
	cloud.google.com/go v0.99.0 // indirect
	cloud.google.com/go/storage v1.10.0 // indirect
	github.com/AlekSi/pointer v1.2.0 // indirect
	github.com/BurntSushi/toml v0.4.1 // indirect
	github.com/MakeNowJust/heredoc v1.0.0 // indirect
	github.com/PuerkitoBio/purell v1.1.1 // indirect
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/antlr/antlr4/runtime/Go/antlr v0.0.0-20210826220005-b48c857c3a0e // indirect
	github.com/armon/circbuf v0.0.0-20190214190532-5111143e8da2 // indirect
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d // indirect
	github.com/aws/aws-sdk-go v1.15.78 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bgentry/go-netrc v0.0.0-20140422174119-9fd32a8b3d3d // indirect
	github.com/blang/semver v3.5.1+incompatible // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/containerd/cri v1.19.0 // indirect
	github.com/containerd/go-cni v1.1.0 // indirect
	github.com/containernetworking/cni v1.0.1 // indirect
	github.com/containernetworking/plugins v1.0.1 // indirect
	github.com/coreos/go-iptables v0.6.0 // indirect
	github.com/cosi-project/runtime v0.0.0-20211216175730-264f8fcd1a4f // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/drone/envsubst/v2 v2.0.0-20210730161058-179042472c46 // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/emicklei/dot v0.16.0 // indirect
	github.com/evanphx/json-patch v5.6.0+incompatible // indirect
	github.com/evanphx/json-patch/v5 v5.6.0 // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/gertd/go-pluralize v0.1.7 // indirect
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/go-logr/logr v1.2.0 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.19.5 // indirect
	github.com/go-openapi/swag v0.19.14 // indirect
	github.com/gobuffalo/flect v0.2.4 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/cel-go v0.9.0 // indirect
	github.com/google/go-cmp v0.5.6 // indirect
	github.com/google/go-github/v33 v33.0.0 // indirect
	github.com/google/go-querystring v1.0.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/googleapis/gax-go/v2 v2.1.1 // indirect
	github.com/googleapis/gnostic v0.5.5 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-getter v1.5.9 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-safetemp v1.0.0 // indirect
	github.com/hashicorp/go-version v1.2.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/insomniacslk/dhcp v0.0.0-20211214070828-5297eed8f489 // indirect
	github.com/jmespath/go-jmespath v0.0.0-20160803190731-bd40a432e4c7 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/josharian/native v0.0.0-20200817173448-b6b71def0850 // indirect
	github.com/jsimonetti/rtnetlink v0.0.0-20211203074127-fd9a11f42291 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.11.13 // indirect
	github.com/magiconair/properties v1.8.5 // indirect
	github.com/mailru/easyjson v0.7.6 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.2-0.20181231171920-c182affec369 // indirect
	github.com/mdlayher/ethtool v0.0.0-20211028163843-288d040e9d60 // indirect
	github.com/mdlayher/genetlink v1.1.0 // indirect
	github.com/mdlayher/netlink v1.4.2 // indirect
	github.com/mdlayher/socket v0.0.0-20211102153432-57e3fa563ecb // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/go-testing-interface v1.0.0 // indirect
	github.com/mitchellh/mapstructure v1.4.3 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/onsi/gomega v1.18.1 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/runtime-spec v1.0.3-0.20210326190908-1c3f411f0417 // indirect
	github.com/pelletier/go-toml v1.9.4 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_golang v1.11.0 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.28.0 // indirect
	github.com/prometheus/procfs v0.7.3 // indirect
	github.com/ryanuber/go-glob v1.0.0 // indirect
	github.com/spf13/afero v1.6.0 // indirect
	github.com/spf13/cast v1.4.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.10.0 // indirect
	github.com/stoewer/go-strcase v1.2.0 // indirect
	github.com/subosito/gotenv v1.2.0 // indirect
	github.com/talos-systems/crypto v0.3.4 // indirect
	github.com/talos-systems/go-blockdevice v0.2.5 // indirect
	github.com/talos-systems/go-cmd v0.1.0 // indirect
	github.com/u-root/uio v0.0.0-20210528114334-82958018845c // indirect
	github.com/ulikunitz/xz v0.5.8 // indirect
	github.com/valyala/fastjson v1.6.3 // indirect
	go.opencensus.io v0.23.0 // indirect
	golang.org/x/crypto v0.0.0-20211202192323-5770296d904e // indirect
	golang.org/x/mod v0.5.1 // indirect
	golang.org/x/net v0.0.0-20220127200216-cd36cc0744dd // indirect
	golang.org/x/oauth2 v0.0.0-20211104180415-d3ed0bb246c8 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.0.0-20220128215802-99c3d69c2c27 // indirect
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/time v0.0.0-20211116232009-f0f3c7e86c11 // indirect
	golang.org/x/tools v0.1.8 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	gomodules.xyz/jsonpatch/v2 v2.2.0 // indirect
	google.golang.org/api v0.62.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20211208223120-3a66f561d7aa // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/ini.v1 v1.66.2 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	honnef.co/go/tools v0.2.2 // indirect
	inet.af/tcpproxy v0.0.0-20200125044825-b6bb9b5b8252 // indirect
	k8s.io/apiserver v0.23.1 // indirect
	k8s.io/cluster-bootstrap v0.23.0 // indirect
	k8s.io/component-base v0.23.1 // indirect
	k8s.io/klog/v2 v2.30.0 // indirect
	k8s.io/kube-openapi v0.0.0-20211115234752-e816edb12b65 // indirect
	k8s.io/utils v0.0.0-20210930125809-cb0fa318a74b // indirect
	sigs.k8s.io/json v0.0.0-20211020170558-c049b76a60c6 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.0 // indirect
	sigs.k8s.io/yaml v1.3.0 // indirect
)
