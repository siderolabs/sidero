## [Sidero 0.3.0-alpha.1](https://github.com/talos-systems/sidero/releases/tag/v0.3.0-alpha.1) (2021-05-20)

Welcome to the v0.3.0-alpha.1 release of Sidero!  
*This is a pre-release of Sidero*



Please try out the release binaries and report any issues at
https://github.com/talos-systems/sidero/issues.

### ServerClass `any` and Environment `default`

Sidero now creates ServerClass `any` which matches all servers.

Environment `default` is created which will supply Talos release that Sidero was built with, as well as default kernel flags.


### Boot from Disk Methods

If the server is configured to PXE boot by default, it might hit the Sidero iPXE server after Talos install, so Sidero has to force the
server to boot from disk.

Sidero 0.2 supports default method via iPXE `exit` command, but this command doesn't always work
([details](https://ipxe.org/appnote/work_around_bios_halting_on_ipxe_exit)).

Sidero 0.3 adds support for two additional methods:

* `http-404` force HTTP 404 response from iPXE server
* `ipxe-sanboot` uses `sanboot` command to boot from the first disk

Variable `SIDERO_CONTROLLER_MANAGER_BOOT_FROM_DISK_METHOD` controls this setting.


### Install and Upgrade Variables

Some aspects of Sidero installation can now be [controlled](https://www.sidero.dev/docs/v0.3/getting-started/installation/) via environment variables or `clusterctl` variables.


### IPMI Discovery and Automatic Setup

Sidero agent by default discovers BMC (IPMI) endpoint when it first runs on the server and provisions user for Sidero automatically.
This feature enables completely automated server and IPMI discovery on the agent boot.
Variable `SIDERO_CONTROLLER_MANAGER_AUTO_BMC_SETUP` can be used to disable this feature during install or upgrade of Sidero.

Additionally Sidero supports storing IMPI credentials in Kubernetes secrets referencing them from the Server object.


### iPXE script

iPXE image delivered by Sidero (either `ipxe.efi` or `undionly.kpxe`) now contains embedded iPXE script to access Sidero iPXE server.
This change allows to simplify DHCP server setup to return only iPXE image without any additional `if`s.


### Label Selector support in ServerClass

ServerClasses now support Kubernetes label selectors: Servers can be selected based on their labels.


### `metal-controller-manager` External Ports Change

Component `metal-metadata-server` was merged into `metal-controller-manager`, and three separate HTTP endpoints were merged into one endpoint on port `:8081`:

* iPXE server (which used to be on port 8081)
* internal gRPC server (Sidero agent uses it to talk back to Sidero service), previously was 50001
* metadata server endpoint (used to be separate deployment and service, docs used port 9091)


### Support for arm64

All components are now capable of running on arm64, including Rasberry Pi 4.

Sidero can provision `amd64` and `arm64` nodes from any platform.
Only UEFI boot is supported for `arm64`.

> Note: Upstream CAPI includes `kube-rbac-proxy` component which is not compatible with `arm64` at the moment of this writing.
A workaround is to patch the deployment to bump the `kube-rbac-proxy` image from v0.4.1 to v0.8.0.


### Contributors

* Andrey Smirnov
* Alexey Palazhchenko
* Spencer Smith
* Artem Chernyshev
* Andrew Rynhard
* Brandon Nason
* Matt Zahorik
* bzub

### Changes
<details><summary>45 commits</summary>
<p>

* [`d51fda5`](https://github.com/talos-systems/sidero/commit/d51fda5cc62e6ca83445604c58150139ff492e84) release(v0.3.0-alpha.1): prepare release
* [`dcc3fde`](https://github.com/talos-systems/sidero/commit/dcc3fde21f0e931be2864ea0ecc829a258f9ba37) feat: add label selector to serverclass
* [`3caa6f5`](https://github.com/talos-systems/sidero/commit/3caa6f529895aebb32f921fe9fa9b84a550400f6) chore: fix markdown linting
* [`a792890`](https://github.com/talos-systems/sidero/commit/a792890011a67d080ede1bc8b866ca032983a421) feat: provide several options to force boot from disk via iPXE
* [`1e8096e`](https://github.com/talos-systems/sidero/commit/1e8096e9214a37a3f3197c03c63f7d24017c0cd1) docs: add Mermaid
* [`c240381`](https://github.com/talos-systems/sidero/commit/c24038157bc8ec43e479fcca8dd2abb8a576e1ac) feat: bump default Talos version to v0.10.2
* [`0a50888`](https://github.com/talos-systems/sidero/commit/0a50888d55e4f6e52afe2a046b20fe88f047d686) docs: move to standardized template
* [`4a1183b`](https://github.com/talos-systems/sidero/commit/4a1183b0834fda539f122d7ee6b2c45697cce59e) feat: ship Environment "default"
* [`2e8c1ed`](https://github.com/talos-systems/sidero/commit/2e8c1ed38a0d4ba92e87c152b3df78a6051302b8) chore: fix a few linter warnings
* [`6bcf9a8`](https://github.com/talos-systems/sidero/commit/6bcf9a8810c71027dad970a457b415f517f458a4) chore: bump dependencies via dependabot
* [`4c0b3de`](https://github.com/talos-systems/sidero/commit/4c0b3de6d92f6cbc5727bf959e6d06289c5b636d) docs: clarify docs around endpoints and metadata server
* [`990263a`](https://github.com/talos-systems/sidero/commit/990263a6b38c140e5db511b35c6eef4b7f789927) feat: merge Sidero HTTP endpoints under a single port
* [`5266a76`](https://github.com/talos-systems/sidero/commit/5266a76a56242f5ef2340780b4fcc13978bfe844) chore: reduce bootstrap node resources in Sfyra
* [`7f3b4b8`](https://github.com/talos-systems/sidero/commit/7f3b4b8e0cf21e8c38f78fb078ff646fa8f881be) fix: remove kube-rbac-proxy
* [`bfa3cd9`](https://github.com/talos-systems/sidero/commit/bfa3cd91228dccad47cf2e96dc0aed7e9dc1c48a) chore: re-sign the .drone.yml file
* [`cf383ab`](https://github.com/talos-systems/sidero/commit/cf383ab630e17990809180c01df43d91af282bd0) chore: use release-tool to generate release notes
* [`056f8c2`](https://github.com/talos-systems/sidero/commit/056f8c2942c34fb1fef58e1749a50edb74ba9921) release(v0.3.0-alpha.0): prepare release
* [`ca75eb4`](https://github.com/talos-systems/sidero/commit/ca75eb49e0cbd6a4a0608db20bed3c063ca840d9) fix: back down resource requests
* [`3a6c5b9`](https://github.com/talos-systems/sidero/commit/3a6c5b93d3df19bc8288c8ab52268b08798f5eec) chore: parse "boolean" variables
* [`148e228`](https://github.com/talos-systems/sidero/commit/148e228d404f3163e1ec57880ee9e2f218205c09) chore: bump dependencies
* [`db28ed3`](https://github.com/talos-systems/sidero/commit/db28ed3778ead07ddc04fba81d5102947220d13c) chore: allow building with debug handlers
* [`2cdde00`](https://github.com/talos-systems/sidero/commit/2cdde009ee00efb550da8a2f5a131f0a251a997f) test: fix the Environment args to support UEFI boot
* [`d15a960`](https://github.com/talos-systems/sidero/commit/d15a960f93bdb557c577bb681da10ced55e18608) chore: add golangci-lint configuration
* [`b3afd17`](https://github.com/talos-systems/sidero/commit/b3afd17a8e80955aac297f4b7dd2db8070f32ec1) feat: ship ServerClass "any"
* [`94ff33b`](https://github.com/talos-systems/sidero/commit/94ff33b4c4cdaabc2aa4e57161107a83a38ae773) feat: enable auto-setup of BMC
* [`52647f9`](https://github.com/talos-systems/sidero/commit/52647f90a2707a6a22b01bcd74b3216723d16ff0) docs: fix install disk name in the examples
* [`44eaa7d`](https://github.com/talos-systems/sidero/commit/44eaa7d7f09e45c1a70b9df49ec20bfed417823e) feat: inject iPXE script into the iPXE binaries
* [`1659b96`](https://github.com/talos-systems/sidero/commit/1659b965949ef196e30b623a93be9f4161084775) docs: fix typo
* [`fb04b24`](https://github.com/talos-systems/sidero/commit/fb04b24d7db86278e785a2f61f2aee5036def7f7) chore: fix `make help` command
* [`f8bc9b1`](https://github.com/talos-systems/sidero/commit/f8bc9b194a618a943d6d5afe3e157007c2803bb3) test: add missing empty tests
* [`b17f370`](https://github.com/talos-systems/sidero/commit/b17f37092c003505e4ef1e7e31d02ebd0772355f) test: port improvements from Talos
* [`c43f9ec`](https://github.com/talos-systems/sidero/commit/c43f9ecf67c4bbd6168a9cb6c541a6899d0bff37) test: run unit tests on CI, skip broken
* [`45cb973`](https://github.com/talos-systems/sidero/commit/45cb97356e2f4635683885e628926b5367c83437) chore: document Makefile target
* [`8e12ab9`](https://github.com/talos-systems/sidero/commit/8e12ab9aa0990c169340576be98482b7ba83248c) chore: disable dependabot rebases
* [`4267ed7`](https://github.com/talos-systems/sidero/commit/4267ed782232db447a68f664aefcc105ba0decc6) chore: update dependecies
* [`4aae477`](https://github.com/talos-systems/sidero/commit/4aae4774e781ee54f8fb2da62f5f4e7297cfe468) chore: add dependabot config
* [`1e33dcd`](https://github.com/talos-systems/sidero/commit/1e33dcd4e7cec0bc6c881eada455c4fe8b8a91d8) feat: pull the Sidero configuration as `clusterctl` variables
* [`23c75e4`](https://github.com/talos-systems/sidero/commit/23c75e437c92a7c8aa1aca419c57bf652b98cdd5) docs: add guide for upgrades
* [`c9aca48`](https://github.com/talos-systems/sidero/commit/c9aca4824d058000c7d2a523812bba61fdd0bec3) docs: fix the docs around CONTROL_PLANE_ENDPOINT
* [`37e4ab7`](https://github.com/talos-systems/sidero/commit/37e4ab7c6b4871bc19bdceda8cfc3ef59c006805) fix: remove erroneous wg.Add in environment controller
* [`74d9bf9`](https://github.com/talos-systems/sidero/commit/74d9bf9be4b3639cfa56bac8d51b64b5a821ad8b) feat: add an option to reference IPMI creds via the secret refs
* [`0613b8f`](https://github.com/talos-systems/sidero/commit/0613b8fb5c39ff194cf2a3f7d380735800f02e67) feat: pull in new version of go-smbios for UUID endiannes fix
* [`f6ca6e8`](https://github.com/talos-systems/sidero/commit/f6ca6e81bc0d1430cc81fdbd7c2a3717e9f05bb6) feat: build Sidero for amd64 and arm64, support arm64 in the agent
* [`8960225`](https://github.com/talos-systems/sidero/commit/89602254728235cdee1f93a566f356b4be07f776) chore: improve build system
* [`a20fcf9`](https://github.com/talos-systems/sidero/commit/a20fcf9fe9563dab2c03e2b21aedbff3fea6b4c5) docs: create v0.2 docs and add note about specifying TALOS_VERSION
</p>
</details>

### Changes since v0.3.0-alpha.0
<details><summary>16 commits</summary>
<p>

* [`d51fda5`](https://github.com/talos-systems/sidero/commit/d51fda5cc62e6ca83445604c58150139ff492e84) release(v0.3.0-alpha.1): prepare release
* [`dcc3fde`](https://github.com/talos-systems/sidero/commit/dcc3fde21f0e931be2864ea0ecc829a258f9ba37) feat: add label selector to serverclass
* [`3caa6f5`](https://github.com/talos-systems/sidero/commit/3caa6f529895aebb32f921fe9fa9b84a550400f6) chore: fix markdown linting
* [`a792890`](https://github.com/talos-systems/sidero/commit/a792890011a67d080ede1bc8b866ca032983a421) feat: provide several options to force boot from disk via iPXE
* [`1e8096e`](https://github.com/talos-systems/sidero/commit/1e8096e9214a37a3f3197c03c63f7d24017c0cd1) docs: add Mermaid
* [`c240381`](https://github.com/talos-systems/sidero/commit/c24038157bc8ec43e479fcca8dd2abb8a576e1ac) feat: bump default Talos version to v0.10.2
* [`0a50888`](https://github.com/talos-systems/sidero/commit/0a50888d55e4f6e52afe2a046b20fe88f047d686) docs: move to standardized template
* [`4a1183b`](https://github.com/talos-systems/sidero/commit/4a1183b0834fda539f122d7ee6b2c45697cce59e) feat: ship Environment "default"
* [`2e8c1ed`](https://github.com/talos-systems/sidero/commit/2e8c1ed38a0d4ba92e87c152b3df78a6051302b8) chore: fix a few linter warnings
* [`6bcf9a8`](https://github.com/talos-systems/sidero/commit/6bcf9a8810c71027dad970a457b415f517f458a4) chore: bump dependencies via dependabot
* [`4c0b3de`](https://github.com/talos-systems/sidero/commit/4c0b3de6d92f6cbc5727bf959e6d06289c5b636d) docs: clarify docs around endpoints and metadata server
* [`990263a`](https://github.com/talos-systems/sidero/commit/990263a6b38c140e5db511b35c6eef4b7f789927) feat: merge Sidero HTTP endpoints under a single port
* [`5266a76`](https://github.com/talos-systems/sidero/commit/5266a76a56242f5ef2340780b4fcc13978bfe844) chore: reduce bootstrap node resources in Sfyra
* [`7f3b4b8`](https://github.com/talos-systems/sidero/commit/7f3b4b8e0cf21e8c38f78fb078ff646fa8f881be) fix: remove kube-rbac-proxy
* [`bfa3cd9`](https://github.com/talos-systems/sidero/commit/bfa3cd91228dccad47cf2e96dc0aed7e9dc1c48a) chore: re-sign the .drone.yml file
* [`cf383ab`](https://github.com/talos-systems/sidero/commit/cf383ab630e17990809180c01df43d91af282bd0) chore: use release-tool to generate release notes
</p>
</details>

### Changes from talos-systems/cluster-api-bootstrap-provider-talos
<details><summary>4 commits</summary>
<p>

* [`63b7459`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/63b7459e66073eb480dd9fcd0547dd9b025d90e6) release(v0.2.0-alpha.12): prepare release
* [`f59baf5`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/f59baf5b0fa8595b1f521f78848d8df9b22e558e) fix: back down resource requests
* [`92f42c4`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/92f42c400acacd5d799e19d5013ea759e5bb7085) fix: ensure secrets are deleted when cluster is dropped
* [`2487307`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/2487307ffc274a7f46ed4e5d47c1f3d9bbb8d4d4) chore: rework build, move to ghcr.io, build for arm64/amd64
</p>
</details>

### Changes from talos-systems/cluster-api-control-plane-provider-talos
<details><summary>3 commits</summary>
<p>

* [`579303c`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/579303cd7efaa76117f274f47a61df58a50b6f0e) release(v0.1.0-alpha.12): prepare release
* [`e0c38b3`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/e0c38b3544e923485612829061b5f476869b33f9) fix: update resources for deployment
* [`fe29dfd`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/fe29dfdf505004205843f35822d8ce4275915c44) fix: use Talos API client correctly (wrapped version)
</p>
</details>

### Changes from talos-systems/go-blockdevice
<details><summary>9 commits</summary>
<p>

* [`1d830a2`](https://github.com/talos-systems/go-blockdevice/commit/1d830a25f64f6fb96a1bedd800c0b40b107dc833) fix: revert mark the EFI partition in PMBR as bootable
* [`bec914f`](https://github.com/talos-systems/go-blockdevice/commit/bec914ffdda42abcfe642bc2cdfc9fcda56a74ee) fix: mark the EFI partition in PMBR as bootable
* [`776b37d`](https://github.com/talos-systems/go-blockdevice/commit/776b37d31de0781f098f5d9d1894fbea3f2dfa1d) feat: add options to probe disk by various sysblock parameters
* [`bb3ad73`](https://github.com/talos-systems/go-blockdevice/commit/bb3ad73f69836acc2785ec659435e24a531359e7) fix: align partition start to physical sector size
* [`8f976c2`](https://github.com/talos-systems/go-blockdevice/commit/8f976c2031108651738ebd4db69fb09758754a28) feat: replace exec.Command with go-cmd module
* [`1cf7f25`](https://github.com/talos-systems/go-blockdevice/commit/1cf7f252c38cf11ef07723de2debc27d1da6b520) fix: properly handle no child processes error from cmd.Wait
* [`04a9851`](https://github.com/talos-systems/go-blockdevice/commit/04a98510c07fe8477f598befbfe6eaec4f4b73a2) feat: implement luks encryption provider
* [`b0375e4`](https://github.com/talos-systems/go-blockdevice/commit/b0375e4267fdc6108bd9ff7a5dc97b80cd924b1d) feat: add an option to open block device with exclusive flock
* [`5a1c7f7`](https://github.com/talos-systems/go-blockdevice/commit/5a1c7f768e016c93f6c0be130ffeaf34109b5b4d) refactor: add devname into gpt.Partition, refactor probe package
</p>
</details>

### Changes from talos-systems/go-debug
<details><summary>4 commits</summary>
<p>

* [`5b292e5`](https://github.com/talos-systems/go-debug/commit/5b292e50198b8ed91c434f00e2772db394dbf0b9) feat: disable memory profiling by default
* [`c6d0ae2`](https://github.com/talos-systems/go-debug/commit/c6d0ae2c0ee099fa0940405401e6a02716a15bd8) fix: linters and CI
* [`d969f95`](https://github.com/talos-systems/go-debug/commit/d969f952af9e02feea59963671298fc236ca4399) feat: initial implementation
* [`b2044b7`](https://github.com/talos-systems/go-debug/commit/b2044b70379c84f9706de74044bd2fd6a8e891cf) Initial commit
</p>
</details>

### Changes from talos-systems/go-kmsg
<details><summary>2 commits</summary>
<p>

* [`2edcd3a`](https://github.com/talos-systems/go-kmsg/commit/2edcd3a913508e2d922776f729bfc4bcab031a8b) feat: add initial version
* [`53cdd8d`](https://github.com/talos-systems/go-kmsg/commit/53cdd8d67b9dbab692471a2d5161e7e0b3d04cca) chore: initial commit
</p>
</details>

### Changes from talos-systems/go-procfs
<details><summary>2 commits</summary>
<p>

* [`8cbc42d`](https://github.com/talos-systems/go-procfs/commit/8cbc42d3dc246a693d9b307c5358f6f7f3cb60bc) feat: provide an option to overwrite some args in AppendAll
* [`24d06a9`](https://github.com/talos-systems/go-procfs/commit/24d06a955782ed7d468f5117e986ec632f316310) refactor: remove talos kernel default args
</p>
</details>

### Changes from talos-systems/go-retry
<details><summary>3 commits</summary>
<p>

* [`7885e16`](https://github.com/talos-systems/go-retry/commit/7885e16b2cb0267bcc8b07cdd0eced14e8005864) feat: add ExpectedErrorf
* [`3d83f61`](https://github.com/talos-systems/go-retry/commit/3d83f6126c1a3a238d1d1d59bfb6273e4087bdac) feat: deprecate UnexpectedError
* [`b9dc1a9`](https://github.com/talos-systems/go-retry/commit/b9dc1a990133dd3399549b4ea199759bdfe58bb8) feat: add support for `context.Context` in Retry
</p>
</details>

### Changes from talos-systems/go-smbios
<details><summary>3 commits</summary>
<p>

* [`d3a32be`](https://github.com/talos-systems/go-smbios/commit/d3a32bea731a0c2a60ce7f5eae60253300ef27e1) fix: return UUID in middle endian only on SMBIOS >= 2.6
* [`fb425d4`](https://github.com/talos-systems/go-smbios/commit/fb425d4727e620b6a2b6ba49e405a2c6f0e46304) feat: add memory device
* [`0bb4f96`](https://github.com/talos-systems/go-smbios/commit/0bb4f96a6679e8fc958903c4f451ca068f8e3c41) feat: add physical memory array
</p>
</details>

### Dependency Changes

* **github.com/go-logr/logr**                                            ee2de8da5be6 -> v0.4.0
* **github.com/hashicorp/go-multierror**                                 v1.1.0 -> v1.1.1
* **github.com/onsi/ginkgo**                                             v1.15.0 -> v1.16.2
* **github.com/onsi/gomega**                                             v1.10.1 -> v1.12.0
* **github.com/stretchr/testify**                                        v1.7.0 **_new_**
* **github.com/talos-systems/cluster-api-bootstrap-provider-talos**      v0.2.0-alpha.11 -> v0.2.0-alpha.12
* **github.com/talos-systems/cluster-api-control-plane-provider-talos**  v0.1.0-alpha.11 -> v0.1.0-alpha.12
* **github.com/talos-systems/go-blockdevice**                            f2728a581972 -> 1d830a25f64f
* **github.com/talos-systems/go-debug**                                  v0.2.0 **_new_**
* **github.com/talos-systems/go-kmsg**                                   v0.1.0 **_new_**
* **github.com/talos-systems/go-procfs**                                 a82654edcec1 -> 8cbc42d3dc24
* **github.com/talos-systems/go-retry**                                  v0.2.0 -> v0.3.0
* **github.com/talos-systems/go-smbios**                                 80196199691e -> d3a32bea731a
* **github.com/talos-systems/talos/pkg/machinery**                       1d8e9674a91b -> 8d73bc5999b4
* **go.uber.org/zap**                                                    v1.14.1 -> v1.16.0
* **golang.org/x/mod**                                                   v0.4.0 **_new_**
* **golang.org/x/net**                                                   0714010a04ed **_new_**
* **golang.org/x/sync**                                                  67f06af15bc9 -> 036812b2e83c
* **golang.org/x/sys**                                                   489259a85091 -> 0981d6026fa6
* **golang.org/x/tools**                                                 2dba1e4ea05c **_new_**
* **google.golang.org/grpc**                                             v1.36.0 -> v1.37.1
* **google.golang.org/protobuf**                                         v1.26.0 **_new_**

Previous release can be found at [v0.2.0](https://github.com/talos-systems/sidero/releases/tag/v0.2.0)

<a name="v0.3.0-alpha.0"></a>
## [v0.3.0-alpha.0](https://github.com/talos-systems/sidero/compare/v0.2.0...v0.3.0-alpha.0) (2021-05-13)

### Chore

* parse "boolean" variables
* bump dependencies
* allow building with debug handlers
* add golangci-lint configuration
* fix `make help` command
* document Makefile target
* disable dependabot rebases
* update dependecies
* add dependabot config
* improve build system

### Docs

* fix install disk name in the examples
* fix typo
* add guide for upgrades
* fix the docs around CONTROL_PLANE_ENDPOINT
* create v0.2 docs and add note about specifying TALOS_VERSION

### Feat

* ship ServerClass "any"
* enable auto-setup of BMC
* inject iPXE script into the iPXE binaries
* pull the Sidero configuration as `clusterctl` variables
* add an option to reference IPMI creds via the secret refs
* pull in new version of go-smbios for UUID endiannes fix
* build Sidero for amd64 and arm64, support arm64 in the agent

### Fix

* back down resource requests
* remove erroneous wg.Add in environment controller

### Test

* fix the Environment args to support UEFI boot
* add missing empty tests
* port improvements from Talos
* run unit tests on CI, skip broken


<a name="v0.2.0"></a>
## [v0.2.0](https://github.com/talos-systems/talos/compare/v0.2.0-beta.0...v0.2.0) (2021-04-09)

### Release

* **v0.2.0:** prepare release


<a name="v0.2.0-beta.0"></a>
## [v0.2.0-beta.0](https://github.com/talos-systems/talos/compare/v0.2.0-alpha.1...v0.2.0-beta.0) (2021-04-08)

### Chore

* use common 'setup-ci' function

### Feat

* add support for DNS resolution in the agent

### Fix

* break the potential endless reconcile loop of ServerClasses
* don't marshal the machine configuration via machinery package

### Release

* **v0.2.0-beta.0:** prepare release

### Test

* drop installer image from the server patch


<a name="v0.2.0-alpha.1"></a>
## [v0.2.0-alpha.1](https://github.com/talos-systems/talos/compare/v0.2.0-alpha.0...v0.2.0-alpha.1) (2021-04-06)


<a name="v0.2.0-alpha.0"></a>
## [v0.2.0-alpha.0](https://github.com/talos-systems/talos/compare/v0.1.0...v0.2.0-alpha.0) (2021-04-06)

### Docs

* update bootstrapping guide
* setup google analytics
* add IPMI info, fix links, update environment
* add note on installation disk
* fix typo on Server Classes page
* describe the command to install talosctl in the bootstrap guide

### Feat

* update Talos and machinery to 0.9.1
* update sidero to use newer talos
* serve assets from TFTP folder in IPXE HTTP server
* bump talos to 0.8.0 release

### Fix

* log retry errors
* add kubernetes version to cluster template
* prevent server orphaning

### Release

* **v0.2.0-alpha.0:** prepare release

### Test

* enable control plane scale down test


<a name="v0.1.0"></a>
## [v0.1.0](https://github.com/talos-systems/talos/compare/v0.1.0-beta.0...v0.1.0) (2021-01-15)

### Fix

* prevent server orphaning

### Release

* **v0.1.0:** prepare release

### Test

* fix test flakiness on workers scale down test


<a name="v0.1.0-beta.0"></a>
## [v0.1.0-beta.0](https://github.com/talos-systems/talos/compare/v0.1.0-alpha.12...v0.1.0-beta.0) (2020-12-23)

### Chore

* bump versions of packages used in the build
* fix analytics
* bump google analytics plugin

### Docs

* update local Sfyra docs
* improve docs on metadata and environments
* add flow charts for PXE and installation
* add CRD documentation
* clarify cluster template environment variables
* clarify process of fetching talosconfig

### Feat

* bump Talos to 0.8.0-beta.0
* send heartbeat from agent while wipe is in progress
* align default subnets with Talos

### Fix

* overwrite kernel args from --extra-agent-kernel-args

### Release

* **v0.1.0-beta.0:** prepare release

### Test

* add simulated power management failures in testing mode
* add a test to deploy and destroy workload cluster


<a name="v0.1.0-alpha.12"></a>
## [v0.1.0-alpha.12](https://github.com/talos-systems/talos/compare/v0.1.0-alpha.11...v0.1.0-alpha.12) (2020-12-02)

### Feat

* rework binding between Servers and MetalMachines

### Release

* **v0.1.0-alpha.12:** prepare release


<a name="v0.1.0-alpha.11"></a>
## [v0.1.0-alpha.11](https://github.com/talos-systems/talos/compare/v0.1.0-alpha.10...v0.1.0-alpha.11) (2020-11-30)

### Docs

* fix typo in asset URLs
* fix patching example

### Fix

* bump resource limits for the sidero pods
* add timeout to retry wipe IPMI commands (PXE + powercycle)

### Refactor

* unify power management under ServerController
* remove dependency on talos Go module

### Release

* **v0.1.0-alpha.11:** prepare release


<a name="v0.1.0-alpha.10"></a>
## [v0.1.0-alpha.10](https://github.com/talos-systems/talos/compare/v0.1.0-alpha.9...v0.1.0-alpha.10) (2020-11-17)

### Chore

* bump Talos package dependency to 0.7-beta.1

### Feat

* boot servers via PXE only once by default
* allow for extra kernel args in agent environment
* add power status to servers
* don't power off the server in discovery agent
* make "insecure-wipe" default, use new fast wipe method

### Fix

* add info log when no servers in serverclass

### Refactor

* use GetDisks from blockdevice library
* add ErrNoServersInServerClass

### Release

* **v0.1.0-alpha.10:** prepare release


<a name="v0.1.0-alpha.9"></a>
## [v0.1.0-alpha.9](https://github.com/talos-systems/talos/compare/v0.1.0-alpha.8...v0.1.0-alpha.9) (2020-11-10)

### Fix

* update pkgs
* ignore missing partition error

### Release

* **v0.1.0-alpha.9:** prepare release


<a name="v0.1.0-alpha.8"></a>
## [v0.1.0-alpha.8](https://github.com/talos-systems/talos/compare/v0.1.0-alpha.7...v0.1.0-alpha.8) (2020-11-07)

### Chore

* fix docker login
* move to ghcr.io

### Docs

* fix typos in bootstrap guides
* update bootstrap docs links
* add note around server acceptance

### Feat

* add option for insecure wipe

### Release

* **v0.1.0-alpha.8:** prepare release


<a name="v0.1.0-alpha.7"></a>
## [v0.1.0-alpha.7](https://github.com/talos-systems/talos/compare/v0.1.0-alpha.6...v0.1.0-alpha.7) (2020-10-31)

### Docs

* update site URL and add robots.txt
* add a metadata section
* expand server acceptance doc
* add non-UEFI clause to PXE example
* add links between concepts and configuration sections
* add chicken and egg note to overview
* add support for versioned docs

### Feat

* wipe disks concurrently in the agent

### Fix

* log error with error log
* wipe disk fully in the agent

### Release

* **v0.1.0-alpha.7:** prepare release


<a name="v0.1.0-alpha.6"></a>
## [v0.1.0-alpha.6](https://github.com/talos-systems/talos/compare/v0.1.0-alpha.5...v0.1.0-alpha.6) (2020-10-21)

### Feat

* add support for control plane port

### Release

* **v0.1.0-alpha.6:** prepare release


<a name="v0.1.0-alpha.5"></a>
## [v0.1.0-alpha.5](https://github.com/talos-systems/talos/compare/v0.1.0-alpha.4...v0.1.0-alpha.5) (2020-10-21)

### Chore

* update talos version

### Feat

* update CAPI providers

### Fix

* don't reset read only disks

### Release

* **v0.1.0-alpha.5:** prepare release


<a name="v0.1.0-alpha.4"></a>
## [v0.1.0-alpha.4](https://github.com/talos-systems/talos/compare/v0.1.0-alpha.3...v0.1.0-alpha.4) (2020-10-19)

### Chore

* fix formatting

### Docs

* add server acceptance note

### Feat

* support config patches at the serverclass level
* discover server IPs on registration, emit server events
* add hostname to the server spec

### Fix

* use APIReader in server controller to avoid caching issues
* proper event patching, powercycle vs. poweron
* handle agent errors gracefully
* use efiboot option

### Release

* **v0.1.0-alpha.4:** prepare release

### Test

* pull in the new version of go-loadbalancer
* check if servers are reset on acceptance
* add flags for modifying workload cluster installer image
* bump Talos to 0.7.0-alpha.6 for Sfyra
* enable verbose logs for CLI version of loadbalancer


<a name="v0.1.0-alpha.3"></a>
## [v0.1.0-alpha.3](https://github.com/talos-systems/talos/compare/v0.1.0-alpha.2...v0.1.0-alpha.3) (2020-10-09)

### Chore

* refactor `reconcile` method in `environment_controller.go`
* introduce RequeueAfter in metalcluster controller
* remove talos dependency from metadata server
* update Sfyra for the new Talos release

### Feat

* support server acceptance
* power down servers !InUse && IsClean
* ensure servers are clean
* set servers as dirty by default

### Fix

* add ipmitool to metal-controller-manager
* exit when port conflicts happen
* delete machine instead of metalmachine in reset test
* handle conflicts on server update

### Refactor

* break apart metadata server code
* split 'sfyra' CLI into multiple subcommands

### Release

* **v0.1.0-alpha.3:** prepare release

### Test

* add scale test
* bump talos for halt & healh check fixes
* apply talos.shutdown=halt kernel argument
* break serverclass creation into function and allow dummy servers
* fix patching for sidero-controller-manager
* provide list of expected control plane/worker nodes
* extract `capi.Cluster` implementation
* use dedicated config for the sfyra tests
* fix nil pointer panic


<a name="v0.1.0-alpha.2"></a>
## [v0.1.0-alpha.2](https://github.com/talos-systems/talos/compare/v0.1.0-alpha.1...v0.1.0-alpha.2) (2020-09-28)

### Chore

* update reference to internal registry
* merge Sfyra into Sidero repository
* add additional logging for metadata server
* update namespace to caps-system

### Docs

* update Sfyra docs
* add a post-pivot/create first workload cluster guide
* add note about templated cluster manifest
* add full bootstrapping guide

### Feat

* add reset functionality
* allow qualifers to be partially equal
* support environment ref at server and server class level
* add serverclass as an owner to servers if needed

### Fix

* close file in TFTP handler
* use proper context in metadata server
* ensure proper checking of "in-use" status when fetching metadata
* address panic in PartialEqual
* revert "update namespace to caps-system"
* update labels and labelselectors for each app
* move to update instead of patching server inuse

### Refactor

* extract common ManagmentClient interface, add HTTP API

### Release

* **v0.1.0-alpha.2:** prepare release

### Test

* add reset test


<a name="v0.1.0-alpha.1"></a>
## [v0.1.0-alpha.1](https://github.com/talos-systems/talos/compare/v0.1.0-alpha.0...v0.1.0-alpha.1) (2020-08-25)

### Chore

* expire server discovery requests

### Docs

* iterate on docs

### Feat

* fetch hostname if available

### Fix

* ensure servers always get patched

### Refactor

* move `internal/app` to `app` so that we can expose API

### Release

* **v0.1.0-alpha.1:** prepare release


<a name="v0.1.0-alpha.0"></a>
## v0.1.0-alpha.0 (2020-08-17)

### Chore

* update drone file for tagged releases
* rename cluster api provider to match capi norms
* add move label to metal controller manager CRDs
* update drone pipeline type
* address confrom errors
* enable drone builds
* format generated files
* address linter errors
* use sidero-system namespace
* use 'sidero' instead of 'arges' and 'metal'
* migrate projects

### Docs

* add website
* address markdown lint errors
* start sidero docs

### Feat

* add cluster-template for sidero

### Fix

* ensure servers only get picked once
* ensure we don't clobber extraargs
* update manifests
* update kustomize configs
* set asset download context timeout to 5 minutes
* patch the release images
* refactor kustomization to fix releases
* remove stray reference to capi namespace.yaml

### Release

* **v0.1.0-alpha.0:** prepare release
