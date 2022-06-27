module github.com/talos-systems/sidero

go 1.17

replace github.com/pensando/goipmi v0.0.0-20200303170213-e858ec1cf0b5 => github.com/talos-systems/goipmi v0.0.0-20211214143420-35f956689e67

require (
	github.com/evanphx/json-patch v5.6.0+incompatible
	github.com/ghodss/yaml v1.0.0
	github.com/go-logr/logr v1.2.3
	github.com/google/go-cmp v0.5.8
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/hashicorp/go-multierror v1.1.1
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.19.0
	github.com/pensando/goipmi v0.0.0-20200303170213-e858ec1cf0b5
	github.com/pin/tftp v2.1.1-0.20200117065540-2f79be2dba4e+incompatible
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.7.5
	github.com/talos-systems/go-blockdevice v0.3.1
	github.com/talos-systems/go-debug v0.2.1
	github.com/talos-systems/go-kmsg v0.1.1
	github.com/talos-systems/go-procfs v0.1.0
	github.com/talos-systems/go-retry v0.3.1
	github.com/talos-systems/go-smbios v0.1.2-0.20220202124743-3f1e775b7e3e
	github.com/talos-systems/grpc-proxy v0.3.0
	github.com/talos-systems/net v0.3.2
	github.com/talos-systems/siderolink v0.1.1-0.20211130121818-9902ad2774f0
	github.com/talos-systems/talos/pkg/machinery v1.0.4
	go.uber.org/zap v1.21.0
	golang.org/x/net v0.0.0-20220425223048-2871e0cb64e4
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/sys v0.0.0-20220502124256-b6088ccd6cba
	golang.zx2c4.com/wireguard/wgctrl v0.0.0-20220330030906-9490840b0b01
	google.golang.org/grpc v1.45.0
	google.golang.org/protobuf v1.28.0
	inet.af/netaddr v0.0.0-20211027220019-c74959edd3b6
	k8s.io/api v0.23.6
	k8s.io/apiextensions-apiserver v0.23.6
	k8s.io/apimachinery v0.23.6
	k8s.io/client-go v0.23.6
	k8s.io/utils v0.0.0-20220210201930-3a6ce19ff2f9
	sigs.k8s.io/cluster-api v1.1.3
	sigs.k8s.io/controller-runtime v0.11.2
)

require (
	cloud.google.com/go v0.93.3 // indirect
	github.com/AlekSi/pointer v1.2.0 // indirect
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/blang/semver v3.5.1+incompatible // indirect
	github.com/cespare/xxhash/v2 v2.1.1 // indirect
	github.com/containerd/go-cni v1.1.3 // indirect
	github.com/containernetworking/cni v1.0.1 // indirect
	github.com/cosi-project/runtime v0.0.0-20211216175730-264f8fcd1a4f // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/digitalocean/go-smbios v0.0.0-20180907143718-390a4f403a8e // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/gertd/go-pluralize v0.1.7 // indirect
	github.com/go-logr/zapr v1.2.0 // indirect
	github.com/gobuffalo/flect v0.2.4 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/googleapis/gnostic v0.5.5 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/josharian/native v1.0.0 // indirect
	github.com/jsimonetti/rtnetlink v1.1.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.2-0.20181231171920-c182affec369 // indirect
	github.com/mdlayher/ethtool v0.0.0-20220213132912-856bd6cb8a38 // indirect
	github.com/mdlayher/genetlink v1.2.0 // indirect
	github.com/mdlayher/netlink v1.6.0 // indirect
	github.com/mdlayher/socket v0.2.3 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/nxadm/tail v1.4.8 // indirect
	github.com/opencontainers/runtime-spec v1.0.3-0.20200929063507-e6143ca7d51d // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_golang v1.11.0 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.28.0 // indirect
	github.com/prometheus/procfs v0.7.2 // indirect
	github.com/ryanuber/go-glob v1.0.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/talos-systems/crypto v0.3.5 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	go4.org/intern v0.0.0-20211027215823-ae77deb06f29 // indirect
	go4.org/unsafe/assume-no-moving-gc v0.0.0-20211027215541-db492cf91b37 // indirect
	golang.org/x/crypto v0.0.0-20220411220226-7b82a4e95df4 // indirect
	golang.org/x/oauth2 v0.0.0-20211104180415-d3ed0bb246c8 // indirect
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/time v0.0.0-20210723032227-1f47c861a9ac // indirect
	golang.zx2c4.com/go118/netip v0.0.0-20211111135330-a4a02eeacf9d // indirect
	golang.zx2c4.com/wintun v0.0.0-20211104114900-415007cec224 // indirect
	golang.zx2c4.com/wireguard v0.0.0-20220202223031-3b95c81cc178 // indirect
	gomodules.xyz/jsonpatch/v2 v2.2.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20220228155957-1da8797a5878 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/cluster-bootstrap v0.23.0 // indirect
	k8s.io/component-base v0.23.6 // indirect
	k8s.io/klog/v2 v2.30.0 // indirect
	k8s.io/kube-openapi v0.0.0-20211115234752-e816edb12b65 // indirect
	sigs.k8s.io/json v0.0.0-20211020170558-c049b76a60c6 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.1 // indirect
	sigs.k8s.io/yaml v1.3.0 // indirect
)
