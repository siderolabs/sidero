## [Sidero 0.4.1-beta.0](https://github.com/talos-systems/sidero/releases/tag/v0.4.1-beta.0) (2021-12-09)

Welcome to the v0.4.1-beta.0 release of Sidero!



Please try out the release binaries and report any issues at
https://github.com/talos-systems/sidero/issues.

### IPMI PXE Method

IPMI PXE method (UEFI, BIOS) can now be configured with `SIDERO_CONTROLLER_MANAGER_IPMI_PXE_METHOD` while installing Sidero.


### Contributors

* Andrey Smirnov
* Artem Chernyshev

### Changes
<details><summary>9 commits</summary>
<p>

* [`9a79c02`](https://github.com/talos-systems/sidero/commit/9a79c02a926e67fcbfa2e39bc4048c6640f9391c) chore: disable check for one commit
* [`b6f860f`](https://github.com/talos-systems/sidero/commit/b6f860f2c4c0c77d75e37ce3ee66c8a59db91ad2) feat: extend information printed in the iPXE script, add retries
* [`fec4d4b`](https://github.com/talos-systems/sidero/commit/fec4d4ba5729b0002dd8ab9803be596e818215e8) feat: provide a way to configure IPMI PXE method
* [`3e2ae6d`](https://github.com/talos-systems/sidero/commit/3e2ae6d3fc900753d81523d31a5712e2550ad0d1) fix: check for server power state when in use
* [`b2a693e`](https://github.com/talos-systems/sidero/commit/b2a693e0092107b74bc8553f0e51a882c8a1695b) fix: update CAPI resources versions to v1alpha4 in the cluster template
* [`4fdcbb3`](https://github.com/talos-systems/sidero/commit/4fdcbb3283f2cf6d3ffdf7f3c5ae8abdf4394669) feat: automatically append talos.config to the Environment
* [`b8553d4`](https://github.com/talos-systems/sidero/commit/b8553d4fee44608f07901a06e2f6ccd4d38f3511) fix: hide password from ipmitool args
* [`198f102`](https://github.com/talos-systems/sidero/commit/198f10233098cd1cf4683db0f13485f8d99933df) fix: drop into the agent for non-allocated servers
* [`ce626cf`](https://github.com/talos-systems/sidero/commit/ce626cf049d7aa5dff08a7092a368f7062c21430) feat: allow defining http server and api server ports separately
</p>
</details>

### Dependency Changes

This release has no dependency changes

Previous release can be found at [v0.4.0](https://github.com/talos-systems/sidero/releases/tag/v0.4.0)

## [Sidero 0.4.0](https://github.com/talos-systems/sidero/releases/tag/v0.4.0) (2021-10-18)

Welcome to the v0.4.0 release of Sidero!



Please try out the release binaries and report any issues at
https://github.com/talos-systems/sidero/issues.

### BMC Port

Sidero now supports the ability to specify the port in a server's BMC info. By default, this value will be determined by talking directly to the BMC if possible, with a fallback to port 623. The value can also simply be specied as part of editing the Server resource directly.


### CAPI v1alpha4

This release of Sidero brings compatibility with CAPI v1alpha4.


### Contributors

* Artem Chernyshev
* Andrey Smirnov
* Alexey Palazhchenko
* Andrey Smirnov
* Spencer Smith
* Artem Chernyshev
* Spencer Smith
* Gerard de Leeuw
* Noel Georgi
* Serge Logvinov
* Gerard de Leeuw
* Khue Doan
* Olli Janatuinen
* Seán C McCord

### Changes
<details><summary>34 commits</summary>
<p>

* [`74e7e10`](https://github.com/talos-systems/sidero/commit/74e7e10494821d9f83a4ec0ba7e5eb5830a14b07) fix: api-port parameter not being used
* [`6e073e7`](https://github.com/talos-systems/sidero/commit/6e073e7ab1c4f8fbb84e4d61b99d943ae19d3393) fix: commit updated release.toml to master
* [`80345c8`](https://github.com/talos-systems/sidero/commit/80345c8b65874f3c50e9cf57492459ea6c1a5e56) release(v0.4.0-alpha.2): prepare release
* [`af646a7`](https://github.com/talos-systems/sidero/commit/af646a7170f5a5a87d6e604f59d5b66fc3d35fd1) chore: bump Go deps, tools, pkgs, alpine versions
* [`bb52e71`](https://github.com/talos-systems/sidero/commit/bb52e71861072104a5f1a8bc460b01fcf8c84932) feat: support CAPI v1alpha4
* [`d923bb0`](https://github.com/talos-systems/sidero/commit/d923bb0f7f5681ada12159b0ffc518ca40c6056d) docs: add notes on version compatibility and improve the bootstrapping guide
* [`a2bb2d0`](https://github.com/talos-systems/sidero/commit/a2bb2d0919c76b7e051c50f600bab05a9c9d0366) chore: fix docker version in the drone pipeline
* [`e52071d`](https://github.com/talos-systems/sidero/commit/e52071d7f9e5a05dc8b3f9d2ee48265a55c6fc14) fix: shutdown sidero-controller-manager when any component fails
* [`afa114b`](https://github.com/talos-systems/sidero/commit/afa114b18f15798a5c9bbcc7d6c1baa9952ee0af) fix: broken url, ServerClass CR spec
* [`b37f43f`](https://github.com/talos-systems/sidero/commit/b37f43fe6274bcd30ffb9affffefb11335b9d3a2) fix: resource config links
* [`aa2b3f0`](https://github.com/talos-systems/sidero/commit/aa2b3f0e3512ea8049b7fbdd2a423b209126365e) chore: bump dependencies via dependabot
* [`f1c1608`](https://github.com/talos-systems/sidero/commit/f1c16082071cec651a7de562ac101a11a1689db7) chore: bump Talos to 0.11.5
* [`8695371`](https://github.com/talos-systems/sidero/commit/8695371253b9fc8732e9a3da493c0ede2c33c50b) release(v0.4.0-alpha.1): prepare release
* [`7bdee0f`](https://github.com/talos-systems/sidero/commit/7bdee0f63fef443f4d51329f40c2db8429030008) fix: update sidero IPMI user to work properly on idrac
* [`acd82e0`](https://github.com/talos-systems/sidero/commit/acd82e0611843d3c273824d7daad53b3abed3352) docs: bump clusterctl from v0.3.14 to v0.3.21 in /docs/website
* [`fb0da3c`](https://github.com/talos-systems/sidero/commit/fb0da3cf5087f67cd6dc40381f228c756abedd4b) release(v0.4.0-alpha.0): prepare release
* [`ee36c74`](https://github.com/talos-systems/sidero/commit/ee36c745016a324220b5d17b7ef5fff9d2ce85a8) docs: redirect latest to v0.3
* [`7cdae7b`](https://github.com/talos-systems/sidero/commit/7cdae7b616db7720627ab2e5b409e2f22b0085e0) feat: add ability to specify ports for BMC
* [`c14e055`](https://github.com/talos-systems/sidero/commit/c14e055f4c26ee7df38742bb80d385b91e6191b7) chore: bump Talos to 0.11.0-beta.3
* [`7170777`](https://github.com/talos-systems/sidero/commit/7170777d603bc52ef6ba3a02669266d421acc65f) fix: make sure powercycle condition gets properly update
* [`90e7804`](https://github.com/talos-systems/sidero/commit/90e78046409fdf5dd4f16b1e8ae90376d6c813ac) chore: bump dependencies in go.mod
* [`dc70167`](https://github.com/talos-systems/sidero/commit/dc70167c3d11613ecd43291c3930519eb356124c) docs: fix duplicate titles in the guides
* [`fd1fae7`](https://github.com/talos-systems/sidero/commit/fd1fae793148fdc7c0eca22a9e0a82618409644f) chore: update Talos to 0.11.0-beta.2
* [`1f8f141`](https://github.com/talos-systems/sidero/commit/1f8f14178c06d96ed47ca1399340f3ca589c140d) docs: promote 0.3 docs to be the latest
* [`ea3016f`](https://github.com/talos-systems/sidero/commit/ea3016fc19bd20436e70f45206eed08e7c974524) fix: update Sfyra to install CAPI v0.3
* [`8e49ddf`](https://github.com/talos-systems/sidero/commit/8e49ddf49b7d02ed08130037fd4ded7982326ca7) chore: update to latest stable talos providers
* [`1155004`](https://github.com/talos-systems/sidero/commit/115500459856f0181d7cb1936cb9629b97a01d59) docs: document using ISO for mgmt plane
* [`a5b3e67`](https://github.com/talos-systems/sidero/commit/a5b3e67311fca76e1f50139f09fc17393035de18) docs: add docs for server removal/decommissioning
* [`c7ae88a`](https://github.com/talos-systems/sidero/commit/c7ae88a49d73fb247be60f11e37566c2056f96b3) docs: ensure tutorial docs are present in sidebar
* [`83b0875`](https://github.com/talos-systems/sidero/commit/83b08757190081d0ab0205ffea741703328c6f1e) docs: rework guides into Tutorial
* [`ee31160`](https://github.com/talos-systems/sidero/commit/ee31160b5a8b3b5844585bdc9567374345caf14a) fix: make sure components of agent environment are of proper arch
* [`dfe2c85`](https://github.com/talos-systems/sidero/commit/dfe2c85438c064d85d85706c313efc55098b8fe5) chore: bump dependencies via dependabot
* [`bf2be1b`](https://github.com/talos-systems/sidero/commit/bf2be1b1569333034cadd2dacc023538221b2930) feat: update default Talos version to v0.10.3
* [`dfeaeec`](https://github.com/talos-systems/sidero/commit/dfeaeec308e2f45ae815943b343266e1c341a262) docs: add notes on running Sidero on RPi4
</p>
</details>

### Changes since v0.4.0-alpha.2
<details><summary>2 commits</summary>
<p>

* [`74e7e10`](https://github.com/talos-systems/sidero/commit/74e7e10494821d9f83a4ec0ba7e5eb5830a14b07) fix: api-port parameter not being used
* [`6e073e7`](https://github.com/talos-systems/sidero/commit/6e073e7ab1c4f8fbb84e4d61b99d943ae19d3393) fix: commit updated release.toml to master
</p>
</details>

### Changes from talos-systems/cluster-api-bootstrap-provider-talos
<details><summary>26 commits</summary>
<p>

* [`2f1364c`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/2f1364c966515d241d21723c44ff3aaab543edfa) release(v0.4.0): prepare release
* [`04742b9`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/04742b96bf757413c88d0f15bee91679644f0337) feat: import fixes/updates from kubeadm bootstrap provider
* [`548b7fb`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/548b7fbd38b89b9790a0daa2380fddb34157cdd5) release(v0.4.0-alpha.0): prepare release
* [`442ee41`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/442ee41bafb2a912e49928c5d61f52c4c61a2593) test: don't set the talosconfig owner ref to the machine
* [`8c7fec8`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/8c7fec8e373bd12609f6274d79ca07d187212d91) fix: don't write incomplete `<cluster>-ca` secret for configtype none
* [`f46c83d`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/f46c83d328ee44db2ccb5eef67b366cc73c13319) feat: bump Talos machinery to 0.12.3
* [`7b760cf`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/7b760cf69ecab93200821dded931171657a5dedc) feat: support CAPI v1alpha4
* [`3147ba4`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/3147ba4fe57b88975133c598c226ff4e397efb44) release(v0.3.0-alpha.1): prepare release
* [`977121a`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/977121ad14dc0637f7c4282e69a4ee26e28372d4) fix: construct properly data secret name
* [`f8c75c8`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/f8c75c89c4653de30165fb1538e906256a4eec66) fix: update metadata.yaml for v0.3 of CABPT
* [`db60f9e`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/db60f9eb0697c4949be9c00cf8dc7787d383bad2) release(v0.3.0-alpha.0): prepare release
* [`755a2dd`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/755a2dd90c3668db89f8eae14f60db4564764475) fix: update Talos machinery to 0.12, fix secrets persistence
* [`f91b032`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/f91b032935776c1224f824cc860bfa4df5e220b1) fix: use bootstrap data secret names
* [`6bff239`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/6bff2393840655c2361def455b601511b86ba71f) chore: use Go 1.17
* [`56fb73b`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/56fb73b53f41b91b12ba2b3c331d7a04b7263a17) test: add test for the second machine
* [`e5b7738`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/e5b773833120fdd7ca4d57e0a0a4fe781495bf7e) test: add more tests
* [`bc4105d`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/bc4105d9e8366d4e840705a6cecfbc81bdcca00a) test: wait for CAPI availability
* [`c82b8ab`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/c82b8ab47bca5313cb96df1b70de0914da285331) chore: make versions configurable
* [`5594c96`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/5594c96daa55fb9fc9af585e8f2fc26551ce9bb5) chore: use codecov uploader from build-container
* [`cced038`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/cced038257d3eec5b7c48bc524de5165b5734496) chore: fix license headers
* [`7b5dc51`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/7b5dc51e83a54a1f5fa707c66a296ca9514c8722) chore: do not run tests on ARM
* [`d6258cf`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/d6258cf21778149a254d9669b03ac10bae9e0955) chore: improve tests runner
* [`c6ce363`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/c6ce36375ef145760647c632d64a9a3c93574e4b) chore: sign Drone CI configuration
* [`ad592d1`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/ad592d16fa8397f88a28e6a4151bc64b0a1c097d) chore: add basic integration test
* [`9fb0d07`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/9fb0d07ca4d2e8333b0b61ee0fe0ba3e6660489f) chore: add missing LICENSE file
* [`acf18d2`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/acf18d2bb09aab64687c1fccf1e628ef76e9cff8) chore: update machinery to v0.11.3
</p>
</details>

### Changes from talos-systems/cluster-api-control-plane-provider-talos
<details><summary>14 commits</summary>
<p>

* [`1bec112`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/1bec1128f39c0abc243a99a92aaf5bf9917439b2) release(v0.3.0): prepare release
* [`1662815`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/166281550865c66ed4f6a015c76c94443b43f0fe) feat: catch up with Kubeadm in terms of conditions
* [`43eb75b`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/43eb75b439d43e87e970df69f49f0abbde047e51) release(v0.3.0-alpha.0): prepare release
* [`48d834b`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/48d834b5dfb364b8e9ae2269771e41a2dc646692) feat: support CAPI v1alpha4
* [`14c6e72`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/14c6e7224703b8838c1790c2014847f991367ff2) release(v0.2.0-alpha.0): prepare release
* [`cd6417d`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/cd6417dd365aa89984703016b61c57e6b33b3b68) fix: update metadata.yaml for v0.2 of CACPPT
* [`8b52b8a`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/8b52b8addd9fa4235c542b0b8554a76f5c76a643) chore: update go to 1.17
* [`86d679a`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/86d679a44e543789474c0b8edaf435a764f7dd2e) chore: update cabpt to v0.3.0
* [`a616f4b`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/a616f4b4bd3b208595cd102eb9e32c8a31b95e18) test: add machine removal test
* [`6ad6aac`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/6ad6aac1315ad5bc8e1264af6162863418cdb280) test: implement scale up and down tests and fix found issues
* [`9435b12`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/9435b1247f010bee00b4a8e4dc592121a0eb2449) chore: add e2e test running on AWS infra
* [`4c7d42c`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/4c7d42caf79ca209f5cda84db2eb712433d3c68b) chore: update bootstrap provider
* [`119b969`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/119b969be2fe152a0e8a63d189563deed55110b4) fix: clean up couple small issues in the etcd member audit code
* [`9be7b88`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/9be7b88bf4a14aec584fe68561c3fda3fbeaf990) chore: update bootstrap provider to stable release
</p>
</details>

### Changes from talos-systems/go-blockdevice
<details><summary>13 commits</summary>
<p>

* [`70d2865`](https://github.com/talos-systems/go-blockdevice/commit/70d28650b398a14469cbb5356417355b0ba62956) fix: try to find cdrom disks
* [`667bf53`](https://github.com/talos-systems/go-blockdevice/commit/667bf539b99ac34b629a0103ef7a7278a5a5f35d) fix: revert gpt partition not found
* [`d7d4cdd`](https://github.com/talos-systems/go-blockdevice/commit/d7d4cdd7ac56c82caab19246b5decd59f12195eb) fix: gpt partition not found
* [`33afba3`](https://github.com/talos-systems/go-blockdevice/commit/33afba347c0dce38a436c46a0aac26d2f99427c1) fix: also open in readonly mode when running `All` lookup method
* [`e367f9d`](https://github.com/talos-systems/go-blockdevice/commit/e367f9dc7fa935f11672de0fdc8a89429285a07a) feat: make probe always open blockdevices in readonly mode
* [`d981156`](https://github.com/talos-systems/go-blockdevice/commit/d9811569588ba44be878a00ce316f59a37abed8b) fix: allow Build for Windows
* [`fe24303`](https://github.com/talos-systems/go-blockdevice/commit/fe2430349e9d734ce6dbf4e7b2e0f8a37bb22679) fix: perform correct PMBR partition calculations
* [`2ec0c3c`](https://github.com/talos-systems/go-blockdevice/commit/2ec0c3cc0ff5ff705ed5c910ca1bcd5d93c7b102) fix: preserve the PMBR bootable flag when opening GPT partition
* [`87816a8`](https://github.com/talos-systems/go-blockdevice/commit/87816a81cefc728cfe3cb221b476d8ed4b609fd8) feat: align partition to minimum I/O size
* [`c34b59f`](https://github.com/talos-systems/go-blockdevice/commit/c34b59fb33a7ad8be18bb19bc8c8d8294b4b3a78) feat: expose more encryption options in the LUKS module
* [`30c2bc3`](https://github.com/talos-systems/go-blockdevice/commit/30c2bc3cb62af52f0aea9ce347923b0649fb7928) feat: mark MBR bootable
* [`1292574`](https://github.com/talos-systems/go-blockdevice/commit/1292574643e06512255fb0f45107e0c296eb5a3b) fix: make disk type matcher parser case insensitive
* [`b77400e`](https://github.com/talos-systems/go-blockdevice/commit/b77400e0a7261bf25da77c1f28c2f393f367bfa9) fix: properly detect nvme and sd card disk types
</p>
</details>

### Changes from talos-systems/go-debug
<details><summary>1 commit</summary>
<p>

* [`3d0a6e1`](https://github.com/talos-systems/go-debug/commit/3d0a6e1bf5e3c521e83ead2c8b7faad3638b8c5d) feat: race build tag flag detector
</p>
</details>

### Changes from talos-systems/go-kmsg
<details><summary>1 commit</summary>
<p>

* [`b08e4d3`](https://github.com/talos-systems/go-kmsg/commit/b08e4d36a2f3df0a3d031b1a3028e2d6e4c26710) feat: replace tab character with space in console output
</p>
</details>

### Changes from talos-systems/go-retry
<details><summary>1 commit</summary>
<p>

* [`c78cc95`](https://github.com/talos-systems/go-retry/commit/c78cc953d9e95992575305b4e8648392c6c9b9e6) fix: implement `errors.Is` for all errors in the set
</p>
</details>

### Dependency Changes

* **github.com/evanphx/json-patch**                                      v4.9.0 -> v4.11.0
* **github.com/onsi/ginkgo**                                             v1.16.3 -> v1.16.4
* **github.com/onsi/gomega**                                             v1.13.0 -> v1.16.0
* **github.com/talos-systems/cluster-api-bootstrap-provider-talos**      v0.2.0 -> v0.4.0
* **github.com/talos-systems/cluster-api-control-plane-provider-talos**  v0.1.0 -> v0.3.0
* **github.com/talos-systems/go-blockdevice**                            1d830a25f64f -> v0.2.4
* **github.com/talos-systems/go-debug**                                  v0.2.0 -> v0.2.1
* **github.com/talos-systems/go-kmsg**                                   v0.1.0 -> v0.1.1
* **github.com/talos-systems/go-retry**                                  v0.3.0 -> v0.3.1
* **github.com/talos-systems/talos/pkg/machinery**                       828772cec9a3 -> v0.13.0
* **golang.org/x/net**                                                   0714010a04ed -> 853a461950ff
* **golang.org/x/sys**                                                   0981d6026fa6 -> 39ccf1dd6fa6
* **google.golang.org/grpc**                                             v1.38.0 -> v1.41.0
* **google.golang.org/protobuf**                                         v1.26.0 -> v1.27.1
* **k8s.io/api**                                                         v0.19.3 -> v0.22.2
* **k8s.io/apiextensions-apiserver**                                     v0.19.1 -> v0.22.2
* **k8s.io/apimachinery**                                                v0.19.3 -> v0.22.2
* **k8s.io/client-go**                                                   v0.19.3 -> v0.22.2
* **k8s.io/utils**                                                       67b214c5f920 -> bdf08cb9a70a
* **sigs.k8s.io/cluster-api**                                            v0.3.12 -> v0.4.3
* **sigs.k8s.io/controller-runtime**                                     v0.6.3 -> v0.9.7

Previous release can be found at [v0.3.2](https://github.com/talos-systems/sidero/releases/tag/v0.3.2)

## [Sidero 0.4.0-alpha.2](https://github.com/talos-systems/sidero/releases/tag/v0.4.0-alpha.2) (2021-10-13)

Welcome to the v0.4.0-alpha.2 release of Sidero!  
*This is a pre-release of Sidero*



Please try out the release binaries and report any issues at
https://github.com/talos-systems/sidero/issues.

### BMC Port

Sidero now supports the ability to specify the port in a server's BMC info. By default, this value will be determined by talking directly to the BMC if possible, with a fallback to port 623. The value can also simply be specied as part of editing the Server resource directly.


### CAPI v1alpha4

This release of Sidero brings compatibility with CAPI v1alpha4.


### Contributors

* Andrey Smirnov
* Artem Chernyshev
* Alexey Palazhchenko
* Andrey Smirnov
* Spencer Smith
* Artem Chernyshev
* Spencer Smith
* Gerard de Leeuw
* Noel Georgi
* Serge Logvinov
* Gerard de Leeuw
* Khue Doan
* Olli Janatuinen
* Seán C McCord

### Changes
<details><summary>31 commits</summary>
<p>

* [`af646a7`](https://github.com/talos-systems/sidero/commit/af646a7170f5a5a87d6e604f59d5b66fc3d35fd1) chore: bump Go deps, tools, pkgs, alpine versions
* [`bb52e71`](https://github.com/talos-systems/sidero/commit/bb52e71861072104a5f1a8bc460b01fcf8c84932) feat: support CAPI v1alpha4
* [`d923bb0`](https://github.com/talos-systems/sidero/commit/d923bb0f7f5681ada12159b0ffc518ca40c6056d) docs: add notes on version compatibility and improve the bootstrapping guide
* [`a2bb2d0`](https://github.com/talos-systems/sidero/commit/a2bb2d0919c76b7e051c50f600bab05a9c9d0366) chore: fix docker version in the drone pipeline
* [`e52071d`](https://github.com/talos-systems/sidero/commit/e52071d7f9e5a05dc8b3f9d2ee48265a55c6fc14) fix: shutdown sidero-controller-manager when any component fails
* [`afa114b`](https://github.com/talos-systems/sidero/commit/afa114b18f15798a5c9bbcc7d6c1baa9952ee0af) fix: broken url, ServerClass CR spec
* [`b37f43f`](https://github.com/talos-systems/sidero/commit/b37f43fe6274bcd30ffb9affffefb11335b9d3a2) fix: resource config links
* [`aa2b3f0`](https://github.com/talos-systems/sidero/commit/aa2b3f0e3512ea8049b7fbdd2a423b209126365e) chore: bump dependencies via dependabot
* [`f1c1608`](https://github.com/talos-systems/sidero/commit/f1c16082071cec651a7de562ac101a11a1689db7) chore: bump Talos to 0.11.5
* [`8695371`](https://github.com/talos-systems/sidero/commit/8695371253b9fc8732e9a3da493c0ede2c33c50b) release(v0.4.0-alpha.1): prepare release
* [`7bdee0f`](https://github.com/talos-systems/sidero/commit/7bdee0f63fef443f4d51329f40c2db8429030008) fix: update sidero IPMI user to work properly on idrac
* [`acd82e0`](https://github.com/talos-systems/sidero/commit/acd82e0611843d3c273824d7daad53b3abed3352) docs: bump clusterctl from v0.3.14 to v0.3.21 in /docs/website
* [`fb0da3c`](https://github.com/talos-systems/sidero/commit/fb0da3cf5087f67cd6dc40381f228c756abedd4b) release(v0.4.0-alpha.0): prepare release
* [`ee36c74`](https://github.com/talos-systems/sidero/commit/ee36c745016a324220b5d17b7ef5fff9d2ce85a8) docs: redirect latest to v0.3
* [`7cdae7b`](https://github.com/talos-systems/sidero/commit/7cdae7b616db7720627ab2e5b409e2f22b0085e0) feat: add ability to specify ports for BMC
* [`c14e055`](https://github.com/talos-systems/sidero/commit/c14e055f4c26ee7df38742bb80d385b91e6191b7) chore: bump Talos to 0.11.0-beta.3
* [`7170777`](https://github.com/talos-systems/sidero/commit/7170777d603bc52ef6ba3a02669266d421acc65f) fix: make sure powercycle condition gets properly update
* [`90e7804`](https://github.com/talos-systems/sidero/commit/90e78046409fdf5dd4f16b1e8ae90376d6c813ac) chore: bump dependencies in go.mod
* [`dc70167`](https://github.com/talos-systems/sidero/commit/dc70167c3d11613ecd43291c3930519eb356124c) docs: fix duplicate titles in the guides
* [`fd1fae7`](https://github.com/talos-systems/sidero/commit/fd1fae793148fdc7c0eca22a9e0a82618409644f) chore: update Talos to 0.11.0-beta.2
* [`1f8f141`](https://github.com/talos-systems/sidero/commit/1f8f14178c06d96ed47ca1399340f3ca589c140d) docs: promote 0.3 docs to be the latest
* [`ea3016f`](https://github.com/talos-systems/sidero/commit/ea3016fc19bd20436e70f45206eed08e7c974524) fix: update Sfyra to install CAPI v0.3
* [`8e49ddf`](https://github.com/talos-systems/sidero/commit/8e49ddf49b7d02ed08130037fd4ded7982326ca7) chore: update to latest stable talos providers
* [`1155004`](https://github.com/talos-systems/sidero/commit/115500459856f0181d7cb1936cb9629b97a01d59) docs: document using ISO for mgmt plane
* [`a5b3e67`](https://github.com/talos-systems/sidero/commit/a5b3e67311fca76e1f50139f09fc17393035de18) docs: add docs for server removal/decommissioning
* [`c7ae88a`](https://github.com/talos-systems/sidero/commit/c7ae88a49d73fb247be60f11e37566c2056f96b3) docs: ensure tutorial docs are present in sidebar
* [`83b0875`](https://github.com/talos-systems/sidero/commit/83b08757190081d0ab0205ffea741703328c6f1e) docs: rework guides into Tutorial
* [`ee31160`](https://github.com/talos-systems/sidero/commit/ee31160b5a8b3b5844585bdc9567374345caf14a) fix: make sure components of agent environment are of proper arch
* [`dfe2c85`](https://github.com/talos-systems/sidero/commit/dfe2c85438c064d85d85706c313efc55098b8fe5) chore: bump dependencies via dependabot
* [`bf2be1b`](https://github.com/talos-systems/sidero/commit/bf2be1b1569333034cadd2dacc023538221b2930) feat: update default Talos version to v0.10.3
* [`dfeaeec`](https://github.com/talos-systems/sidero/commit/dfeaeec308e2f45ae815943b343266e1c341a262) docs: add notes on running Sidero on RPi4
</p>
</details>

### Changes since v0.4.0-alpha.1
<details><summary>9 commits</summary>
<p>

* [`af646a7`](https://github.com/talos-systems/sidero/commit/af646a7170f5a5a87d6e604f59d5b66fc3d35fd1) chore: bump Go deps, tools, pkgs, alpine versions
* [`bb52e71`](https://github.com/talos-systems/sidero/commit/bb52e71861072104a5f1a8bc460b01fcf8c84932) feat: support CAPI v1alpha4
* [`d923bb0`](https://github.com/talos-systems/sidero/commit/d923bb0f7f5681ada12159b0ffc518ca40c6056d) docs: add notes on version compatibility and improve the bootstrapping guide
* [`a2bb2d0`](https://github.com/talos-systems/sidero/commit/a2bb2d0919c76b7e051c50f600bab05a9c9d0366) chore: fix docker version in the drone pipeline
* [`e52071d`](https://github.com/talos-systems/sidero/commit/e52071d7f9e5a05dc8b3f9d2ee48265a55c6fc14) fix: shutdown sidero-controller-manager when any component fails
* [`afa114b`](https://github.com/talos-systems/sidero/commit/afa114b18f15798a5c9bbcc7d6c1baa9952ee0af) fix: broken url, ServerClass CR spec
* [`b37f43f`](https://github.com/talos-systems/sidero/commit/b37f43fe6274bcd30ffb9affffefb11335b9d3a2) fix: resource config links
* [`aa2b3f0`](https://github.com/talos-systems/sidero/commit/aa2b3f0e3512ea8049b7fbdd2a423b209126365e) chore: bump dependencies via dependabot
* [`f1c1608`](https://github.com/talos-systems/sidero/commit/f1c16082071cec651a7de562ac101a11a1689db7) chore: bump Talos to 0.11.5
</p>
</details>

### Changes from talos-systems/cluster-api-bootstrap-provider-talos
<details><summary>26 commits</summary>
<p>

* [`2f1364c`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/2f1364c966515d241d21723c44ff3aaab543edfa) release(v0.4.0): prepare release
* [`04742b9`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/04742b96bf757413c88d0f15bee91679644f0337) feat: import fixes/updates from kubeadm bootstrap provider
* [`548b7fb`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/548b7fbd38b89b9790a0daa2380fddb34157cdd5) release(v0.4.0-alpha.0): prepare release
* [`442ee41`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/442ee41bafb2a912e49928c5d61f52c4c61a2593) test: don't set the talosconfig owner ref to the machine
* [`8c7fec8`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/8c7fec8e373bd12609f6274d79ca07d187212d91) fix: don't write incomplete `<cluster>-ca` secret for configtype none
* [`f46c83d`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/f46c83d328ee44db2ccb5eef67b366cc73c13319) feat: bump Talos machinery to 0.12.3
* [`7b760cf`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/7b760cf69ecab93200821dded931171657a5dedc) feat: support CAPI v1alpha4
* [`3147ba4`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/3147ba4fe57b88975133c598c226ff4e397efb44) release(v0.3.0-alpha.1): prepare release
* [`977121a`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/977121ad14dc0637f7c4282e69a4ee26e28372d4) fix: construct properly data secret name
* [`f8c75c8`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/f8c75c89c4653de30165fb1538e906256a4eec66) fix: update metadata.yaml for v0.3 of CABPT
* [`db60f9e`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/db60f9eb0697c4949be9c00cf8dc7787d383bad2) release(v0.3.0-alpha.0): prepare release
* [`755a2dd`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/755a2dd90c3668db89f8eae14f60db4564764475) fix: update Talos machinery to 0.12, fix secrets persistence
* [`f91b032`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/f91b032935776c1224f824cc860bfa4df5e220b1) fix: use bootstrap data secret names
* [`6bff239`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/6bff2393840655c2361def455b601511b86ba71f) chore: use Go 1.17
* [`56fb73b`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/56fb73b53f41b91b12ba2b3c331d7a04b7263a17) test: add test for the second machine
* [`e5b7738`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/e5b773833120fdd7ca4d57e0a0a4fe781495bf7e) test: add more tests
* [`bc4105d`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/bc4105d9e8366d4e840705a6cecfbc81bdcca00a) test: wait for CAPI availability
* [`c82b8ab`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/c82b8ab47bca5313cb96df1b70de0914da285331) chore: make versions configurable
* [`5594c96`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/5594c96daa55fb9fc9af585e8f2fc26551ce9bb5) chore: use codecov uploader from build-container
* [`cced038`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/cced038257d3eec5b7c48bc524de5165b5734496) chore: fix license headers
* [`7b5dc51`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/7b5dc51e83a54a1f5fa707c66a296ca9514c8722) chore: do not run tests on ARM
* [`d6258cf`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/d6258cf21778149a254d9669b03ac10bae9e0955) chore: improve tests runner
* [`c6ce363`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/c6ce36375ef145760647c632d64a9a3c93574e4b) chore: sign Drone CI configuration
* [`ad592d1`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/ad592d16fa8397f88a28e6a4151bc64b0a1c097d) chore: add basic integration test
* [`9fb0d07`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/9fb0d07ca4d2e8333b0b61ee0fe0ba3e6660489f) chore: add missing LICENSE file
* [`acf18d2`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/acf18d2bb09aab64687c1fccf1e628ef76e9cff8) chore: update machinery to v0.11.3
</p>
</details>

### Changes from talos-systems/cluster-api-control-plane-provider-talos
<details><summary>14 commits</summary>
<p>

* [`1bec112`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/1bec1128f39c0abc243a99a92aaf5bf9917439b2) release(v0.3.0): prepare release
* [`1662815`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/166281550865c66ed4f6a015c76c94443b43f0fe) feat: catch up with Kubeadm in terms of conditions
* [`43eb75b`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/43eb75b439d43e87e970df69f49f0abbde047e51) release(v0.3.0-alpha.0): prepare release
* [`48d834b`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/48d834b5dfb364b8e9ae2269771e41a2dc646692) feat: support CAPI v1alpha4
* [`14c6e72`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/14c6e7224703b8838c1790c2014847f991367ff2) release(v0.2.0-alpha.0): prepare release
* [`cd6417d`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/cd6417dd365aa89984703016b61c57e6b33b3b68) fix: update metadata.yaml for v0.2 of CACPPT
* [`8b52b8a`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/8b52b8addd9fa4235c542b0b8554a76f5c76a643) chore: update go to 1.17
* [`86d679a`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/86d679a44e543789474c0b8edaf435a764f7dd2e) chore: update cabpt to v0.3.0
* [`a616f4b`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/a616f4b4bd3b208595cd102eb9e32c8a31b95e18) test: add machine removal test
* [`6ad6aac`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/6ad6aac1315ad5bc8e1264af6162863418cdb280) test: implement scale up and down tests and fix found issues
* [`9435b12`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/9435b1247f010bee00b4a8e4dc592121a0eb2449) chore: add e2e test running on AWS infra
* [`4c7d42c`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/4c7d42caf79ca209f5cda84db2eb712433d3c68b) chore: update bootstrap provider
* [`119b969`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/119b969be2fe152a0e8a63d189563deed55110b4) fix: clean up couple small issues in the etcd member audit code
* [`9be7b88`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/9be7b88bf4a14aec584fe68561c3fda3fbeaf990) chore: update bootstrap provider to stable release
</p>
</details>

### Changes from talos-systems/go-blockdevice
<details><summary>13 commits</summary>
<p>

* [`70d2865`](https://github.com/talos-systems/go-blockdevice/commit/70d28650b398a14469cbb5356417355b0ba62956) fix: try to find cdrom disks
* [`667bf53`](https://github.com/talos-systems/go-blockdevice/commit/667bf539b99ac34b629a0103ef7a7278a5a5f35d) fix: revert gpt partition not found
* [`d7d4cdd`](https://github.com/talos-systems/go-blockdevice/commit/d7d4cdd7ac56c82caab19246b5decd59f12195eb) fix: gpt partition not found
* [`33afba3`](https://github.com/talos-systems/go-blockdevice/commit/33afba347c0dce38a436c46a0aac26d2f99427c1) fix: also open in readonly mode when running `All` lookup method
* [`e367f9d`](https://github.com/talos-systems/go-blockdevice/commit/e367f9dc7fa935f11672de0fdc8a89429285a07a) feat: make probe always open blockdevices in readonly mode
* [`d981156`](https://github.com/talos-systems/go-blockdevice/commit/d9811569588ba44be878a00ce316f59a37abed8b) fix: allow Build for Windows
* [`fe24303`](https://github.com/talos-systems/go-blockdevice/commit/fe2430349e9d734ce6dbf4e7b2e0f8a37bb22679) fix: perform correct PMBR partition calculations
* [`2ec0c3c`](https://github.com/talos-systems/go-blockdevice/commit/2ec0c3cc0ff5ff705ed5c910ca1bcd5d93c7b102) fix: preserve the PMBR bootable flag when opening GPT partition
* [`87816a8`](https://github.com/talos-systems/go-blockdevice/commit/87816a81cefc728cfe3cb221b476d8ed4b609fd8) feat: align partition to minimum I/O size
* [`c34b59f`](https://github.com/talos-systems/go-blockdevice/commit/c34b59fb33a7ad8be18bb19bc8c8d8294b4b3a78) feat: expose more encryption options in the LUKS module
* [`30c2bc3`](https://github.com/talos-systems/go-blockdevice/commit/30c2bc3cb62af52f0aea9ce347923b0649fb7928) feat: mark MBR bootable
* [`1292574`](https://github.com/talos-systems/go-blockdevice/commit/1292574643e06512255fb0f45107e0c296eb5a3b) fix: make disk type matcher parser case insensitive
* [`b77400e`](https://github.com/talos-systems/go-blockdevice/commit/b77400e0a7261bf25da77c1f28c2f393f367bfa9) fix: properly detect nvme and sd card disk types
</p>
</details>

### Changes from talos-systems/go-debug
<details><summary>1 commit</summary>
<p>

* [`3d0a6e1`](https://github.com/talos-systems/go-debug/commit/3d0a6e1bf5e3c521e83ead2c8b7faad3638b8c5d) feat: race build tag flag detector
</p>
</details>

### Changes from talos-systems/go-kmsg
<details><summary>1 commit</summary>
<p>

* [`b08e4d3`](https://github.com/talos-systems/go-kmsg/commit/b08e4d36a2f3df0a3d031b1a3028e2d6e4c26710) feat: replace tab character with space in console output
</p>
</details>

### Changes from talos-systems/go-retry
<details><summary>1 commit</summary>
<p>

* [`c78cc95`](https://github.com/talos-systems/go-retry/commit/c78cc953d9e95992575305b4e8648392c6c9b9e6) fix: implement `errors.Is` for all errors in the set
</p>
</details>

### Dependency Changes

* **github.com/evanphx/json-patch**                                      v4.9.0 -> v4.11.0
* **github.com/onsi/ginkgo**                                             v1.16.3 -> v1.16.4
* **github.com/onsi/gomega**                                             v1.13.0 -> v1.16.0
* **github.com/talos-systems/cluster-api-bootstrap-provider-talos**      v0.2.0 -> v0.4.0
* **github.com/talos-systems/cluster-api-control-plane-provider-talos**  v0.1.0 -> v0.3.0
* **github.com/talos-systems/go-blockdevice**                            1d830a25f64f -> v0.2.4
* **github.com/talos-systems/go-debug**                                  v0.2.0 -> v0.2.1
* **github.com/talos-systems/go-kmsg**                                   v0.1.0 -> v0.1.1
* **github.com/talos-systems/go-retry**                                  v0.3.0 -> v0.3.1
* **github.com/talos-systems/talos/pkg/machinery**                       828772cec9a3 -> v0.13.0
* **golang.org/x/net**                                                   0714010a04ed -> 853a461950ff
* **golang.org/x/sys**                                                   0981d6026fa6 -> 39ccf1dd6fa6
* **google.golang.org/grpc**                                             v1.38.0 -> v1.41.0
* **google.golang.org/protobuf**                                         v1.26.0 -> v1.27.1
* **k8s.io/api**                                                         v0.19.3 -> v0.22.2
* **k8s.io/apiextensions-apiserver**                                     v0.19.1 -> v0.22.2
* **k8s.io/apimachinery**                                                v0.19.3 -> v0.22.2
* **k8s.io/client-go**                                                   v0.19.3 -> v0.22.2
* **k8s.io/utils**                                                       67b214c5f920 -> bdf08cb9a70a
* **sigs.k8s.io/cluster-api**                                            v0.3.12 -> v0.4.3
* **sigs.k8s.io/controller-runtime**                                     v0.6.3 -> v0.9.7

Previous release can be found at [v0.3.2](https://github.com/talos-systems/sidero/releases/tag/v0.3.2)


## [Sidero 0.4.0-alpha.1](https://github.com/talos-systems/sidero/releases/tag/v0.4.0-alpha.1) (2021-08-04)

Welcome to the v0.4.0-alpha.1 release of Sidero!  
*This is a pre-release of Sidero*



Please try out the release binaries and report any issues at
https://github.com/talos-systems/sidero/issues.

### BMC Port

Sidero now supports the ability to specify the port in a server's BMC info. By default, this value will be determined by talking directly to the BMC if possible, with a fallback to port 623. The value can also simply be specied as part of editing the Server resource directly.


### Contributors

* Andrey Smirnov
* Spencer Smith
* Artem Chernyshev
* Khue Doan
* Seán C McCord

### Changes
<details><summary>21 commits</summary>
<p>

* [`7bdee0f`](https://github.com/talos-systems/sidero/commit/7bdee0f63fef443f4d51329f40c2db8429030008) fix: update sidero IPMI user to work properly on idrac
* [`acd82e0`](https://github.com/talos-systems/sidero/commit/acd82e0611843d3c273824d7daad53b3abed3352) docs: bump clusterctl from v0.3.14 to v0.3.21 in /docs/website
* [`fb0da3c`](https://github.com/talos-systems/sidero/commit/fb0da3cf5087f67cd6dc40381f228c756abedd4b) release(v0.4.0-alpha.0): prepare release
* [`ee36c74`](https://github.com/talos-systems/sidero/commit/ee36c745016a324220b5d17b7ef5fff9d2ce85a8) docs: redirect latest to v0.3
* [`7cdae7b`](https://github.com/talos-systems/sidero/commit/7cdae7b616db7720627ab2e5b409e2f22b0085e0) feat: add ability to specify ports for BMC
* [`c14e055`](https://github.com/talos-systems/sidero/commit/c14e055f4c26ee7df38742bb80d385b91e6191b7) chore: bump Talos to 0.11.0-beta.3
* [`7170777`](https://github.com/talos-systems/sidero/commit/7170777d603bc52ef6ba3a02669266d421acc65f) fix: make sure powercycle condition gets properly update
* [`90e7804`](https://github.com/talos-systems/sidero/commit/90e78046409fdf5dd4f16b1e8ae90376d6c813ac) chore: bump dependencies in go.mod
* [`dc70167`](https://github.com/talos-systems/sidero/commit/dc70167c3d11613ecd43291c3930519eb356124c) docs: fix duplicate titles in the guides
* [`fd1fae7`](https://github.com/talos-systems/sidero/commit/fd1fae793148fdc7c0eca22a9e0a82618409644f) chore: update Talos to 0.11.0-beta.2
* [`1f8f141`](https://github.com/talos-systems/sidero/commit/1f8f14178c06d96ed47ca1399340f3ca589c140d) docs: promote 0.3 docs to be the latest
* [`ea3016f`](https://github.com/talos-systems/sidero/commit/ea3016fc19bd20436e70f45206eed08e7c974524) fix: update Sfyra to install CAPI v0.3
* [`8e49ddf`](https://github.com/talos-systems/sidero/commit/8e49ddf49b7d02ed08130037fd4ded7982326ca7) chore: update to latest stable talos providers
* [`1155004`](https://github.com/talos-systems/sidero/commit/115500459856f0181d7cb1936cb9629b97a01d59) docs: document using ISO for mgmt plane
* [`a5b3e67`](https://github.com/talos-systems/sidero/commit/a5b3e67311fca76e1f50139f09fc17393035de18) docs: add docs for server removal/decommissioning
* [`c7ae88a`](https://github.com/talos-systems/sidero/commit/c7ae88a49d73fb247be60f11e37566c2056f96b3) docs: ensure tutorial docs are present in sidebar
* [`83b0875`](https://github.com/talos-systems/sidero/commit/83b08757190081d0ab0205ffea741703328c6f1e) docs: rework guides into Tutorial
* [`ee31160`](https://github.com/talos-systems/sidero/commit/ee31160b5a8b3b5844585bdc9567374345caf14a) fix: make sure components of agent environment are of proper arch
* [`dfe2c85`](https://github.com/talos-systems/sidero/commit/dfe2c85438c064d85d85706c313efc55098b8fe5) chore: bump dependencies via dependabot
* [`bf2be1b`](https://github.com/talos-systems/sidero/commit/bf2be1b1569333034cadd2dacc023538221b2930) feat: update default Talos version to v0.10.3
* [`dfeaeec`](https://github.com/talos-systems/sidero/commit/dfeaeec308e2f45ae815943b343266e1c341a262) docs: add notes on running Sidero on RPi4
</p>
</details>

### Changes since v0.4.0-alpha.0
<details><summary>2 commits</summary>
<p>

* [`7bdee0f`](https://github.com/talos-systems/sidero/commit/7bdee0f63fef443f4d51329f40c2db8429030008) fix: update sidero IPMI user to work properly on idrac
* [`acd82e0`](https://github.com/talos-systems/sidero/commit/acd82e0611843d3c273824d7daad53b3abed3352) docs: bump clusterctl from v0.3.14 to v0.3.21 in /docs/website
</p>
</details>

### Changes from talos-systems/go-blockdevice
<details><summary>3 commits</summary>
<p>

* [`30c2bc3`](https://github.com/talos-systems/go-blockdevice/commit/30c2bc3cb62af52f0aea9ce347923b0649fb7928) feat: mark MBR bootable
* [`1292574`](https://github.com/talos-systems/go-blockdevice/commit/1292574643e06512255fb0f45107e0c296eb5a3b) fix: make disk type matcher parser case insensitive
* [`b77400e`](https://github.com/talos-systems/go-blockdevice/commit/b77400e0a7261bf25da77c1f28c2f393f367bfa9) fix: properly detect nvme and sd card disk types
</p>
</details>

### Changes from talos-systems/go-debug
<details><summary>1 commit</summary>
<p>

* [`3d0a6e1`](https://github.com/talos-systems/go-debug/commit/3d0a6e1bf5e3c521e83ead2c8b7faad3638b8c5d) feat: race build tag flag detector
</p>
</details>

### Changes from talos-systems/go-kmsg
<details><summary>1 commit</summary>
<p>

* [`b08e4d3`](https://github.com/talos-systems/go-kmsg/commit/b08e4d36a2f3df0a3d031b1a3028e2d6e4c26710) feat: replace tab character with space in console output
</p>
</details>

### Changes from talos-systems/go-retry
<details><summary>1 commit</summary>
<p>

* [`c78cc95`](https://github.com/talos-systems/go-retry/commit/c78cc953d9e95992575305b4e8648392c6c9b9e6) fix: implement `errors.Is` for all errors in the set
</p>
</details>

### Dependency Changes

* **github.com/evanphx/json-patch**                 v4.9.0 -> v4.11.0
* **github.com/onsi/ginkgo**                        v1.16.3 -> v1.16.4
* **github.com/talos-systems/go-blockdevice**       1d830a25f64f -> v0.2.1
* **github.com/talos-systems/go-debug**             v0.2.0 -> v0.2.1
* **github.com/talos-systems/go-kmsg**              v0.1.0 -> v0.1.1
* **github.com/talos-systems/go-retry**             v0.3.0 -> v0.3.1
* **github.com/talos-systems/talos/pkg/machinery**  828772cec9a3 -> v0.11.0-beta.3
* **golang.org/x/net**                              0714010a04ed -> abc453219eb5
* **golang.org/x/sys**                              0981d6026fa6 -> 59db8d763f22
* **google.golang.org/grpc**                        v1.38.0 -> v1.39.0
* **google.golang.org/protobuf**                    v1.26.0 -> v1.27.1
* **k8s.io/api**                                    v0.19.3 -> v0.17.9
* **k8s.io/apiextensions-apiserver**                v0.19.1 -> v0.17.9
* **k8s.io/apimachinery**                           v0.19.3 -> v0.17.9
* **k8s.io/client-go**                              v0.19.3 -> v0.17.9
* **sigs.k8s.io/cluster-api**                       v0.3.12 -> v0.3.20
* **sigs.k8s.io/controller-runtime**                v0.6.3 -> v0.5.14

Previous release can be found at [v0.3.0](https://github.com/talos-systems/sidero/releases/tag/v0.3.0)

## [Sidero 0.4.0-alpha.0](https://github.com/talos-systems/sidero/releases/tag/v0.4.0-alpha.0) (2021-07-19)

Welcome to the v0.4.0-alpha.0 release of Sidero!  
*This is a pre-release of Sidero*



Please try out the release binaries and report any issues at
https://github.com/talos-systems/sidero/issues.

### BMC Port

Sidero now supports the ability to specify the port in a server's BMC info. By default, this value will be determined by talking directly to the BMC if possible, with a fallback to port 623. The value can also simply be specied as part of editing the Server resource directly.


### Contributors

* Andrey Smirnov
* Spencer Smith
* Artem Chernyshev
* Seán C McCord

### Changes
<details><summary>18 commits</summary>
<p>

* [`ee36c74`](https://github.com/talos-systems/sidero/commit/ee36c745016a324220b5d17b7ef5fff9d2ce85a8) docs: redirect latest to v0.3
* [`7cdae7b`](https://github.com/talos-systems/sidero/commit/7cdae7b616db7720627ab2e5b409e2f22b0085e0) feat: add ability to specify ports for BMC
* [`c14e055`](https://github.com/talos-systems/sidero/commit/c14e055f4c26ee7df38742bb80d385b91e6191b7) chore: bump Talos to 0.11.0-beta.3
* [`7170777`](https://github.com/talos-systems/sidero/commit/7170777d603bc52ef6ba3a02669266d421acc65f) fix: make sure powercycle condition gets properly update
* [`90e7804`](https://github.com/talos-systems/sidero/commit/90e78046409fdf5dd4f16b1e8ae90376d6c813ac) chore: bump dependencies in go.mod
* [`dc70167`](https://github.com/talos-systems/sidero/commit/dc70167c3d11613ecd43291c3930519eb356124c) docs: fix duplicate titles in the guides
* [`fd1fae7`](https://github.com/talos-systems/sidero/commit/fd1fae793148fdc7c0eca22a9e0a82618409644f) chore: update Talos to 0.11.0-beta.2
* [`1f8f141`](https://github.com/talos-systems/sidero/commit/1f8f14178c06d96ed47ca1399340f3ca589c140d) docs: promote 0.3 docs to be the latest
* [`ea3016f`](https://github.com/talos-systems/sidero/commit/ea3016fc19bd20436e70f45206eed08e7c974524) fix: update Sfyra to install CAPI v0.3
* [`8e49ddf`](https://github.com/talos-systems/sidero/commit/8e49ddf49b7d02ed08130037fd4ded7982326ca7) chore: update to latest stable talos providers
* [`1155004`](https://github.com/talos-systems/sidero/commit/115500459856f0181d7cb1936cb9629b97a01d59) docs: document using ISO for mgmt plane
* [`a5b3e67`](https://github.com/talos-systems/sidero/commit/a5b3e67311fca76e1f50139f09fc17393035de18) docs: add docs for server removal/decommissioning
* [`c7ae88a`](https://github.com/talos-systems/sidero/commit/c7ae88a49d73fb247be60f11e37566c2056f96b3) docs: ensure tutorial docs are present in sidebar
* [`83b0875`](https://github.com/talos-systems/sidero/commit/83b08757190081d0ab0205ffea741703328c6f1e) docs: rework guides into Tutorial
* [`ee31160`](https://github.com/talos-systems/sidero/commit/ee31160b5a8b3b5844585bdc9567374345caf14a) fix: make sure components of agent environment are of proper arch
* [`dfe2c85`](https://github.com/talos-systems/sidero/commit/dfe2c85438c064d85d85706c313efc55098b8fe5) chore: bump dependencies via dependabot
* [`bf2be1b`](https://github.com/talos-systems/sidero/commit/bf2be1b1569333034cadd2dacc023538221b2930) feat: update default Talos version to v0.10.3
* [`dfeaeec`](https://github.com/talos-systems/sidero/commit/dfeaeec308e2f45ae815943b343266e1c341a262) docs: add notes on running Sidero on RPi4
</p>
</details>

### Changes from talos-systems/go-blockdevice
<details><summary>3 commits</summary>
<p>

* [`30c2bc3`](https://github.com/talos-systems/go-blockdevice/commit/30c2bc3cb62af52f0aea9ce347923b0649fb7928) feat: mark MBR bootable
* [`1292574`](https://github.com/talos-systems/go-blockdevice/commit/1292574643e06512255fb0f45107e0c296eb5a3b) fix: make disk type matcher parser case insensitive
* [`b77400e`](https://github.com/talos-systems/go-blockdevice/commit/b77400e0a7261bf25da77c1f28c2f393f367bfa9) fix: properly detect nvme and sd card disk types
</p>
</details>

### Changes from talos-systems/go-debug
<details><summary>1 commit</summary>
<p>

* [`3d0a6e1`](https://github.com/talos-systems/go-debug/commit/3d0a6e1bf5e3c521e83ead2c8b7faad3638b8c5d) feat: race build tag flag detector
</p>
</details>

### Changes from talos-systems/go-kmsg
<details><summary>1 commit</summary>
<p>

* [`b08e4d3`](https://github.com/talos-systems/go-kmsg/commit/b08e4d36a2f3df0a3d031b1a3028e2d6e4c26710) feat: replace tab character with space in console output
</p>
</details>

### Changes from talos-systems/go-retry
<details><summary>1 commit</summary>
<p>

* [`c78cc95`](https://github.com/talos-systems/go-retry/commit/c78cc953d9e95992575305b4e8648392c6c9b9e6) fix: implement `errors.Is` for all errors in the set
</p>
</details>

### Dependency Changes

* **github.com/evanphx/json-patch**                 v4.9.0 -> v4.11.0
* **github.com/onsi/ginkgo**                        v1.16.3 -> v1.16.4
* **github.com/talos-systems/go-blockdevice**       1d830a25f64f -> v0.2.1
* **github.com/talos-systems/go-debug**             v0.2.0 -> v0.2.1
* **github.com/talos-systems/go-kmsg**              v0.1.0 -> v0.1.1
* **github.com/talos-systems/go-retry**             v0.3.0 -> v0.3.1
* **github.com/talos-systems/talos/pkg/machinery**  828772cec9a3 -> v0.11.0-beta.3
* **golang.org/x/net**                              0714010a04ed -> abc453219eb5
* **golang.org/x/sys**                              0981d6026fa6 -> 59db8d763f22
* **google.golang.org/grpc**                        v1.38.0 -> v1.39.0
* **google.golang.org/protobuf**                    v1.26.0 -> v1.27.1
* **k8s.io/api**                                    v0.19.3 -> v0.17.9
* **k8s.io/apiextensions-apiserver**                v0.19.1 -> v0.17.9
* **k8s.io/apimachinery**                           v0.19.3 -> v0.17.9
* **k8s.io/client-go**                              v0.19.3 -> v0.17.9
* **sigs.k8s.io/cluster-api**                       v0.3.12 -> v0.3.20
* **sigs.k8s.io/controller-runtime**                v0.6.3 -> v0.5.14

Previous release can be found at [v0.3.0](https://github.com/talos-systems/sidero/releases/tag/v0.3.0)

