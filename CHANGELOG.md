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

