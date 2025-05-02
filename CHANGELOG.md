## [Sidero 0.6.9](https://github.com/talos-systems/sidero/releases/tag/v0.6.9) (2025-05-02)

Welcome to the v0.6.9 release of Sidero!



Please try out the release binaries and report any issues at
https://github.com/talos-systems/sidero/issues.

### Contributors

* Noel Georgi

### Changes
<details><summary>1 commit</summary>
<p>

* [`4e082be`](https://github.com/talos-systems/sidero/commit/4e082be6d1277372479d0820ae24ad7ba978c519) fix(ci): arm64 container image asset
</p>
</details>

### Dependency Changes

This release has no dependency changes

Previous release can be found at [v0.6.8](https://github.com/talos-systems/sidero/releases/tag/v0.6.8)

## [Sidero 0.6.8](https://github.com/talos-systems/sidero/releases/tag/v0.6.8) (2025-05-01)

Welcome to the v0.6.8 release of Sidero!



Please try out the release binaries and report any issues at
https://github.com/talos-systems/sidero/issues.

### Cluster API

Sidero Metal is now built and tested with Cluster API 1.10.0.


### Talos Linux

Sidero Metal now supports Talos Linux v1.10.x.


### Contributors

* Andrey Smirnov
* Noel Georgi
* Utku Ozdemir

### Changes
<details><summary>4 commits</summary>
<p>

* [`520bb35`](https://github.com/talos-systems/sidero/commit/520bb357a07b07b1aff1f29c2aa5281cd903fcb6) feat: update to Talos 1.10.0 final
* [`3784341`](https://github.com/talos-systems/sidero/commit/37843417f965ce88c0873e530976e8e2e3bab662) feat: update Talos and CAPI to 1.10
* [`b27fca4`](https://github.com/talos-systems/sidero/commit/b27fca4f564d25ba831f5036e1c27f5aa4573004) feat: use kres to manage github actions
* [`a46f257`](https://github.com/talos-systems/sidero/commit/a46f2576f9b73b5304bc586b94c3dcc6ce0d6f74) docs: add Omni bare metal infra provider link to maintenance warning
</p>
</details>

### Dependency Changes

* **github.com/google/go-cmp**                   v0.6.0 -> v0.7.0
* **github.com/insomniacslk/dhcp**               a481575ed0ef -> 5f8cf70e8c5f
* **github.com/siderolabs/gen**                  v0.7.0 -> v0.8.0
* **github.com/siderolabs/go-debug**             v0.4.0 -> v0.5.0
* **github.com/siderolabs/go-pointer**           v1.0.0 -> v1.0.1
* **github.com/siderolabs/siderolink**           v0.3.11 -> v0.3.14
* **github.com/siderolabs/talos/pkg/machinery**  v1.9.0 -> v1.10.0
* **github.com/spf13/pflag**                     v1.0.5 -> v1.0.6
* **golang.org/x/net**                           v0.32.0 -> v0.39.0
* **golang.org/x/sync**                          v0.10.0 -> v0.13.0
* **golang.org/x/sys**                           v0.28.0 -> v0.32.0
* **golang.zx2c4.com/wireguard/wgctrl**          925a1e7659e6 -> a9ab2273dd10
* **google.golang.org/grpc**                     v1.68.1 -> v1.72.0
* **google.golang.org/protobuf**                 v1.35.2 -> v1.36.6
* **k8s.io/api**                                 v0.31.3 -> v0.32.3
* **k8s.io/apiextensions-apiserver**             v0.31.3 -> v0.32.3
* **k8s.io/apimachinery**                        v0.31.3 -> v0.32.3
* **k8s.io/client-go**                           v0.31.3 -> v0.32.3
* **k8s.io/component-base**                      v0.31.3 -> v0.32.3
* **sigs.k8s.io/cluster-api**                    v1.9.1 -> v1.10.1
* **sigs.k8s.io/controller-runtime**             v0.19.3 -> v0.20.4

Previous release can be found at [v0.6.7](https://github.com/talos-systems/sidero/releases/tag/v0.6.7)

## [Sidero 0.6.7](https://github.com/talos-systems/sidero/releases/tag/v0.6.7) (2024-12-20)

Welcome to the v0.6.7 release of Sidero!



Please try out the release binaries and report any issues at
https://github.com/talos-systems/sidero/issues.

### Cluster API

Sidero Metal is now built and tested with Cluster API 1.9.0.


### Talos Linux

Sidero Metal now supports Talos Linux v1.9.x.


### Contributors

* Andrey Smirnov

### Changes
<details><summary>2 commits</summary>
<p>

* [`de84f7e`](https://github.com/talos-systems/sidero/commit/de84f7e54254da74c3f5b3f659af97032e317617) feat: update to Talos 1.9.0
* [`b8daa0c`](https://github.com/talos-systems/sidero/commit/b8daa0cd2ce3f6f3c1194d5b3492257ae9651bb8) feat: update dependencies for Talos 1.9.0-beta.0
</p>
</details>

### Dependency Changes

* **github.com/insomniacslk/dhcp**               a3a4c1f04475 -> a481575ed0ef
* **github.com/siderolabs/gen**                  v0.5.0 -> v0.7.0
* **github.com/siderolabs/go-cmd**               v0.1.1 -> v0.1.3
* **github.com/siderolabs/go-debug**             v0.3.0 -> v0.4.0
* **github.com/siderolabs/go-smbios**            v0.3.2 -> v0.3.3
* **github.com/siderolabs/grpc-proxy**           v0.4.1 -> v0.5.1
* **github.com/siderolabs/siderolink**           v0.3.9 -> v0.3.11
* **github.com/siderolabs/talos/pkg/machinery**  v1.8.0 -> v1.9.0
* **github.com/stretchr/testify**                v1.9.0 -> v1.10.0
* **golang.org/x/net**                           v0.29.0 -> v0.32.0
* **golang.org/x/sync**                          v0.8.0 -> v0.10.0
* **golang.org/x/sys**                           v0.25.0 -> v0.28.0
* **google.golang.org/grpc**                     v1.66.0 -> v1.68.1
* **google.golang.org/protobuf**                 v1.34.2 -> v1.35.2
* **k8s.io/api**                                 v0.31.0 -> v0.31.3
* **k8s.io/apiextensions-apiserver**             v0.31.0 -> v0.31.3
* **k8s.io/client-go**                           v0.31.0 -> v0.31.3
* **k8s.io/component-base**                      v0.31.0 -> v0.31.3
* **sigs.k8s.io/cluster-api**                    v1.8.3 -> v1.9.1
* **sigs.k8s.io/controller-runtime**             v0.19.0 -> v0.19.3

Previous release can be found at [v0.6.6](https://github.com/talos-systems/sidero/releases/tag/v0.6.6)

## [Sidero 0.6.6](https://github.com/talos-systems/sidero/releases/tag/v0.6.6) (2024-09-24)

Welcome to the v0.6.6 release of Sidero!



Please try out the release binaries and report any issues at
https://github.com/talos-systems/sidero/issues.

### Talos Linux

Sidero Metal now supports Talos Linux v1.8.x.


### Contributors

* Andrey Smirnov
* Spencer Smith
* Steve Francis

### Changes
<details><summary>6 commits</summary>
<p>

* [`b6fe15d`](https://github.com/talos-systems/sidero/commit/b6fe15df7e3907b5a8814b55fc690d450d8afb2a) feat: update for Talos 1.8
* [`9a35d37`](https://github.com/talos-systems/sidero/commit/9a35d37cf6a221da96195fd96326d50b9de4e008) docs: add deprecation notice
* [`d439b10`](https://github.com/talos-systems/sidero/commit/d439b10b6023910489347ef7ad980d2849f9c18b) release(v0.6.5): prepare release
* [`a30f4d9`](https://github.com/talos-systems/sidero/commit/a30f4d9bbacc7f7755d779c29fc6dc93bd6f3aca) feat: provide negative address filter support
* [`cb907cf`](https://github.com/talos-systems/sidero/commit/cb907cf6e05c3ae0dbd939da5a02e0242af0ff80) chore: update office hours
* [`d91ab75`](https://github.com/talos-systems/sidero/commit/d91ab75b0771109d5c280a95098c244a2d28d6c4) fix: update HTTP boot filename
</p>
</details>

### Changes since v0.6.5
<details><summary>2 commits</summary>
<p>

* [`b6fe15d`](https://github.com/talos-systems/sidero/commit/b6fe15df7e3907b5a8814b55fc690d450d8afb2a) feat: update for Talos 1.8
* [`9a35d37`](https://github.com/talos-systems/sidero/commit/9a35d37cf6a221da96195fd96326d50b9de4e008) docs: add deprecation notice
</p>
</details>

### Dependency Changes

* **github.com/go-logr/logr**                    v1.4.1 -> v1.4.2
* **github.com/insomniacslk/dhcp**               c728f5dd21c8 -> a3a4c1f04475
* **github.com/jsimonetti/rtnetlink**            v1.4.1 -> v1.4.2
* **github.com/siderolabs/gen**                  v0.4.8 -> v0.5.0
* **github.com/siderolabs/go-blockdevice**       v0.4.7 -> v0.4.8
* **github.com/siderolabs/grpc-proxy**           v0.4.0 -> v0.4.1
* **github.com/siderolabs/siderolink**           v0.3.5 -> v0.3.9
* **github.com/siderolabs/talos/pkg/machinery**  v1.7.0 -> v1.8.0-beta.0
* **golang.org/x/net**                           v0.24.0 -> v0.29.0
* **golang.org/x/sync**                          v0.7.0 -> v0.8.0
* **golang.org/x/sys**                           v0.19.0 -> v0.25.0
* **google.golang.org/grpc**                     v1.62.1 -> v1.66.0
* **google.golang.org/protobuf**                 v1.33.0 -> v1.34.2
* **k8s.io/api**                                 v0.29.3 -> v0.31.0
* **k8s.io/apiextensions-apiserver**             v0.29.3 -> v0.31.0
* **k8s.io/apimachinery**                        v0.29.3 -> v0.31.0
* **k8s.io/client-go**                           v0.29.3 -> v0.31.0
* **k8s.io/component-base**                      v0.29.3 -> v0.31.0
* **k8s.io/klog/v2**                             v2.110.1 -> v2.130.1
* **sigs.k8s.io/cluster-api**                    v1.7.0 -> v1.8.3
* **sigs.k8s.io/controller-runtime**             v0.17.3 -> v0.19.0

Previous release can be found at [v0.6.4](https://github.com/talos-systems/sidero/releases/tag/v0.6.4)

## [Sidero 0.6.5](https://github.com/talos-systems/sidero/releases/tag/v0.6.5) (2024-05-23)

Welcome to the v0.6.5 release of Sidero!



Please try out the release binaries and report any issues at
https://github.com/talos-systems/sidero/issues.

### Contributors

* Andrey Smirnov
* Spencer Smith

### Changes
<details><summary>3 commits</summary>
<p>

* [`a30f4d9`](https://github.com/talos-systems/sidero/commit/a30f4d9bbacc7f7755d779c29fc6dc93bd6f3aca) feat: provide negative address filter support
* [`cb907cf`](https://github.com/talos-systems/sidero/commit/cb907cf6e05c3ae0dbd939da5a02e0242af0ff80) chore: update office hours
* [`d91ab75`](https://github.com/talos-systems/sidero/commit/d91ab75b0771109d5c280a95098c244a2d28d6c4) fix: update HTTP boot filename
</p>
</details>

### Dependency Changes

* **github.com/insomniacslk/dhcp**               c728f5dd21c8 -> f1cffa2c0c49
* **github.com/jsimonetti/rtnetlink**            v1.4.1 -> v1.4.2
* **github.com/siderolabs/talos/pkg/machinery**  v1.7.0 -> v1.7.2
* **golang.org/x/sys**                           v0.19.0 -> v0.20.0
* **sigs.k8s.io/cluster-api**                    v1.7.0 -> v1.7.2

Previous release can be found at [v0.6.4](https://github.com/talos-systems/sidero/releases/tag/v0.6.4)

## [Sidero 0.6.4](https://github.com/talos-systems/sidero/releases/tag/v0.6.4) (2024-04-19)

Welcome to the v0.6.4 release of Sidero!



Please try out the release binaries and report any issues at
https://github.com/talos-systems/sidero/issues.

### Patches

Sidero Metal now supports Talos Linux machine configuration strategic merge patches via 'strategicPatches' field on the `Server` and `ServerClass` CRDs.


### Contributors

* Andrey Smirnov
* Andrew Rynhard
* Ksawery Kuczyński
* Luke Carrier

### Changes
<details><summary>5 commits</summary>
<p>

* [`62d34bd`](https://github.com/talos-systems/sidero/commit/62d34bd29b26380a6b2e94f9b534ed2fbb5b70c7) feat: update to final Talos 1.7.0
* [`b3f0131`](https://github.com/talos-systems/sidero/commit/b3f01313937cb5ae4c52a6815f32939c61218f8e) feat: add support for strategic merge patches
* [`5f9acdf`](https://github.com/talos-systems/sidero/commit/5f9acdf5648dccdb23d64186ba444e2251a7ea8c) feat: update to Talos 1.7.0-beta.1
* [`b19e58a`](https://github.com/talos-systems/sidero/commit/b19e58a85aa2d422d7bec354f72893fabb7519f4) chore: add notice to README
* [`376fd6e`](https://github.com/talos-systems/sidero/commit/376fd6e6bf441e08491ef164efb6a2f129cf0982) docs: correct "which { much => must }" typo
</p>
</details>

### Dependency Changes

* **github.com/go-logr/logr**                    v1.3.0 -> v1.4.1
* **github.com/insomniacslk/dhcp**               8c70d406f6d2 -> c728f5dd21c8
* **github.com/jsimonetti/rtnetlink**            v1.4.0 -> v1.4.1
* **github.com/siderolabs/gen**                  v0.4.7 -> v0.4.8
* **github.com/siderolabs/go-debug**             v0.2.3 -> v0.3.0
* **github.com/siderolabs/go-kmsg**              v0.1.3 -> v0.1.4
* **github.com/siderolabs/go-procfs**            v0.1.1 -> v0.1.2
* **github.com/siderolabs/siderolink**           v0.3.4 -> v0.3.5
* **github.com/siderolabs/talos/pkg/machinery**  v1.6.0 -> v1.7.0
* **github.com/stretchr/testify**                v1.8.4 -> v1.9.0
* **go.uber.org/zap**                            v1.26.0 -> v1.27.0
* **golang.org/x/net**                           v0.18.0 -> v0.24.0
* **golang.org/x/sync**                          v0.5.0 -> v0.7.0
* **golang.org/x/sys**                           e4099bfacb8c -> v0.19.0
* **google.golang.org/grpc**                     v1.59.0 -> v1.62.1
* **google.golang.org/protobuf**                 v1.31.0 -> v1.33.0
* **k8s.io/api**                                 v0.28.4 -> v0.29.3
* **k8s.io/apiextensions-apiserver**             v0.28.4 -> v0.29.3
* **k8s.io/apimachinery**                        v0.28.4 -> v0.29.3
* **k8s.io/client-go**                           v0.28.4 -> v0.29.3
* **k8s.io/component-base**                      v0.28.4 -> v0.29.3
* **k8s.io/klog/v2**                             v2.110.1 **_new_**
* **sigs.k8s.io/cluster-api**                    v1.6.0 -> v1.7.0
* **sigs.k8s.io/controller-runtime**             v0.16.3 -> v0.17.3

Previous release can be found at [v0.6.3](https://github.com/talos-systems/sidero/releases/tag/v0.6.3)

## [Sidero 0.6.3](https://github.com/talos-systems/sidero/releases/tag/v0.6.3) (2024-01-23)

Welcome to the v0.6.3 release of Sidero!



Please try out the release binaries and report any issues at
https://github.com/talos-systems/sidero/issues.

### Contributors

* Alexandra
* Andrey Smirnov
* Anthony ARNAUD
* rmvangun

### Changes
<details><summary>4 commits</summary>
<p>

* [`30d077f`](https://github.com/talos-systems/sidero/commit/30d077fb50a2b4200d19ab7a07beddc172ce0ce5) fix: set a default controller runtime log
* [`aabe7bc`](https://github.com/talos-systems/sidero/commit/aabe7bc6945328b1b23b0aded7dbbe39b0c01e02) fix: sidero.endpoint and sidero.mac was inverted
* [`90d7a7d`](https://github.com/talos-systems/sidero/commit/90d7a7d5b3bb42c1515bd421f3791da191183ec8) docs: update to reflect master => controlplane
* [`dc8b95e`](https://github.com/talos-systems/sidero/commit/dc8b95edf73b59efea6a0d7b29e640e2bbd06c2e) docs: fix wrong raspberry installation url
</p>
</details>

### Dependency Changes

This release has no dependency changes

Previous release can be found at [v0.6.2](https://github.com/talos-systems/sidero/releases/tag/v0.6.2)

## [Sidero 0.6.2](https://github.com/talos-systems/sidero/releases/tag/v0.6.2) (2023-12-15)

Welcome to the v0.6.2 release of Sidero!



Please try out the release binaries and report any issues at
https://github.com/talos-systems/sidero/issues.

### Contributors

* Nathan Lee
* Andrey Smirnov
* Oscar Utbult

### Changes
<details><summary>4 commits</summary>
<p>

* [`01981eb`](https://github.com/talos-systems/sidero/commit/01981eb1c63db76e5aa8ab45d366ad028c683e9b) feat: update Talos to 1.6.0
* [`130e608`](https://github.com/talos-systems/sidero/commit/130e608a0d3e66bf9501e26ce1bda7a7f844b3ee) feat: add firmware to agent initramfs for QLogic NICs
* [`0f7973c`](https://github.com/talos-systems/sidero/commit/0f7973c86c174a8d859def174f41d9c708406eba) docs: fix website redirect to latest version
* [`7f6f787`](https://github.com/talos-systems/sidero/commit/7f6f7875054a5c04edb43407e58ffd0647f47ed3) fix: handle network interfaces without dhcp during pxeboot
</p>
</details>

### Dependency Changes

* **github.com/go-logr/logr**                    v1.2.4 -> v1.3.0
* **github.com/google/go-cmp**                   v0.5.9 -> v0.6.0
* **github.com/insomniacslk/dhcp**               65c27093e38a -> 8c70d406f6d2
* **github.com/jsimonetti/rtnetlink**            v1.3.5 -> v1.4.0
* **github.com/siderolabs/gen**                  v0.4.5 -> v0.4.7
* **github.com/siderolabs/go-blockdevice**       v0.4.6 -> v0.4.7
* **github.com/siderolabs/go-pointer**           v1.0.0 **_new_**
* **github.com/siderolabs/go-retry**             v0.3.2 -> v0.3.3
* **github.com/siderolabs/siderolink**           v0.3.1 -> v0.3.4
* **github.com/siderolabs/talos/pkg/machinery**  v1.5.3 -> v1.6.0
* **github.com/spf13/pflag**                     v1.0.5 **_new_**
* **golang.org/x/net**                           v0.17.0 -> v0.18.0
* **golang.org/x/sync**                          v0.4.0 -> v0.5.0
* **golang.org/x/sys**                           v0.13.0 -> e4099bfacb8c
* **google.golang.org/grpc**                     v1.58.3 -> v1.59.0
* **k8s.io/api**                                 v0.27.2 -> v0.28.4
* **k8s.io/apiextensions-apiserver**             v0.27.2 -> v0.28.4
* **k8s.io/apimachinery**                        v0.27.2 -> v0.28.4
* **k8s.io/client-go**                           v0.27.2 -> v0.28.4
* **k8s.io/component-base**                      v0.28.4 **_new_**
* **k8s.io/utils**                               a36077c30491 -> d93618cff8a2
* **sigs.k8s.io/cluster-api**                    v1.5.2 -> v1.6.0
* **sigs.k8s.io/controller-runtime**             v0.15.1 -> v0.16.3

Previous release can be found at [v0.6.1](https://github.com/talos-systems/sidero/releases/tag/v0.6.1)


## [Sidero 0.6.1](https://github.com/talos-systems/sidero/releases/tag/v0.6.1) (2023-10-12)

Welcome to the v0.6.1 release of Sidero!



Please try out the release binaries and report any issues at
https://github.com/talos-systems/sidero/issues.

### DHCP Proxy

Sidero Controller Manager now provides a way to disable DHCP Proxy with `SIDERO_CONTROLLER_MANAGER_DISABLE_DHCP_PROXY` variable on
installation.


### Contributors

* Andrey Smirnov
* Daniel Höxtermann
* Victor Seva

### Changes
<details><summary>5 commits</summary>
<p>

* [`56744cf`](https://github.com/talos-systems/sidero/commit/56744cfb03f2dd562c4ee88bc56168f3a5b135b1) feat: provide a way to disable DHCP proxy
* [`7880ee2`](https://github.com/talos-systems/sidero/commit/7880ee21b7fa2690a663f260598de9dd247d0060) feat: update kube-rbac-proxy to support arm64
* [`42685a8`](https://github.com/talos-systems/sidero/commit/42685a8c1c4500bee736edd97958cb8b656533e9) docs: fix port count for DHCP proxy
* [`891edce`](https://github.com/talos-systems/sidero/commit/891edce4f16975349e3a20d4fdc68f5118dcb77a) docs: remove excess bash continuation backslashes
* [`8820a2b`](https://github.com/talos-systems/sidero/commit/8820a2bd5a5d9807920a11242d1a9a9a005d9d33) docs: update documentation for Sidero v0.6
</p>
</details>

### Dependency Changes

* **github.com/insomniacslk/dhcp**  b3ca2534940d -> 65c27093e38a
* **go.uber.org/zap**               v1.25.0 -> v1.26.0
* **golang.org/x/net**              v0.14.0 -> v0.17.0
* **golang.org/x/sync**             v0.3.0 -> v0.4.0
* **golang.org/x/sys**              v0.11.0 -> v0.13.0
* **google.golang.org/grpc**        v1.57.0 -> v1.58.3
* **sigs.k8s.io/cluster-api**       v1.5.0 -> v1.5.2

Previous release can be found at [v0.6.0](https://github.com/talos-systems/sidero/releases/tag/v0.6.0)

## [Sidero 0.6.0](https://github.com/talos-systems/sidero/releases/tag/v0.6.0) (2023-08-29)

Welcome to the v0.6.0 release of Sidero!



Please try out the release binaries and report any issues at
https://github.com/talos-systems/sidero/issues.

### Sidero Agent

Sidero Agent now runs DHCP client in the userland, on the link which was used to PXE boot the machine.
This allows to run Sidero Agent on the machine with several autoconfigured network interfaces, when one of them is used for the management network.


### DHCP Proxy

Sidero Controller Manager now includes DHCP proxy which augments DHCP response with additional PXE boot options.
When enabled, DHCP server in the environment only handles IP allocation and network configuration, while DHCP proxy
provides PXE boot information automatically based on the architecture and boot method.


### Metadata Server

Sidero Metadata Server no longer depends on the version of Talos machinery library it is built with.
Sidero should be able to process machine config for future versions of Talos.


### New API Version for `metal.sidero.dev` Resources

Resources under `metal.sidero.dev` (`Server`, `ServerClass`, `Environment`) now have a new version `v1alpha2`.
Old version `v1alpha1` is still supported, but it is recommended to update templates to use the new resource version.

#### `Server` Changes

Hardware information was restructured and extended when compared with `v1alpha1`:

* `.spec.systemInformation` -> `.spec.hardware.system`
* `.spec.cpu` -> `.spec.hardware.compute.processors[]`

#### `ServerClass` Changes

* `.spec.qualifiers.systemInformation` -> `.spec.qualifiers.system`
* `.spec.qualifiers.cpu` -> `.spec.qualifiers.hardware.compute.processors[]`


### Contributors

* Andrey Smirnov
* Spencer Smith
* Gerard de Leeuw
* Tim Jones
* Andrey Smirnov
* Artem Chernyshev
* Noel Georgi
* JJGadgets
* Martin Sweeny
* Michael Vorburger ⛑️
* Russell Troxel
* Sander Flobbe
* Steve Francis
* Utku Ozdemir
* Zach Bernstein
* bzub

### Changes
<details><summary>65 commits</summary>
<p>

* [`adea408`](https://github.com/talos-systems/sidero/commit/adea408cf803b9d599d86bc10d01671ca28b1383) release(v0.6.0-beta.0): prepare release
* [`f88ee6d`](https://github.com/talos-systems/sidero/commit/f88ee6d2e6e8b2908ccb069e0d7a6bdbde312aa5) feat: add proxy DHCP server
* [`baaece8`](https://github.com/talos-systems/sidero/commit/baaece89a6cda5e77393318e42fae051866e8d46) feat: update for Talos 1.5.0
* [`f711a64`](https://github.com/talos-systems/sidero/commit/f711a64b60f20b3b44e0e218d85a7767a730fcb2) docs: code snippet lacks "sudo"
* [`45ac3d8`](https://github.com/talos-systems/sidero/commit/45ac3d899f4082de13f1ce37d8bacd1b0f7f1dda) fix: build siderolink endpoint correctly
* [`049a5fa`](https://github.com/talos-systems/sidero/commit/049a5faad74941c5f37f47522008cfe364a0adcd) release(v0.6.0-alpha.0): prepare release
* [`b2c53bf`](https://github.com/talos-systems/sidero/commit/b2c53bf9d7d016d418365ee3caff352ddae86f68) fix: allow 'empty user' slots to be used in IPMI setup
* [`2b3dedc`](https://github.com/talos-systems/sidero/commit/2b3dedc23b1bf1c109fe7426471e1209acf91455) feat: update to final Talos 1.4.0 release
* [`6c15a40`](https://github.com/talos-systems/sidero/commit/6c15a401ebe59b43cb3be46a1e5aed97c3462cf6) fix: parse Talos events into Machine conditions
* [`24e1893`](https://github.com/talos-systems/sidero/commit/24e189331df34f55c4583a0a45bd6dc3b096b116) fix: break the link of metadata server to Talos machine config
* [`ef65ff0`](https://github.com/talos-systems/sidero/commit/ef65ff05a91bc49d0bc8474cefe3f5b0c880c4cb) chore: add v0.6.x to metadata, fix metrics service
* [`e433504`](https://github.com/talos-systems/sidero/commit/e433504087ae39299d1c2e492db6ed5714dd73a0) feat: provide 'snp.efi' and recommend it over 'ipxe.efi'
* [`71390a3`](https://github.com/talos-systems/sidero/commit/71390a33dba10a5e2f02ade07805cd4435d69d5a) docs: remove docs v0.1-v0.3
* [`b305f2c`](https://github.com/talos-systems/sidero/commit/b305f2c5f0970b0f779bd465537546476e99d67c) feat: rework the way Sidero Agent boots and configures networking
* [`b6235eb`](https://github.com/talos-systems/sidero/commit/b6235eb4bdc562eb33c32fbf217bb81e2c68b579) test: drop tests with old Talos compatibility
* [`9121a3b`](https://github.com/talos-systems/sidero/commit/9121a3b7d53b3897fe60c734b6bfa3599954710f) chore: bump dependencies
* [`18116bc`](https://github.com/talos-systems/sidero/commit/18116bcabe77f1e7aa8d33b4d10cf2ff4a7cd18b) fix: use updated pkgs with fixed ipmitool
* [`fbcd02a`](https://github.com/talos-systems/sidero/commit/fbcd02a45d49a774ee9a6327347d342edd4ed2fa) feat: update for Talos 1.3.0
* [`c95184a`](https://github.com/talos-systems/sidero/commit/c95184a318544a45203825c94e187fc50c0b1919) docs: detail how to disable IPMI magic (re. #988)
* [`51ac816`](https://github.com/talos-systems/sidero/commit/51ac816e6f1bd2625030be3f0044cabe7849e55d) fix: docs: /etc/dhcp -> config to preserve across firmware update
* [`d8ef68b`](https://github.com/talos-systems/sidero/commit/d8ef68b6b951df093d288c1dae6c0d32a9e00aec) feat: provide a way to override IPMI PXE boot method on `Server`
* [`831761a`](https://github.com/talos-systems/sidero/commit/831761aeffaa9534e3072d1b894df97a0b495809) fix: docs sitemap www prefix
* [`2b2ed86`](https://github.com/talos-systems/sidero/commit/2b2ed86e1ceb9b6f1a74be85a6546a6a61f73c8c) docs: fix dead link in serverclasses guide
* [`5527432`](https://github.com/talos-systems/sidero/commit/5527432ddfc0fca78ad8bd9a1034a431117fcc06) fix: canonical links in docs site
* [`757706c`](https://github.com/talos-systems/sidero/commit/757706cae125abfdcb8ff49167539abd832b7c6a) feat: finalize update to Talos 1.2.0
* [`f9d12f3`](https://github.com/talos-systems/sidero/commit/f9d12f3c4635281015c6e3795aadfbc61dad805c) feat: update Talos to 1.2.0-beta.2
* [`44f1962`](https://github.com/talos-systems/sidero/commit/44f1962ceaac518f1e59205a43db3a1bf1e79cd6) fix: properly inspect readonly flag of the disk
* [`7cb28fc`](https://github.com/talos-systems/sidero/commit/7cb28fccfea5be7f5a22ddd9f895cd04cda68dc1) docs: handle latest version banner and canonical urls
* [`824d059`](https://github.com/talos-systems/sidero/commit/824d05920a8add702a9e705321585f247ea4ca48) fix: filter out SideroLink address from Machine addresses
* [`08cac99`](https://github.com/talos-systems/sidero/commit/08cac998b1e5228cf4daa0d1911ef7d0a8b4ae16) feat: update base Talos to 1.1.1
* [`3c65480`](https://github.com/talos-systems/sidero/commit/3c654806a52bbeb8604b78f27f2fafc0a75f2eb9) fix: skip read-only block devices on wiping
* [`4e714a1`](https://github.com/talos-systems/sidero/commit/4e714a1e4d657e01fcfe0f7a3e64ca58d6af89f5) fix: resolve wireguard endpoint to IP
* [`3eddbbd`](https://github.com/talos-systems/sidero/commit/3eddbbd69ecab561bc6583d75262442ed4b73228) fix: use http response writer directly
* [`a06135a`](https://github.com/talos-systems/sidero/commit/a06135a96e0d19c07d5e34908ebffee1ea15571d) docs: use ipxe-sanboot for raspberry pi 4's
* [`c0dfd23`](https://github.com/talos-systems/sidero/commit/c0dfd23fa5f8b775bf1ea84ce14c93927b821a88) test: integrate new controlplane loadbalancer
* [`fd0086d`](https://github.com/talos-systems/sidero/commit/fd0086d6f0579f3888cb0dfd8476b6c91cdd5c28) fix: append Talos default kernel args even if there is something defined
* [`6c6b354`](https://github.com/talos-systems/sidero/commit/6c6b3546e40665b605fa44cc5508647c90d20223) docs: edit for clarity and conciseness
* [`56d27ce`](https://github.com/talos-systems/sidero/commit/56d27ce71876df1626edcc56cb19c5c7b53d5aca) chore: bump dependencies
* [`b8843e2`](https://github.com/talos-systems/sidero/commit/b8843e2f875c0dc73a3acdba5ed65a3d9a869e04) chore: bump dependencies
* [`03a5c90`](https://github.com/talos-systems/sidero/commit/03a5c90b2c961d92709e3ff318e340ec6a5d0105) chore: bump dependencies
* [`6092854`](https://github.com/talos-systems/sidero/commit/609285488e13423626dc9131fef5c77e91ab9f5a) docs: fix spelling mistake
* [`c29d464`](https://github.com/talos-systems/sidero/commit/c29d4645d9dbb3715e6829948ba42fc0f518817b) feat: make MetalMachineTemplate immutable
* [`835d5cf`](https://github.com/talos-systems/sidero/commit/835d5cf227961d9a423985277d80767407dc7eac) docs: fix links and getting started overview
* [`7c7a377`](https://github.com/talos-systems/sidero/commit/7c7a37785711a5a90afd701921995ba88755ccce) feat: add extended hardware information to Server and ServerClass CRDs
* [`e4bb416`](https://github.com/talos-systems/sidero/commit/e4bb4165e9d25cceabd2e4cd73e5a53b4d5756ef) feat: restructure HW information in `Server` resources (v1alpha2)
* [`6c81518`](https://github.com/talos-systems/sidero/commit/6c81518b172ba09efbc46f15e33706d7340cdecf) docs: fork docs for Sidero 0.6
* [`061ee8e`](https://github.com/talos-systems/sidero/commit/061ee8e57debf11457e732f74be6e654e31e0203) refactor: remove pre-ServerBinding leftovers, use remote cached clients
* [`24449aa`](https://github.com/talos-systems/sidero/commit/24449aa3993f75299a8cf2113aa82a8aba65e551) chore: bump dependencies
* [`fd34779`](https://github.com/talos-systems/sidero/commit/fd34779bf7ada1914e320997c73fecc624813e8a) docs: add info about GitHub Org rename to v0.5 docs
* [`c86953e`](https://github.com/talos-systems/sidero/commit/c86953e11c0f9a82e7279b1f3977d869a33870bb) docs: update algolia info
* [`511ddfc`](https://github.com/talos-systems/sidero/commit/511ddfc0828a3e9b5ac26814d34ae9c5b117071a) feat: enable webhooks in `sidero-controller-manager`
* [`37a1d52`](https://github.com/talos-systems/sidero/commit/37a1d526962c0d3f6cbef28a43c6ca86d555df81) docs: add sitemap override
* [`54f896d`](https://github.com/talos-systems/sidero/commit/54f896dd6750874a2d065dc00ec96a56494a3f9d) feat: allow configuring Sidero deployment strategy
* [`3be5e6e`](https://github.com/talos-systems/sidero/commit/3be5e6e817fc8db6a8e4e2566cbad311d54ee7ad) chore: bump dependencies
* [`b7cc8b2`](https://github.com/talos-systems/sidero/commit/b7cc8b2e40d687b51558bb73efc014fc53e60183) docs: switch code block to text for DHCP prereq
* [`15f6730`](https://github.com/talos-systems/sidero/commit/15f67308f0bdbd52a4565deff8923d67a4c4c755) chore: use consistent naming for imports
* [`8985a04`](https://github.com/talos-systems/sidero/commit/8985a04231ed91f79477d5556a3256919287efb4) fix: ipxe prompt on ARM64
* [`9f9c922`](https://github.com/talos-systems/sidero/commit/9f9c92205c40e289975e9470625f4da0bf6cad64) docs: move docs to hugo
* [`2a475db`](https://github.com/talos-systems/sidero/commit/2a475db79ee6df559cecd082581fb7ba67d4f35b) feat: update to Talos 1.0
* [`7a1a930`](https://github.com/talos-systems/sidero/commit/7a1a930989d1b93667645c8b1d075a99dd31d799) chore: update more registries
* [`9438ca8`](https://github.com/talos-systems/sidero/commit/9438ca8e8aaafff3270bac85c5982918f3333278) chore: fix gpg check and sfyra build
* [`5fb30a9`](https://github.com/talos-systems/sidero/commit/5fb30a996c3ed4a4239ddfd48c67dd8f73f4d0b5) docs: update README for 0.5 version compatibility
* [`4d11603`](https://github.com/talos-systems/sidero/commit/4d116038e09d2c3aa62486ae4eac4f8ea863f1df) chore: bump cert-manager to v1
* [`13e11d3`](https://github.com/talos-systems/sidero/commit/13e11d3320c529c0934f0b2e5e5b80c0e36a14cc) chore: bump dependencies
* [`d3e75dc`](https://github.com/talos-systems/sidero/commit/d3e75dc9beff977b3c0e754805ccf94fb2e77001) docs: make 0.5 docs the default, add what's new
</p>
</details>

### Changes since v0.6.0-beta.0
<details><summary>0 commit</summary>
<p>

</p>
</details>

### Dependency Changes

* **github.com/go-logr/logr**                       v0.4.0 -> v1.2.4
* **github.com/google/go-cmp**                      v0.5.9 **_new_**
* **github.com/grpc-ecosystem/go-grpc-middleware**  v1.3.0 -> v1.4.0
* **github.com/insomniacslk/dhcp**                  b3ca2534940d **_new_**
* **github.com/jsimonetti/rtnetlink**               v1.3.5 **_new_**
* **github.com/siderolabs/gen**                     v0.4.5 **_new_**
* **github.com/siderolabs/go-blockdevice**          v0.4.6 **_new_**
* **github.com/siderolabs/go-cmd**                  v0.1.1 **_new_**
* **github.com/siderolabs/go-debug**                v0.2.3 **_new_**
* **github.com/siderolabs/go-kmsg**                 v0.1.3 **_new_**
* **github.com/siderolabs/go-procfs**               v0.1.1 **_new_**
* **github.com/siderolabs/go-retry**                v0.3.2 **_new_**
* **github.com/siderolabs/go-smbios**               v0.3.2 **_new_**
* **github.com/siderolabs/grpc-proxy**              v0.4.0 **_new_**
* **github.com/siderolabs/siderolink**              v0.3.1 **_new_**
* **github.com/siderolabs/talos/pkg/machinery**     v1.5.0 **_new_**
* **github.com/stretchr/testify**                   v1.7.0 -> v1.8.4
* **go.uber.org/zap**                               v1.20.0 -> v1.25.0
* **golang.org/x/net**                              cd36cc0744dd -> v0.14.0
* **golang.org/x/sync**                             036812b2e83c -> v0.3.0
* **golang.org/x/sys**                              99c3d69c2c27 -> v0.11.0
* **golang.zx2c4.com/wireguard/wgctrl**             daad0b7ba671 -> 925a1e7659e6
* **google.golang.org/grpc**                        v1.44.0 -> v1.57.0
* **google.golang.org/protobuf**                    v1.27.1 -> v1.31.0
* **gopkg.in/yaml.v3**                              v3.0.1 **_new_**
* **k8s.io/api**                                    v0.22.2 -> v0.27.2
* **k8s.io/apiextensions-apiserver**                v0.22.2 -> v0.27.2
* **k8s.io/apimachinery**                           v0.22.2 -> v0.27.2
* **k8s.io/client-go**                              v0.22.2 -> v0.27.2
* **k8s.io/utils**                                  cb0fa318a74b -> a36077c30491
* **sigs.k8s.io/cluster-api**                       v1.0.4 -> v1.5.0
* **sigs.k8s.io/controller-runtime**                v0.10.3 -> v0.15.1

Previous release can be found at [v0.5.0](https://github.com/talos-systems/sidero/releases/tag/v0.5.0)

## [Sidero 0.6.0-beta.0](https://github.com/talos-systems/sidero/releases/tag/v0.6.0-beta.0) (2023-08-24)

Welcome to the v0.6.0-beta.0 release of Sidero!  
*This is a pre-release of Sidero*



Please try out the release binaries and report any issues at
https://github.com/talos-systems/sidero/issues.

### Sidero Agent

Sidero Agent now runs DHCP client in the userland, on the link which was used to PXE boot the machine.
This allows to run Sidero Agent on the machine with several autoconfigured network interfaces, when one of them is used for the management network.


### DHCP Proxy

Sidero Controller Manager now includes DHCP proxy which augments DHCP response with additional PXE boot options.
When enabled, DHCP server in the environment only handles IP allocation and network configuration, while DHCP proxy
provides PXE boot information automatically based on the architecture and boot method.


### Metadata Server

Sidero Metadata Server no longer depends on the version of Talos machinery library it is built with.
Sidero should be able to process machine config for future versions of Talos.


### New API Version for `metal.sidero.dev` Resources

Resources under `metal.sidero.dev` (`Server`, `ServerClass`, `Environment`) now have a new version `v1alpha2`.
Old version `v1alpha1` is still supported, but it is recommended to update templates to use the new resource version.

#### `Server` Changes

Hardware information was restructured and extended when compared with `v1alpha1`:

* `.spec.systemInformation` -> `.spec.hardware.system`
* `.spec.cpu` -> `.spec.hardware.compute.processors[]`

#### `ServerClass` Changes

* `.spec.qualifiers.systemInformation` -> `.spec.qualifiers.system`
* `.spec.qualifiers.cpu` -> `.spec.qualifiers.hardware.compute.processors[]`


### Contributors

* Andrey Smirnov
* Spencer Smith
* Gerard de Leeuw
* Tim Jones
* Artem Chernyshev
* Noel Georgi
* Andrey Smirnov
* JJGadgets
* Martin Sweeny
* Michael Vorburger ⛑️
* Russell Troxel
* Sander Flobbe
* Steve Francis
* Utku Ozdemir
* Zach Bernstein
* bzub

### Changes
<details><summary>64 commits</summary>
<p>

* [`f88ee6d`](https://github.com/talos-systems/sidero/commit/f88ee6d2e6e8b2908ccb069e0d7a6bdbde312aa5) feat: add proxy DHCP server
* [`baaece8`](https://github.com/talos-systems/sidero/commit/baaece89a6cda5e77393318e42fae051866e8d46) feat: update for Talos 1.5.0
* [`f711a64`](https://github.com/talos-systems/sidero/commit/f711a64b60f20b3b44e0e218d85a7767a730fcb2) docs: code snippet lacks "sudo"
* [`45ac3d8`](https://github.com/talos-systems/sidero/commit/45ac3d899f4082de13f1ce37d8bacd1b0f7f1dda) fix: build siderolink endpoint correctly
* [`049a5fa`](https://github.com/talos-systems/sidero/commit/049a5faad74941c5f37f47522008cfe364a0adcd) release(v0.6.0-alpha.0): prepare release
* [`b2c53bf`](https://github.com/talos-systems/sidero/commit/b2c53bf9d7d016d418365ee3caff352ddae86f68) fix: allow 'empty user' slots to be used in IPMI setup
* [`2b3dedc`](https://github.com/talos-systems/sidero/commit/2b3dedc23b1bf1c109fe7426471e1209acf91455) feat: update to final Talos 1.4.0 release
* [`6c15a40`](https://github.com/talos-systems/sidero/commit/6c15a401ebe59b43cb3be46a1e5aed97c3462cf6) fix: parse Talos events into Machine conditions
* [`24e1893`](https://github.com/talos-systems/sidero/commit/24e189331df34f55c4583a0a45bd6dc3b096b116) fix: break the link of metadata server to Talos machine config
* [`ef65ff0`](https://github.com/talos-systems/sidero/commit/ef65ff05a91bc49d0bc8474cefe3f5b0c880c4cb) chore: add v0.6.x to metadata, fix metrics service
* [`e433504`](https://github.com/talos-systems/sidero/commit/e433504087ae39299d1c2e492db6ed5714dd73a0) feat: provide 'snp.efi' and recommend it over 'ipxe.efi'
* [`71390a3`](https://github.com/talos-systems/sidero/commit/71390a33dba10a5e2f02ade07805cd4435d69d5a) docs: remove docs v0.1-v0.3
* [`b305f2c`](https://github.com/talos-systems/sidero/commit/b305f2c5f0970b0f779bd465537546476e99d67c) feat: rework the way Sidero Agent boots and configures networking
* [`b6235eb`](https://github.com/talos-systems/sidero/commit/b6235eb4bdc562eb33c32fbf217bb81e2c68b579) test: drop tests with old Talos compatibility
* [`9121a3b`](https://github.com/talos-systems/sidero/commit/9121a3b7d53b3897fe60c734b6bfa3599954710f) chore: bump dependencies
* [`18116bc`](https://github.com/talos-systems/sidero/commit/18116bcabe77f1e7aa8d33b4d10cf2ff4a7cd18b) fix: use updated pkgs with fixed ipmitool
* [`fbcd02a`](https://github.com/talos-systems/sidero/commit/fbcd02a45d49a774ee9a6327347d342edd4ed2fa) feat: update for Talos 1.3.0
* [`c95184a`](https://github.com/talos-systems/sidero/commit/c95184a318544a45203825c94e187fc50c0b1919) docs: detail how to disable IPMI magic (re. #988)
* [`51ac816`](https://github.com/talos-systems/sidero/commit/51ac816e6f1bd2625030be3f0044cabe7849e55d) fix: docs: /etc/dhcp -> config to preserve across firmware update
* [`d8ef68b`](https://github.com/talos-systems/sidero/commit/d8ef68b6b951df093d288c1dae6c0d32a9e00aec) feat: provide a way to override IPMI PXE boot method on `Server`
* [`831761a`](https://github.com/talos-systems/sidero/commit/831761aeffaa9534e3072d1b894df97a0b495809) fix: docs sitemap www prefix
* [`2b2ed86`](https://github.com/talos-systems/sidero/commit/2b2ed86e1ceb9b6f1a74be85a6546a6a61f73c8c) docs: fix dead link in serverclasses guide
* [`5527432`](https://github.com/talos-systems/sidero/commit/5527432ddfc0fca78ad8bd9a1034a431117fcc06) fix: canonical links in docs site
* [`757706c`](https://github.com/talos-systems/sidero/commit/757706cae125abfdcb8ff49167539abd832b7c6a) feat: finalize update to Talos 1.2.0
* [`f9d12f3`](https://github.com/talos-systems/sidero/commit/f9d12f3c4635281015c6e3795aadfbc61dad805c) feat: update Talos to 1.2.0-beta.2
* [`44f1962`](https://github.com/talos-systems/sidero/commit/44f1962ceaac518f1e59205a43db3a1bf1e79cd6) fix: properly inspect readonly flag of the disk
* [`7cb28fc`](https://github.com/talos-systems/sidero/commit/7cb28fccfea5be7f5a22ddd9f895cd04cda68dc1) docs: handle latest version banner and canonical urls
* [`824d059`](https://github.com/talos-systems/sidero/commit/824d05920a8add702a9e705321585f247ea4ca48) fix: filter out SideroLink address from Machine addresses
* [`08cac99`](https://github.com/talos-systems/sidero/commit/08cac998b1e5228cf4daa0d1911ef7d0a8b4ae16) feat: update base Talos to 1.1.1
* [`3c65480`](https://github.com/talos-systems/sidero/commit/3c654806a52bbeb8604b78f27f2fafc0a75f2eb9) fix: skip read-only block devices on wiping
* [`4e714a1`](https://github.com/talos-systems/sidero/commit/4e714a1e4d657e01fcfe0f7a3e64ca58d6af89f5) fix: resolve wireguard endpoint to IP
* [`3eddbbd`](https://github.com/talos-systems/sidero/commit/3eddbbd69ecab561bc6583d75262442ed4b73228) fix: use http response writer directly
* [`a06135a`](https://github.com/talos-systems/sidero/commit/a06135a96e0d19c07d5e34908ebffee1ea15571d) docs: use ipxe-sanboot for raspberry pi 4's
* [`c0dfd23`](https://github.com/talos-systems/sidero/commit/c0dfd23fa5f8b775bf1ea84ce14c93927b821a88) test: integrate new controlplane loadbalancer
* [`fd0086d`](https://github.com/talos-systems/sidero/commit/fd0086d6f0579f3888cb0dfd8476b6c91cdd5c28) fix: append Talos default kernel args even if there is something defined
* [`6c6b354`](https://github.com/talos-systems/sidero/commit/6c6b3546e40665b605fa44cc5508647c90d20223) docs: edit for clarity and conciseness
* [`56d27ce`](https://github.com/talos-systems/sidero/commit/56d27ce71876df1626edcc56cb19c5c7b53d5aca) chore: bump dependencies
* [`b8843e2`](https://github.com/talos-systems/sidero/commit/b8843e2f875c0dc73a3acdba5ed65a3d9a869e04) chore: bump dependencies
* [`03a5c90`](https://github.com/talos-systems/sidero/commit/03a5c90b2c961d92709e3ff318e340ec6a5d0105) chore: bump dependencies
* [`6092854`](https://github.com/talos-systems/sidero/commit/609285488e13423626dc9131fef5c77e91ab9f5a) docs: fix spelling mistake
* [`c29d464`](https://github.com/talos-systems/sidero/commit/c29d4645d9dbb3715e6829948ba42fc0f518817b) feat: make MetalMachineTemplate immutable
* [`835d5cf`](https://github.com/talos-systems/sidero/commit/835d5cf227961d9a423985277d80767407dc7eac) docs: fix links and getting started overview
* [`7c7a377`](https://github.com/talos-systems/sidero/commit/7c7a37785711a5a90afd701921995ba88755ccce) feat: add extended hardware information to Server and ServerClass CRDs
* [`e4bb416`](https://github.com/talos-systems/sidero/commit/e4bb4165e9d25cceabd2e4cd73e5a53b4d5756ef) feat: restructure HW information in `Server` resources (v1alpha2)
* [`6c81518`](https://github.com/talos-systems/sidero/commit/6c81518b172ba09efbc46f15e33706d7340cdecf) docs: fork docs for Sidero 0.6
* [`061ee8e`](https://github.com/talos-systems/sidero/commit/061ee8e57debf11457e732f74be6e654e31e0203) refactor: remove pre-ServerBinding leftovers, use remote cached clients
* [`24449aa`](https://github.com/talos-systems/sidero/commit/24449aa3993f75299a8cf2113aa82a8aba65e551) chore: bump dependencies
* [`fd34779`](https://github.com/talos-systems/sidero/commit/fd34779bf7ada1914e320997c73fecc624813e8a) docs: add info about GitHub Org rename to v0.5 docs
* [`c86953e`](https://github.com/talos-systems/sidero/commit/c86953e11c0f9a82e7279b1f3977d869a33870bb) docs: update algolia info
* [`511ddfc`](https://github.com/talos-systems/sidero/commit/511ddfc0828a3e9b5ac26814d34ae9c5b117071a) feat: enable webhooks in `sidero-controller-manager`
* [`37a1d52`](https://github.com/talos-systems/sidero/commit/37a1d526962c0d3f6cbef28a43c6ca86d555df81) docs: add sitemap override
* [`54f896d`](https://github.com/talos-systems/sidero/commit/54f896dd6750874a2d065dc00ec96a56494a3f9d) feat: allow configuring Sidero deployment strategy
* [`3be5e6e`](https://github.com/talos-systems/sidero/commit/3be5e6e817fc8db6a8e4e2566cbad311d54ee7ad) chore: bump dependencies
* [`b7cc8b2`](https://github.com/talos-systems/sidero/commit/b7cc8b2e40d687b51558bb73efc014fc53e60183) docs: switch code block to text for DHCP prereq
* [`15f6730`](https://github.com/talos-systems/sidero/commit/15f67308f0bdbd52a4565deff8923d67a4c4c755) chore: use consistent naming for imports
* [`8985a04`](https://github.com/talos-systems/sidero/commit/8985a04231ed91f79477d5556a3256919287efb4) fix: ipxe prompt on ARM64
* [`9f9c922`](https://github.com/talos-systems/sidero/commit/9f9c92205c40e289975e9470625f4da0bf6cad64) docs: move docs to hugo
* [`2a475db`](https://github.com/talos-systems/sidero/commit/2a475db79ee6df559cecd082581fb7ba67d4f35b) feat: update to Talos 1.0
* [`7a1a930`](https://github.com/talos-systems/sidero/commit/7a1a930989d1b93667645c8b1d075a99dd31d799) chore: update more registries
* [`9438ca8`](https://github.com/talos-systems/sidero/commit/9438ca8e8aaafff3270bac85c5982918f3333278) chore: fix gpg check and sfyra build
* [`5fb30a9`](https://github.com/talos-systems/sidero/commit/5fb30a996c3ed4a4239ddfd48c67dd8f73f4d0b5) docs: update README for 0.5 version compatibility
* [`4d11603`](https://github.com/talos-systems/sidero/commit/4d116038e09d2c3aa62486ae4eac4f8ea863f1df) chore: bump cert-manager to v1
* [`13e11d3`](https://github.com/talos-systems/sidero/commit/13e11d3320c529c0934f0b2e5e5b80c0e36a14cc) chore: bump dependencies
* [`d3e75dc`](https://github.com/talos-systems/sidero/commit/d3e75dc9beff977b3c0e754805ccf94fb2e77001) docs: make 0.5 docs the default, add what's new
</p>
</details>

### Changes since v0.6.0-alpha.0-42-g5b4ff34
<details><summary>2 commits</summary>
<p>

* [`f88ee6d`](https://github.com/talos-systems/sidero/commit/f88ee6d2e6e8b2908ccb069e0d7a6bdbde312aa5) feat: add proxy DHCP server
* [`baaece8`](https://github.com/talos-systems/sidero/commit/baaece89a6cda5e77393318e42fae051866e8d46) feat: update for Talos 1.5.0
</p>
</details>

### Dependency Changes

* **github.com/go-logr/logr**                       v0.4.0 -> v1.2.4
* **github.com/google/go-cmp**                      v0.5.9 **_new_**
* **github.com/grpc-ecosystem/go-grpc-middleware**  v1.3.0 -> v1.4.0
* **github.com/insomniacslk/dhcp**                  b3ca2534940d **_new_**
* **github.com/jsimonetti/rtnetlink**               v1.3.5 **_new_**
* **github.com/siderolabs/gen**                     v0.4.5 **_new_**
* **github.com/siderolabs/go-blockdevice**          v0.4.6 **_new_**
* **github.com/siderolabs/go-cmd**                  v0.1.1 **_new_**
* **github.com/siderolabs/go-debug**                v0.2.3 **_new_**
* **github.com/siderolabs/go-kmsg**                 v0.1.3 **_new_**
* **github.com/siderolabs/go-procfs**               v0.1.1 **_new_**
* **github.com/siderolabs/go-retry**                v0.3.2 **_new_**
* **github.com/siderolabs/go-smbios**               v0.3.2 **_new_**
* **github.com/siderolabs/grpc-proxy**              v0.4.0 **_new_**
* **github.com/siderolabs/siderolink**              v0.3.1 **_new_**
* **github.com/siderolabs/talos/pkg/machinery**     v1.5.0 **_new_**
* **github.com/stretchr/testify**                   v1.7.0 -> v1.8.4
* **go.uber.org/zap**                               v1.20.0 -> v1.25.0
* **golang.org/x/net**                              cd36cc0744dd -> v0.14.0
* **golang.org/x/sync**                             036812b2e83c -> v0.3.0
* **golang.org/x/sys**                              99c3d69c2c27 -> v0.11.0
* **golang.zx2c4.com/wireguard/wgctrl**             daad0b7ba671 -> 925a1e7659e6
* **google.golang.org/grpc**                        v1.44.0 -> v1.57.0
* **google.golang.org/protobuf**                    v1.27.1 -> v1.31.0
* **gopkg.in/yaml.v3**                              v3.0.1 **_new_**
* **k8s.io/api**                                    v0.22.2 -> v0.27.2
* **k8s.io/apiextensions-apiserver**                v0.22.2 -> v0.27.2
* **k8s.io/apimachinery**                           v0.22.2 -> v0.27.2
* **k8s.io/client-go**                              v0.22.2 -> v0.27.2
* **k8s.io/utils**                                  cb0fa318a74b -> a36077c30491
* **sigs.k8s.io/cluster-api**                       v1.0.4 -> v1.5.0
* **sigs.k8s.io/controller-runtime**                v0.10.3 -> v0.15.1

Previous release can be found at [v0.5.0](https://github.com/talos-systems/sidero/releases/tag/v0.5.0)

## [Sidero 0.6.0-alpha.0](https://github.com/talos-systems/sidero/releases/tag/v0.6.0-alpha.0) (2023-04-20)

Welcome to the v0.6.0-alpha.0 release of Sidero!  
*This is a pre-release of Sidero*



Please try out the release binaries and report any issues at
https://github.com/talos-systems/sidero/issues.

### Sidero Agent

Sidero Agent now runs DHCP client in the userland, on the link which was used to PXE boot the machine.
This allows to run Sidero Agent on the machine with several autoconfigured network interfaces, when one of them is used for the management network.


### Metadata Server

Sidero Metadata Server no longer depends on the version of Talos machinery library it is built with.
Sidero should be able to process machine config for future versions of Talos.


### New API Version for `metal.sidero.dev` Resources

Resources under `metal.sidero.dev` (`Server`, `ServerClass`, `Environment`) now have a new version `v1alpha2`.
Old version `v1alpha1` is still supported, but it is recommended to update templates to use the new resource version.

#### `Server` Changes

Hardware information was restructured and extended when compared with `v1alpha1`:

* `.spec.systemInformation` -> `.spec.hardware.system`
* `.spec.cpu` -> `.spec.hardware.compute.processors[]`

#### `ServerClass` Changes

* `.spec.qualifiers.systemInformation` -> `.spec.qualifiers.system`
* `.spec.qualifiers.cpu` -> `.spec.qualifiers.hardware.compute.processors[]`


### Contributors

* Andrey Smirnov
* Spencer Smith
* Gerard de Leeuw
* Tim Jones
* Artem Chernyshev
* Noel Georgi
* JJGadgets
* Martin Sweeny
* Michael Vorburger ⛑️
* Russell Troxel
* Steve Francis
* Utku Ozdemir
* Zach Bernstein
* bzub

### Changes
<details><summary>59 commits</summary>
<p>

* [`b2c53bf`](https://github.com/talos-systems/sidero/commit/b2c53bf9d7d016d418365ee3caff352ddae86f68) fix: allow 'empty user' slots to be used in IPMI setup
* [`2b3dedc`](https://github.com/talos-systems/sidero/commit/2b3dedc23b1bf1c109fe7426471e1209acf91455) feat: update to final Talos 1.4.0 release
* [`6c15a40`](https://github.com/talos-systems/sidero/commit/6c15a401ebe59b43cb3be46a1e5aed97c3462cf6) fix: parse Talos events into Machine conditions
* [`24e1893`](https://github.com/talos-systems/sidero/commit/24e189331df34f55c4583a0a45bd6dc3b096b116) fix: break the link of metadata server to Talos machine config
* [`ef65ff0`](https://github.com/talos-systems/sidero/commit/ef65ff05a91bc49d0bc8474cefe3f5b0c880c4cb) chore: add v0.6.x to metadata, fix metrics service
* [`e433504`](https://github.com/talos-systems/sidero/commit/e433504087ae39299d1c2e492db6ed5714dd73a0) feat: provide 'snp.efi' and recommend it over 'ipxe.efi'
* [`71390a3`](https://github.com/talos-systems/sidero/commit/71390a33dba10a5e2f02ade07805cd4435d69d5a) docs: remove docs v0.1-v0.3
* [`b305f2c`](https://github.com/talos-systems/sidero/commit/b305f2c5f0970b0f779bd465537546476e99d67c) feat: rework the way Sidero Agent boots and configures networking
* [`b6235eb`](https://github.com/talos-systems/sidero/commit/b6235eb4bdc562eb33c32fbf217bb81e2c68b579) test: drop tests with old Talos compatibility
* [`9121a3b`](https://github.com/talos-systems/sidero/commit/9121a3b7d53b3897fe60c734b6bfa3599954710f) chore: bump dependencies
* [`18116bc`](https://github.com/talos-systems/sidero/commit/18116bcabe77f1e7aa8d33b4d10cf2ff4a7cd18b) fix: use updated pkgs with fixed ipmitool
* [`fbcd02a`](https://github.com/talos-systems/sidero/commit/fbcd02a45d49a774ee9a6327347d342edd4ed2fa) feat: update for Talos 1.3.0
* [`c95184a`](https://github.com/talos-systems/sidero/commit/c95184a318544a45203825c94e187fc50c0b1919) docs: detail how to disable IPMI magic (re. #988)
* [`51ac816`](https://github.com/talos-systems/sidero/commit/51ac816e6f1bd2625030be3f0044cabe7849e55d) fix: docs: /etc/dhcp -> config to preserve across firmware update
* [`d8ef68b`](https://github.com/talos-systems/sidero/commit/d8ef68b6b951df093d288c1dae6c0d32a9e00aec) feat: provide a way to override IPMI PXE boot method on `Server`
* [`831761a`](https://github.com/talos-systems/sidero/commit/831761aeffaa9534e3072d1b894df97a0b495809) fix: docs sitemap www prefix
* [`2b2ed86`](https://github.com/talos-systems/sidero/commit/2b2ed86e1ceb9b6f1a74be85a6546a6a61f73c8c) docs: fix dead link in serverclasses guide
* [`5527432`](https://github.com/talos-systems/sidero/commit/5527432ddfc0fca78ad8bd9a1034a431117fcc06) fix: canonical links in docs site
* [`757706c`](https://github.com/talos-systems/sidero/commit/757706cae125abfdcb8ff49167539abd832b7c6a) feat: finalize update to Talos 1.2.0
* [`f9d12f3`](https://github.com/talos-systems/sidero/commit/f9d12f3c4635281015c6e3795aadfbc61dad805c) feat: update Talos to 1.2.0-beta.2
* [`44f1962`](https://github.com/talos-systems/sidero/commit/44f1962ceaac518f1e59205a43db3a1bf1e79cd6) fix: properly inspect readonly flag of the disk
* [`7cb28fc`](https://github.com/talos-systems/sidero/commit/7cb28fccfea5be7f5a22ddd9f895cd04cda68dc1) docs: handle latest version banner and canonical urls
* [`824d059`](https://github.com/talos-systems/sidero/commit/824d05920a8add702a9e705321585f247ea4ca48) fix: filter out SideroLink address from Machine addresses
* [`08cac99`](https://github.com/talos-systems/sidero/commit/08cac998b1e5228cf4daa0d1911ef7d0a8b4ae16) feat: update base Talos to 1.1.1
* [`3c65480`](https://github.com/talos-systems/sidero/commit/3c654806a52bbeb8604b78f27f2fafc0a75f2eb9) fix: skip read-only block devices on wiping
* [`4e714a1`](https://github.com/talos-systems/sidero/commit/4e714a1e4d657e01fcfe0f7a3e64ca58d6af89f5) fix: resolve wireguard endpoint to IP
* [`3eddbbd`](https://github.com/talos-systems/sidero/commit/3eddbbd69ecab561bc6583d75262442ed4b73228) fix: use http response writer directly
* [`a06135a`](https://github.com/talos-systems/sidero/commit/a06135a96e0d19c07d5e34908ebffee1ea15571d) docs: use ipxe-sanboot for raspberry pi 4's
* [`c0dfd23`](https://github.com/talos-systems/sidero/commit/c0dfd23fa5f8b775bf1ea84ce14c93927b821a88) test: integrate new controlplane loadbalancer
* [`fd0086d`](https://github.com/talos-systems/sidero/commit/fd0086d6f0579f3888cb0dfd8476b6c91cdd5c28) fix: append Talos default kernel args even if there is something defined
* [`6c6b354`](https://github.com/talos-systems/sidero/commit/6c6b3546e40665b605fa44cc5508647c90d20223) docs: edit for clarity and conciseness
* [`56d27ce`](https://github.com/talos-systems/sidero/commit/56d27ce71876df1626edcc56cb19c5c7b53d5aca) chore: bump dependencies
* [`b8843e2`](https://github.com/talos-systems/sidero/commit/b8843e2f875c0dc73a3acdba5ed65a3d9a869e04) chore: bump dependencies
* [`03a5c90`](https://github.com/talos-systems/sidero/commit/03a5c90b2c961d92709e3ff318e340ec6a5d0105) chore: bump dependencies
* [`6092854`](https://github.com/talos-systems/sidero/commit/609285488e13423626dc9131fef5c77e91ab9f5a) docs: fix spelling mistake
* [`c29d464`](https://github.com/talos-systems/sidero/commit/c29d4645d9dbb3715e6829948ba42fc0f518817b) feat: make MetalMachineTemplate immutable
* [`835d5cf`](https://github.com/talos-systems/sidero/commit/835d5cf227961d9a423985277d80767407dc7eac) docs: fix links and getting started overview
* [`7c7a377`](https://github.com/talos-systems/sidero/commit/7c7a37785711a5a90afd701921995ba88755ccce) feat: add extended hardware information to Server and ServerClass CRDs
* [`e4bb416`](https://github.com/talos-systems/sidero/commit/e4bb4165e9d25cceabd2e4cd73e5a53b4d5756ef) feat: restructure HW information in `Server` resources (v1alpha2)
* [`6c81518`](https://github.com/talos-systems/sidero/commit/6c81518b172ba09efbc46f15e33706d7340cdecf) docs: fork docs for Sidero 0.6
* [`061ee8e`](https://github.com/talos-systems/sidero/commit/061ee8e57debf11457e732f74be6e654e31e0203) refactor: remove pre-ServerBinding leftovers, use remote cached clients
* [`24449aa`](https://github.com/talos-systems/sidero/commit/24449aa3993f75299a8cf2113aa82a8aba65e551) chore: bump dependencies
* [`fd34779`](https://github.com/talos-systems/sidero/commit/fd34779bf7ada1914e320997c73fecc624813e8a) docs: add info about GitHub Org rename to v0.5 docs
* [`c86953e`](https://github.com/talos-systems/sidero/commit/c86953e11c0f9a82e7279b1f3977d869a33870bb) docs: update algolia info
* [`511ddfc`](https://github.com/talos-systems/sidero/commit/511ddfc0828a3e9b5ac26814d34ae9c5b117071a) feat: enable webhooks in `sidero-controller-manager`
* [`37a1d52`](https://github.com/talos-systems/sidero/commit/37a1d526962c0d3f6cbef28a43c6ca86d555df81) docs: add sitemap override
* [`54f896d`](https://github.com/talos-systems/sidero/commit/54f896dd6750874a2d065dc00ec96a56494a3f9d) feat: allow configuring Sidero deployment strategy
* [`3be5e6e`](https://github.com/talos-systems/sidero/commit/3be5e6e817fc8db6a8e4e2566cbad311d54ee7ad) chore: bump dependencies
* [`b7cc8b2`](https://github.com/talos-systems/sidero/commit/b7cc8b2e40d687b51558bb73efc014fc53e60183) docs: switch code block to text for DHCP prereq
* [`15f6730`](https://github.com/talos-systems/sidero/commit/15f67308f0bdbd52a4565deff8923d67a4c4c755) chore: use consistent naming for imports
* [`8985a04`](https://github.com/talos-systems/sidero/commit/8985a04231ed91f79477d5556a3256919287efb4) fix: ipxe prompt on ARM64
* [`9f9c922`](https://github.com/talos-systems/sidero/commit/9f9c92205c40e289975e9470625f4da0bf6cad64) docs: move docs to hugo
* [`2a475db`](https://github.com/talos-systems/sidero/commit/2a475db79ee6df559cecd082581fb7ba67d4f35b) feat: update to Talos 1.0
* [`7a1a930`](https://github.com/talos-systems/sidero/commit/7a1a930989d1b93667645c8b1d075a99dd31d799) chore: update more registries
* [`9438ca8`](https://github.com/talos-systems/sidero/commit/9438ca8e8aaafff3270bac85c5982918f3333278) chore: fix gpg check and sfyra build
* [`5fb30a9`](https://github.com/talos-systems/sidero/commit/5fb30a996c3ed4a4239ddfd48c67dd8f73f4d0b5) docs: update README for 0.5 version compatibility
* [`4d11603`](https://github.com/talos-systems/sidero/commit/4d116038e09d2c3aa62486ae4eac4f8ea863f1df) chore: bump cert-manager to v1
* [`13e11d3`](https://github.com/talos-systems/sidero/commit/13e11d3320c529c0934f0b2e5e5b80c0e36a14cc) chore: bump dependencies
* [`d3e75dc`](https://github.com/talos-systems/sidero/commit/d3e75dc9beff977b3c0e754805ccf94fb2e77001) docs: make 0.5 docs the default, add what's new
</p>
</details>

### Dependency Changes

* **github.com/go-logr/logr**                       v0.4.0 -> v1.2.3
* **github.com/google/go-cmp**                      v0.5.9 **_new_**
* **github.com/grpc-ecosystem/go-grpc-middleware**  v1.3.0 -> v1.4.0
* **github.com/insomniacslk/dhcp**                  74ae03f2425e **_new_**
* **github.com/jsimonetti/rtnetlink**               v1.3.2 **_new_**
* **github.com/siderolabs/gen**                     v0.4.3 **_new_**
* **github.com/siderolabs/go-blockdevice**          v0.4.4 **_new_**
* **github.com/siderolabs/go-cmd**                  v0.1.1 **_new_**
* **github.com/siderolabs/go-debug**                v0.2.2 **_new_**
* **github.com/siderolabs/go-kmsg**                 v0.1.3 **_new_**
* **github.com/siderolabs/go-procfs**               v0.1.1 **_new_**
* **github.com/siderolabs/go-retry**                v0.3.2 **_new_**
* **github.com/siderolabs/go-smbios**               v0.3.2 **_new_**
* **github.com/siderolabs/grpc-proxy**              v0.4.0 **_new_**
* **github.com/siderolabs/siderolink**              v0.3.1 **_new_**
* **github.com/siderolabs/talos/pkg/machinery**     v1.4.0 **_new_**
* **github.com/stretchr/testify**                   v1.7.0 -> v1.8.2
* **go.uber.org/zap**                               v1.20.0 -> v1.24.0
* **golang.org/x/net**                              cd36cc0744dd -> v0.9.0
* **golang.org/x/sync**                             036812b2e83c -> v0.1.0
* **golang.org/x/sys**                              99c3d69c2c27 -> v0.7.0
* **golang.zx2c4.com/wireguard/wgctrl**             daad0b7ba671 -> 9c5414ab4bde
* **google.golang.org/grpc**                        v1.44.0 -> v1.54.0
* **google.golang.org/protobuf**                    v1.27.1 -> v1.30.0
* **gopkg.in/yaml.v3**                              v3.0.1 **_new_**
* **k8s.io/api**                                    v0.22.2 -> v0.26.1
* **k8s.io/apiextensions-apiserver**                v0.22.2 -> v0.26.1
* **k8s.io/apimachinery**                           v0.22.2 -> v0.26.1
* **k8s.io/client-go**                              v0.22.2 -> v0.26.1
* **k8s.io/utils**                                  cb0fa318a74b -> 99ec85e7a448
* **sigs.k8s.io/cluster-api**                       v1.0.4 -> v1.4.1
* **sigs.k8s.io/controller-runtime**                v0.10.3 -> v0.14.5

Previous release can be found at [v0.5.0](https://github.com/talos-systems/sidero/releases/tag/v0.5.0)

## [Sidero 0.5.0-alpha.2](https://github.com/talos-systems/sidero/releases/tag/v0.5.0-alpha.2) (2022-02-04)

Welcome to the v0.5.0-alpha.2 release of Sidero!  
*This is a pre-release of Sidero*



Please try out the release binaries and report any issues at
https://github.com/talos-systems/sidero/issues.

### Machine Addresses

Sidero now populates `MetalMachine` addresses with the ones discovered from Siderolink server events.
Which is then propagated to CAPI `Machine` resources.

Requires Talos >= v0.14.


### iPXE Boot From Disk Method

iPXE boot from disk method can now be set not only on the global level, but also in the Server and ServerClass specs.


### CAPI v1beta1

This release of CACPPT brings compatibility with CAPI v1beta1.


### New `MetalMachines` Conditions

New set of conditions is now available which can simplify cluster troubleshooting:

- `TalosConfigLoaded` is set to false when the config load has failed.
- `TalosConfigValidated` is set to false when the config validation
fails on the node.
- `TalosInstalled` is set to true/false when talos installer finishes.


### IPMI PXE Method

IPMI PXE method (UEFI, BIOS) can now be configured with `SIDERO_CONTROLLER_MANAGER_IPMI_PXE_METHOD` while installing Sidero.


### Retry PXE Boot

Sidero server controller now keeps track of Talos installation progress.
Now the node will be PXE booted until Talos installation succeeds.


### Siderolink

Sidero now connects to all servers using Siderolink.
This enables streaming of all dmesg logs and events back to sidero.

All server logs can now be viewed by getting logs of one of the container of the `sidero-controller-manager`:

```
kubectl logs -f -n sidero-system deployment/sidero-controller-manager serverlogs
```

Events:

```
kubectl logs -f -n sidero-system deployment/sidero-controller-manager serverevents
```


### Contributors

* Andrey Smirnov
* Artem Chernyshev
* Michal Witkowski
* Andrey Smirnov
* Noel Georgi
* Spencer Smith
* Andrey Smirnov
* Steve Francis
* Charlie Haley
* Daniel Low
* Jack Wink
* Rubens Farias
* Serge Logvinov
* Seán C McCord

### Changes
<details><summary>49 commits</summary>
<p>

* [`0a967a8`](https://github.com/talos-systems/sidero/commit/0a967a8c3fce2fca2fbfdf35ef2dc2ffe90932e9) feat: configure boot from disk method on Server/ServerClass level
* [`7912509`](https://github.com/talos-systems/sidero/commit/7912509347539eca3de691ec14ca03ee2a721309) refactor: cleanup and deduplicate the code which matches SideroLink IPs
* [`003f6a7`](https://github.com/talos-systems/sidero/commit/003f6a73ec1486c3e8a8ba1930a15358ae3ca583) fix: don't boot from not ready environments
* [`e44f350`](https://github.com/talos-systems/sidero/commit/e44f350d08ec51c8cfd091fc3ef28f5e61eae8a1) fix: use controller-runtime standard healthz endpoints
* [`c73d8e5`](https://github.com/talos-systems/sidero/commit/c73d8e52b19314d11079cb56aa42813c4cac6b05) docs: update to Sidero from Talos
* [`5e3f68d`](https://github.com/talos-systems/sidero/commit/5e3f68ddac339ad471625164a3c04ae1ee0ffc4e) fix: add move label to BMC secrets
* [`f28d7b0`](https://github.com/talos-systems/sidero/commit/f28d7b0b53d44ae805d927c14997c389fd044525) docs: update README and index page
* [`194e870`](https://github.com/talos-systems/sidero/commit/194e87069b01761f856c6ca5e3d834820ff6fcb4) chore: bump dependencies
* [`b30fbe4`](https://github.com/talos-systems/sidero/commit/b30fbe4317f3cc832ce6d463c84aeba7d63ce5cb) feat: set server PXEBooted condition only after Talos gets installed
* [`34f7822`](https://github.com/talos-systems/sidero/commit/34f7822c3f243eae9ffd13bfdcbc83bb31178421) docs: fixes to the homepage, footer, log
* [`682594c`](https://github.com/talos-systems/sidero/commit/682594c4fd6edb24150f4d0c1cb65493286eb01f) docs: update index.md and sync across versions
* [`dea2179`](https://github.com/talos-systems/sidero/commit/dea21796d4f5600522cf7bf0f3a4d4409ff01ab4) refactor: rewrite conditions update function in the adapter
* [`59ed3cd`](https://github.com/talos-systems/sidero/commit/59ed3cde2bda13c15e8d33d2f19d9eb8afe1e403) release(v0.5.0-alpha.1): prepare release
* [`1f7588f`](https://github.com/talos-systems/sidero/commit/1f7588f1fa9f99db7ddf0b7ce1e5bb74b117940e) docs: update office hours
* [`fe41335`](https://github.com/talos-systems/sidero/commit/fe41335457a76f429177a24627044758773fe557) feat: introduce new conditions in the `metalmachine`
* [`6454dee`](https://github.com/talos-systems/sidero/commit/6454dee29c994b57e09fe9a7673dd95274366820) feat: populate addresses and hostname in ServerBindings
* [`d69264f`](https://github.com/talos-systems/sidero/commit/d69264f7b7ee0c4d32e8de054a6cf70017a77f60) docs: fix patching examples
* [`04d90fd`](https://github.com/talos-systems/sidero/commit/04d90fdc616c3a60302523e326f978b216d5e062) docs: add patching examples
* [`41b7451`](https://github.com/talos-systems/sidero/commit/41b7451b50d13ac979e1df4a16299bf424b2be2e) docs: update docs for upcoming Sidero 0.4.1 release
* [`d5f8f4e`](https://github.com/talos-systems/sidero/commit/d5f8f4e96de346385bc5f9f014d989e2726d35c6) release(v0.5.0-alpha.0): prepare release
* [`229dae1`](https://github.com/talos-systems/sidero/commit/229dae1fd6e90e3acca5eee3de05ba05dd39c809) fix: ignore BMC info if username or password are not set
* [`650743a`](https://github.com/talos-systems/sidero/commit/650743ad9a16d326bf9b327d45126a7a8f5e8c21) fix: use environment variables in the ipmitool
* [`ed39a3c`](https://github.com/talos-systems/sidero/commit/ed39a3c4018b4db9c8fa5debad636cc812e2ca55) fix: ensure we setup BMC info *after* disk wiping
* [`025ff04`](https://github.com/talos-systems/sidero/commit/025ff047d0a2d0b239d4ecff801d37d145040bae) feat: additional printer colums
* [`189793e`](https://github.com/talos-systems/sidero/commit/189793e47fb57757b8de39814d3ef831919b4bb2) fix: wait for link up in iPXE script
* [`9a27861`](https://github.com/talos-systems/sidero/commit/9a2786123a598cc3dde6e96c308b73c332f1d70a) fix: make iPXE script replace script on chain request
* [`1bbe3be`](https://github.com/talos-systems/sidero/commit/1bbe3be26fd5a925dec5f8de41ebafd1d852a5f2) feat: extend information printed in the iPXE script, add retries
* [`4cfdeda`](https://github.com/talos-systems/sidero/commit/4cfdedaf97e267e308e4641d1df28fbbbc36d922) feat: provide a way to configure IPMI PXE method
* [`2ff14c4`](https://github.com/talos-systems/sidero/commit/2ff14c4528e31bd5964073b9335791f2d683f021) docs: reflect sidero runs on ARM
* [`274ae33`](https://github.com/talos-systems/sidero/commit/274ae33fc7c3b4b8f5b517914f730a4db3a9840a) fix: check for server power state when in use
* [`d0df929`](https://github.com/talos-systems/sidero/commit/d0df929eb1d1067636caaa2a95c7126be6c02713) feat: implement events manager container
* [`ab12b81`](https://github.com/talos-systems/sidero/commit/ab12b81ef00ad1762aaa251fbaa6b022c298ae62) feat: implement Talos kernel log receiver
* [`5bf7c21`](https://github.com/talos-systems/sidero/commit/5bf7c21f8f61002fd655863cb0ed6ef8f8b649fe) docs: fix clusterctl usage
* [`e77bf54`](https://github.com/talos-systems/sidero/commit/e77bf54a31076dc70dc726f30e924e57e25a14ec) feat: support cordoning server
* [`ab29103`](https://github.com/talos-systems/sidero/commit/ab291034e5a6bafef7eaea86dcb6d594a8afc420) feat: implement SideroLink
* [`adc73b6`](https://github.com/talos-systems/sidero/commit/adc73b67c5ae3b2302be0ced3b87913de4f15b0e) fix: update CAPI resources versions to v1alpha4 in the cluster template
* [`04dbaf0`](https://github.com/talos-systems/sidero/commit/04dbaf065a20e1d5dfbe745d9126b5a174456afd) test: fix Sidero components version in sfyra docs
* [`acb5f57`](https://github.com/talos-systems/sidero/commit/acb5f57f65a0a226d604ad124f189afe5752974a) feat: automatically append talos.config to the Environment
* [`0e7f8a6`](https://github.com/talos-systems/sidero/commit/0e7f8a6820dc77a28d0d264b7e2dd216dab54eb0) docs: metadata server -> sidero
* [`a826281`](https://github.com/talos-systems/sidero/commit/a82628186d84a4f2f49c51b1cb219cf482a3653e) fix: hide password from ipmitool args
* [`ef270df`](https://github.com/talos-systems/sidero/commit/ef270df5b7af6e143bf1449e34a0a577441ab03f) docs: fork docs for 0.5
* [`a0bf382`](https://github.com/talos-systems/sidero/commit/a0bf3828623c3c4ba425601767960f64fc9d85e6) docs: dhcp option-93
* [`bbbe814`](https://github.com/talos-systems/sidero/commit/bbbe814eb98884c525301e49db09b65bc3f0a7b3) chore: conformance check
* [`701d01b`](https://github.com/talos-systems/sidero/commit/701d01ba370c0bd6baa173207575bdaf53c72970) fix: drop into the agent for non-allocated servers
* [`b0e3611`](https://github.com/talos-systems/sidero/commit/b0e3611d2383061533e10d3bc1e642c99e9f70f9) docs: update help message for extra kernel args
* [`bb87567`](https://github.com/talos-systems/sidero/commit/bb87567e242c684977a3e688fc9553f6a77c81e6) chore: bump dependencies
* [`62ab9be`](https://github.com/talos-systems/sidero/commit/62ab9beacfa33e8fb348b52582e63d6234af9146) fix: update extension for controlplane.yam l talosctl generates YAML files with the .yaml extension, update to the apply-config command to reflect this
* [`0511d78`](https://github.com/talos-systems/sidero/commit/0511d78ef29651bce4d1b0ce8d24592582cfdb2e) feat: allow defining http server and api server ports separately
* [`432ca2a`](https://github.com/talos-systems/sidero/commit/432ca2a467a7fb4ee3cc448b7508872aa2674654) docs: create `v0.4` docs and set them as latest
</p>
</details>

### Changes since v0.5.0-alpha.1
<details><summary>12 commits</summary>
<p>

* [`0a967a8`](https://github.com/talos-systems/sidero/commit/0a967a8c3fce2fca2fbfdf35ef2dc2ffe90932e9) feat: configure boot from disk method on Server/ServerClass level
* [`7912509`](https://github.com/talos-systems/sidero/commit/7912509347539eca3de691ec14ca03ee2a721309) refactor: cleanup and deduplicate the code which matches SideroLink IPs
* [`003f6a7`](https://github.com/talos-systems/sidero/commit/003f6a73ec1486c3e8a8ba1930a15358ae3ca583) fix: don't boot from not ready environments
* [`e44f350`](https://github.com/talos-systems/sidero/commit/e44f350d08ec51c8cfd091fc3ef28f5e61eae8a1) fix: use controller-runtime standard healthz endpoints
* [`c73d8e5`](https://github.com/talos-systems/sidero/commit/c73d8e52b19314d11079cb56aa42813c4cac6b05) docs: update to Sidero from Talos
* [`5e3f68d`](https://github.com/talos-systems/sidero/commit/5e3f68ddac339ad471625164a3c04ae1ee0ffc4e) fix: add move label to BMC secrets
* [`f28d7b0`](https://github.com/talos-systems/sidero/commit/f28d7b0b53d44ae805d927c14997c389fd044525) docs: update README and index page
* [`194e870`](https://github.com/talos-systems/sidero/commit/194e87069b01761f856c6ca5e3d834820ff6fcb4) chore: bump dependencies
* [`b30fbe4`](https://github.com/talos-systems/sidero/commit/b30fbe4317f3cc832ce6d463c84aeba7d63ce5cb) feat: set server PXEBooted condition only after Talos gets installed
* [`34f7822`](https://github.com/talos-systems/sidero/commit/34f7822c3f243eae9ffd13bfdcbc83bb31178421) docs: fixes to the homepage, footer, log
* [`682594c`](https://github.com/talos-systems/sidero/commit/682594c4fd6edb24150f4d0c1cb65493286eb01f) docs: update index.md and sync across versions
* [`dea2179`](https://github.com/talos-systems/sidero/commit/dea21796d4f5600522cf7bf0f3a4d4409ff01ab4) refactor: rewrite conditions update function in the adapter
</p>
</details>

### Changes from talos-systems/cluster-api-bootstrap-provider-talos
<details><summary>9 commits</summary>
<p>

* [`1776117`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/17761177f0cfe157498acb6440d07feac8f9a6f1) release(v0.5.1): prepare release
* [`1b88f9f`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/1b88f9f8a37c3c3fefe3d84fc310c44d1dcb0ded) feat: update Talos to 0.14.0
* [`6d27c57`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/6d27c57584b99ac9aab5313191c701ccd780bc5d) release(v0.5.0): prepare release
* [`f6dc0a3`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/f6dc0a3372dba82306a4abc9b2a064f1e337421c) fix: regenerate manifests
* [`2a4115f`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/2a4115f1211a20e5058a7b0430c4dc4081acfcfe) release(v0.5.0-alpha.0): prepare release
* [`d124c07`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/d124c072c9db8d402b353a73646d2d197bae76a4) docs: update README with usage and compatibility matrix
* [`20792f3`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/20792f345b7ff3c8ffa9d65c9ca8dcab1932f49e) feat: generate talosconfig as a secret with proper endpoints
* [`abd206f`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/abd206fd8a98f5478f8ffd0f8686e32be3b7defe) feat: update to CAPI v1.0.x contract (v1beta1)
* [`b7faf9e`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/b7faf9e730b7c9f50ffa94be194ddcf908708a2c) feat: update Talos machinery to 0.13.0
</p>
</details>

### Changes from talos-systems/cluster-api-control-plane-provider-talos
<details><summary>19 commits</summary>
<p>

* [`adea239`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/adea239aafeb8b274d3b15e39e3106a96f22c1fc) release(v0.4.3): prepare release
* [`efa0345`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/efa03451f88f7c0b1eb6b94302f674466660a9af) fix: fall back to old scheme of getting talsoconfig for older templates
* [`89f793e`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/89f793ef54173d613949af715c95aa0581630758) release(v0.4.2): prepare release
* [`a77ddde`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/a77ddde2607165396c254c292de4e726c2c29f38) refactor: change reconcile loop flow
* [`ea7842f`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/ea7842f6fabefa1775c0a2c7bd43e1a7e0615fe7) release(v0.4.1): prepare release
* [`7f63ad0`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/7f63ad0a391dcf0901edf9009717fb57f503f652) fix: avoid long backoff when trying to bootstrap the cluster
* [`8fc9a6c`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/8fc9a6c3bf7e9de5d074c50c7e7e70d025a5369c) release(v0.4.0): prepare release
* [`b63f1d2`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/b63f1d211002b6ee048de304987151c7deda5db1) release(v0.4.0-beta.2): prepare release
* [`f5f5b2d`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/f5f5b2d441ccfb939e2573ef780a69af371775da) fix: patch the status and use APIReader to get resource
* [`d606d32`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/d606d32a79a1ba567748c782c6ffc4155ce0c81a) release(v0.4.0-beta.1): prepare release
* [`333fc02`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/333fc0221a68e978c5d31399e9befe7d6b29aebe) fix: ensure that bootstrap is called only a single time
* [`77b0bba`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/77b0bba8d7c24026d458acb14c9d2793c0450aa8) test: update templates to v1beta1
* [`a5af5e4`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/a5af5e4b09f8450b73b825d95e49b75b045cb47b) release(v0.4.0-beta.0): prepare release
* [`80b24a0`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/80b24a0abb4a9176a9f1635891a236f299d7dc64) fix: introduce a separate step for release builds
* [`a24dad3`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/a24dad3328d52a3744f685ddde913d66dd17b176) fix: do not allow scaling down controlplane to zero
* [`8a73e6a`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/8a73e6a26e52151b1dd4604c4d0737036f119c30) feat: get rid of init nodes and use bootstrap API to setup cluster
* [`205f4be`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/205f4be2057b3ea81c4dcf47004db6864ff31801) release(v0.4.0-alpha.0): prepare release
* [`b8db449`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/b8db4492d55f910e8a7d2a3b69ab08740963683e) fix: properly pick talos client configuration
* [`61fb582`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/61fb5826391e4434b64619f0590683f7fa7b82b6) feat: support clusterapi v1beta1
</p>
</details>

### Changes from talos-systems/go-blockdevice
<details><summary>2 commits</summary>
<p>

* [`15b182d`](https://github.com/talos-systems/go-blockdevice/commit/15b182db0cd233b163ed83d1724c7e28cf29d71a) fix: return partition table not exist when trying to read an empty dev
* [`b9517d5`](https://github.com/talos-systems/go-blockdevice/commit/b9517d51120d385f97b0026f99ce3c4782940c37) fix: resize partition
</p>
</details>

### Changes from talos-systems/go-smbios
<details><summary>1 commit</summary>
<p>

* [`fd5ec8c`](https://github.com/talos-systems/go-smbios/commit/fd5ec8ce4873790b7fbd46dba9d7f49c9de7176a) fix: remove useless (?) goroutines leading to data race error
</p>
</details>

### Changes from talos-systems/grpc-proxy
<details><summary>44 commits</summary>
<p>

* [`ca3bc61`](https://github.com/talos-systems/grpc-proxy/commit/ca3bc6131f052aa000517339211335aaa4ebb640) fix: ignore some errors so that we don't spam the logs
* [`5c579a7`](https://github.com/talos-systems/grpc-proxy/commit/5c579a7a61475bde3ec9c1efe000d2a55e2a3cb2) feat: allow different formats for messages streaming/unary
* [`6c9f7b3`](https://github.com/talos-systems/grpc-proxy/commit/6c9f7b399173dd5769dbc4e8e366e78f05cead85) fix: allow mode to be set for each request being proxied
* [`cc91c09`](https://github.com/talos-systems/grpc-proxy/commit/cc91c09782824e261bf1c861961a272aedb2b123) refactor: provide better public API, enforce proxying mode
* [`d8d3a75`](https://github.com/talos-systems/grpc-proxy/commit/d8d3a751d1e71d006ba90379eed388c487bbb246) chore: update import paths after repo move
* [`dbf07a4`](https://github.com/talos-systems/grpc-proxy/commit/dbf07a4d9e16fe3cf7407b9921c1746aa24ffaf6) Merge pull request  [#7](https://github.com/talos-systems/grpc-proxy/pull/7) from smira/one2many-4
* [`fc0d27d`](https://github.com/talos-systems/grpc-proxy/commit/fc0d27dc6b5b9db35173f3e78778784a9e7c95bf) More tests, small code fixes, updated README.
* [`d9ce0b1`](https://github.com/talos-systems/grpc-proxy/commit/d9ce0b1053a7f15ea65bf46e94cfe4154493bad7) Merge pull request  [#6](https://github.com/talos-systems/grpc-proxy/pull/6) from smira/one2many-3
* [`2d37ba4`](https://github.com/talos-systems/grpc-proxy/commit/2d37ba444528a00f988671f3a01666e692739a37) Support for one2many streaming calls, tests.
* [`817b035`](https://github.com/talos-systems/grpc-proxy/commit/817b03553ed7d97bd0da09283776d54592d7b5d4) Merge pull request  [#5](https://github.com/talos-systems/grpc-proxy/pull/5) from smira/one2many-2
* [`436b338`](https://github.com/talos-systems/grpc-proxy/commit/436b3383a39fd860f3b2379ffab80a44ae1809f7) More unary one-2-many tests, error propagation.
* [`1f0cb46`](https://github.com/talos-systems/grpc-proxy/commit/1f0cb466268f046e8e9fb78b1902411ac3a753ba) Merge pull request  [#4](https://github.com/talos-systems/grpc-proxy/pull/4) from smira/one2many-1
* [`992a975`](https://github.com/talos-systems/grpc-proxy/commit/992a975ccf0b97e4be329c84bd3018652e8e50ae) Proxying one to many: first iteration
* [`a0988ff`](https://github.com/talos-systems/grpc-proxy/commit/a0988ff2b29839892a7913acd76f26f4e7edcc3a) Merge pull request  [#3](https://github.com/talos-systems/grpc-proxy/pull/3) from smira/small-fixups
* [`e3111ef`](https://github.com/talos-systems/grpc-proxy/commit/e3111ef2c16f0ee4bba597a2ab1ab6a2818c2734) Small fixups in preparation to add one-to-many proxying.
* [`6d76ffc`](https://github.com/talos-systems/grpc-proxy/commit/6d76ffcff89f6636d3689ed1c9b0eebe87722114) Merge pull request  [#2](https://github.com/talos-systems/grpc-proxy/pull/2) from smira/backend-concept
* [`2aad63a`](https://github.com/talos-systems/grpc-proxy/commit/2aad63ac5bae09232ea5ac80b42338e9e3af67c4) Add concept of a 'Backend', but still one to one proxying
* [`7cc4610`](https://github.com/talos-systems/grpc-proxy/commit/7cc46101114a2779d6393e0e8f841bf3febb2753) Merge pull request  [#1](https://github.com/talos-systems/grpc-proxy/pull/1) from smira/build
* [`37f01f3`](https://github.com/talos-systems/grpc-proxy/commit/37f01f3aab3b978a8fecb428fca4d4c722141229) Rework build to use GitHub Actions, linting updates.
* [`0f1106e`](https://github.com/talos-systems/grpc-proxy/commit/0f1106ef9c766333b9acb4b81e705da4bade7215) Move error checking further up (#34)
* [`d5b35f6`](https://github.com/talos-systems/grpc-proxy/commit/d5b35f634383bf8931f8798797daaf9c1a59235e) Update gRPC and fix tests (#27)
* [`67591eb`](https://github.com/talos-systems/grpc-proxy/commit/67591eb23c48346a480470e462289835d96f70da) Break StreamDirector interface, fix metadata propagation for gRPC-Go>1.5. (#20)
* [`97396d9`](https://github.com/talos-systems/grpc-proxy/commit/97396d94749c00db659393ba5123f707062f829f) Merge pull request  [#11](https://github.com/talos-systems/grpc-proxy/pull/11) from mwitkow/fix-close-bug
* [`3fcbd37`](https://github.com/talos-systems/grpc-proxy/commit/3fcbd3737ec6baff505795417e48f162a7a3183c) fixup closing conns
* [`a8f5f87`](https://github.com/talos-systems/grpc-proxy/commit/a8f5f87a2f5e6bc3643b78d64594195b2395a238) fixup tests, extend readme
* [`428fa1c`](https://github.com/talos-systems/grpc-proxy/commit/428fa1c450320041e0ad8e251d6aed435401174e) Fix a channel closing bug
* [`af55d61`](https://github.com/talos-systems/grpc-proxy/commit/af55d612de6c5723a5a59340704db7bc771023ff) Merge pull request  [#10](https://github.com/talos-systems/grpc-proxy/pull/10) from mwitkow/bugfix/streaming-fix
* [`de4d3db`](https://github.com/talos-systems/grpc-proxy/commit/de4d3db538565636e1e977102f6f0bd1ed0ce9c2) remove spurious printfs
* [`84242c4`](https://github.com/talos-systems/grpc-proxy/commit/84242c4e690da18d16d2ab8f2fa47e45986220b6) fix the "i don't know who finished" case
* [`9b22f41`](https://github.com/talos-systems/grpc-proxy/commit/9b22f41d8535fa3e40908c78ae66066c7972b6d9) fix full duplex streaming
* [`c2f7c98`](https://github.com/talos-systems/grpc-proxy/commit/c2f7c98b0b6cd180659aed31e98cbbc18d616b1c) update readme
* [`d654141`](https://github.com/talos-systems/grpc-proxy/commit/d654141edcb92b7fa2bba9d3e690e569c72f8e9d) update README
* [`f457856`](https://github.com/talos-systems/grpc-proxy/commit/f4578565f2d34dc89774128db2bfda3a328cba40) move to proxy subdirectory
* [`4889d78`](https://github.com/talos-systems/grpc-proxy/commit/4889d78e468681601b8229c81807dcf37b00ff63) Add fixup scripts
* [`ef60a37`](https://github.com/talos-systems/grpc-proxy/commit/ef60a37547d137e52873be183f2d7a5626d7c034) version 2 of the grpc-proxy, this time with fewer grpc upstream deps
* [`07aeac1`](https://github.com/talos-systems/grpc-proxy/commit/07aeac13e988c0c0b3a886c79972e20408a765e0) Merge pull request  [#2](https://github.com/talos-systems/grpc-proxy/pull/2) from daniellowtw/master
* [`e5c3df5`](https://github.com/talos-systems/grpc-proxy/commit/e5c3df5b2f0a1ffc4cb755cbe6b30b435e35de37) Fix compatibility with latest grpc library
* [`52be0a5`](https://github.com/talos-systems/grpc-proxy/commit/52be0a559a85f0e2480bde6725f3f144396aa6ef) bugfix: fix gRPC Java deadlock, due to different dispatch logic
* [`822df7d`](https://github.com/talos-systems/grpc-proxy/commit/822df7d86b556b703fc11798a3bdcbaeb60c18a6) Fix reference to mwitkow.
* [`28341d1`](https://github.com/talos-systems/grpc-proxy/commit/28341d171dd4c1a52f46371ddfb5fd2240b79731) move out forward logic to method, allowing for use as `grpc.Server` not found handler.
* [`89e28b4`](https://github.com/talos-systems/grpc-proxy/commit/89e28b42ee9dda8e36522b77e3771d9debc645e0) add reference to upstream grpc bug
* [`00dd588`](https://github.com/talos-systems/grpc-proxy/commit/00dd588ae68adf4187a7fca87db45a73af4c834d) merge upstream `grpc.Server` changes changing the dispatch logic
* [`77edc97`](https://github.com/talos-systems/grpc-proxy/commit/77edc9715de187dcbc9969e2f0e8a04d2087fd13) move to upstream `protobuf` from `gogo`
* [`db71c3e`](https://github.com/talos-systems/grpc-proxy/commit/db71c3e7e812db8d75cb282dac38d953fcb436b3) initial commit, tested and working.
</p>
</details>

### Changes from talos-systems/net
<details><summary>2 commits</summary>
<p>

* [`b4b7181`](https://github.com/talos-systems/net/commit/b4b718179a1aa68e4f54422baf08ca3761723d2d) feat: add a way to filter list of IPs for the machine
* [`0abe5bd`](https://github.com/talos-systems/net/commit/0abe5bdae8f85e4e976bc4d90e95dcb4be8fb853) feat: implement FilterIPs function
</p>
</details>

### Changes from talos-systems/siderolink
<details><summary>7 commits</summary>
<p>

* [`9902ad2`](https://github.com/talos-systems/siderolink/commit/9902ad2774f0655e050233854b9d28dad0431f6c) feat: pass request context and node address to the events sink adapter
* [`d0612a7`](https://github.com/talos-systems/siderolink/commit/d0612a724a1b1336a2bc6a99ed3178e3e40f6d9b) refactor: pass in listener to the log receiver
* [`d86cdd5`](https://github.com/talos-systems/siderolink/commit/d86cdd59ee7a0e0504b739a913991c272c7fb3f5) feat: implement logreceiver for kernel logs
* [`f7cadbc`](https://github.com/talos-systems/siderolink/commit/f7cadbcdfbb84d367e27b5af32e89c138d72d9d7) fix: handle duplicate peer updates
* [`0755b24`](https://github.com/talos-systems/siderolink/commit/0755b24d4682410b251a2a9d662960da15153106) feat: initial implementation of SideroLink
* [`ee73ea9`](https://github.com/talos-systems/siderolink/commit/ee73ea9575a81be7685f24936b2c48a4508a159e) feat: add Talos events sink proto files and the reference implementation
* [`1e2cd9d`](https://github.com/talos-systems/siderolink/commit/1e2cd9d38621234a0a6010e33b1bab264f4d9bdf) Initial commit
</p>
</details>

### Dependency Changes

* **github.com/evanphx/json-patch**                                      v4.11.0 -> v5.6.0
* **github.com/grpc-ecosystem/go-grpc-middleware**                       v1.3.0 **_new_**
* **github.com/onsi/ginkgo**                                             v1.16.4 -> v1.16.5
* **github.com/onsi/gomega**                                             v1.16.0 -> v1.18.1
* **github.com/talos-systems/cluster-api-bootstrap-provider-talos**      v0.4.0 -> v0.5.1
* **github.com/talos-systems/cluster-api-control-plane-provider-talos**  v0.3.0 -> v0.4.3
* **github.com/talos-systems/go-blockdevice**                            v0.2.4 -> v0.2.5
* **github.com/talos-systems/go-smbios**                                 d3a32bea731a -> v0.1.1
* **github.com/talos-systems/grpc-proxy**                                v0.2.0 **_new_**
* **github.com/talos-systems/net**                                       v0.3.0 -> v0.3.1
* **github.com/talos-systems/siderolink**                                9902ad2774f0 **_new_**
* **go.uber.org/zap**                                                    v1.20.0 **_new_**
* **golang.org/x/net**                                                   853a461950ff -> cd36cc0744dd
* **golang.org/x/sys**                                                   39ccf1dd6fa6 -> 99c3d69c2c27
* **golang.zx2c4.com/wireguard/wgctrl**                                  daad0b7ba671 **_new_**
* **google.golang.org/grpc**                                             v1.41.0 -> v1.44.0
* **inet.af/netaddr**                                                    c74959edd3b6 **_new_**
* **k8s.io/utils**                                                       bdf08cb9a70a -> cb0fa318a74b
* **sigs.k8s.io/cluster-api**                                            v0.4.3 -> v1.0.4
* **sigs.k8s.io/controller-runtime**                                     v0.9.7 -> v0.10.3

Previous release can be found at [v0.4.0](https://github.com/talos-systems/sidero/releases/tag/v0.4.0)

## [Sidero 0.5.0-alpha.1](https://github.com/talos-systems/sidero/releases/tag/v0.5.0-alpha.1) (2022-01-11)

Welcome to the v0.5.0-alpha.1 release of Sidero!  
*This is a pre-release of Sidero*



Please try out the release binaries and report any issues at
https://github.com/talos-systems/sidero/issues.

### Machine Addresses

Sidero now populates `MetalMachine` addresses with the ones discovered from Siderolink server events.
Which is then propagated to CAPI `Machine` resources.

Requires Talos >= v0.14.


### CAPI v1beta1

This release of CACPPT brings compatibility with CAPI v1beta1.


### New `MetalMachines` Conditions

New set of conditions is now available which can simplify cluster troubleshooting:

- `TalosConfigLoaded` is set to false when the config load has failed.
- `TalosConfigValidated` is set to false when the config validation
fails on the node.
- `TalosInstalled` is set to true/false when talos installer finishes.


### IPMI PXE Method

IPMI PXE method (UEFI, BIOS) can now be configured with `SIDERO_CONTROLLER_MANAGER_IPMI_PXE_METHOD` while installing Sidero.


### Siderolink

Sidero now connects to all servers using Siderolink.
This enables streaming of all dmesg logs and events back to sidero.

All server logs can now be viewed by getting logs of one of the container of the `sidero-controller-manager`:

```
kubectl logs -f -n sidero-system deployment/sidero-controller-manager serverlogs
```

Events:

```
kubectl logs -f -n sidero-system deployment/sidero-controller-manager serverevents
```


### Contributors

* Andrey Smirnov
* Michal Witkowski
* Artem Chernyshev
* Andrey Smirnov
* Noel Georgi
* Andrey Smirnov
* Spencer Smith
* Charlie Haley
* Daniel Low
* Jack Wink
* Rubens Farias
* Serge Logvinov
* Seán C McCord

### Changes
<details><summary>36 commits</summary>
<p>

* [`1f7588f`](https://github.com/talos-systems/sidero/commit/1f7588f1fa9f99db7ddf0b7ce1e5bb74b117940e) docs: update office hours
* [`fe41335`](https://github.com/talos-systems/sidero/commit/fe41335457a76f429177a24627044758773fe557) feat: introduce new conditions in the `metalmachine`
* [`6454dee`](https://github.com/talos-systems/sidero/commit/6454dee29c994b57e09fe9a7673dd95274366820) feat: populate addresses and hostname in ServerBindings
* [`d69264f`](https://github.com/talos-systems/sidero/commit/d69264f7b7ee0c4d32e8de054a6cf70017a77f60) docs: fix patching examples
* [`04d90fd`](https://github.com/talos-systems/sidero/commit/04d90fdc616c3a60302523e326f978b216d5e062) docs: add patching examples
* [`41b7451`](https://github.com/talos-systems/sidero/commit/41b7451b50d13ac979e1df4a16299bf424b2be2e) docs: update docs for upcoming Sidero 0.4.1 release
* [`d5f8f4e`](https://github.com/talos-systems/sidero/commit/d5f8f4e96de346385bc5f9f014d989e2726d35c6) release(v0.5.0-alpha.0): prepare release
* [`229dae1`](https://github.com/talos-systems/sidero/commit/229dae1fd6e90e3acca5eee3de05ba05dd39c809) fix: ignore BMC info if username or password are not set
* [`650743a`](https://github.com/talos-systems/sidero/commit/650743ad9a16d326bf9b327d45126a7a8f5e8c21) fix: use environment variables in the ipmitool
* [`ed39a3c`](https://github.com/talos-systems/sidero/commit/ed39a3c4018b4db9c8fa5debad636cc812e2ca55) fix: ensure we setup BMC info *after* disk wiping
* [`025ff04`](https://github.com/talos-systems/sidero/commit/025ff047d0a2d0b239d4ecff801d37d145040bae) feat: additional printer colums
* [`189793e`](https://github.com/talos-systems/sidero/commit/189793e47fb57757b8de39814d3ef831919b4bb2) fix: wait for link up in iPXE script
* [`9a27861`](https://github.com/talos-systems/sidero/commit/9a2786123a598cc3dde6e96c308b73c332f1d70a) fix: make iPXE script replace script on chain request
* [`1bbe3be`](https://github.com/talos-systems/sidero/commit/1bbe3be26fd5a925dec5f8de41ebafd1d852a5f2) feat: extend information printed in the iPXE script, add retries
* [`4cfdeda`](https://github.com/talos-systems/sidero/commit/4cfdedaf97e267e308e4641d1df28fbbbc36d922) feat: provide a way to configure IPMI PXE method
* [`2ff14c4`](https://github.com/talos-systems/sidero/commit/2ff14c4528e31bd5964073b9335791f2d683f021) docs: reflect sidero runs on ARM
* [`274ae33`](https://github.com/talos-systems/sidero/commit/274ae33fc7c3b4b8f5b517914f730a4db3a9840a) fix: check for server power state when in use
* [`d0df929`](https://github.com/talos-systems/sidero/commit/d0df929eb1d1067636caaa2a95c7126be6c02713) feat: implement events manager container
* [`ab12b81`](https://github.com/talos-systems/sidero/commit/ab12b81ef00ad1762aaa251fbaa6b022c298ae62) feat: implement Talos kernel log receiver
* [`5bf7c21`](https://github.com/talos-systems/sidero/commit/5bf7c21f8f61002fd655863cb0ed6ef8f8b649fe) docs: fix clusterctl usage
* [`e77bf54`](https://github.com/talos-systems/sidero/commit/e77bf54a31076dc70dc726f30e924e57e25a14ec) feat: support cordoning server
* [`ab29103`](https://github.com/talos-systems/sidero/commit/ab291034e5a6bafef7eaea86dcb6d594a8afc420) feat: implement SideroLink
* [`adc73b6`](https://github.com/talos-systems/sidero/commit/adc73b67c5ae3b2302be0ced3b87913de4f15b0e) fix: update CAPI resources versions to v1alpha4 in the cluster template
* [`04dbaf0`](https://github.com/talos-systems/sidero/commit/04dbaf065a20e1d5dfbe745d9126b5a174456afd) test: fix Sidero components version in sfyra docs
* [`acb5f57`](https://github.com/talos-systems/sidero/commit/acb5f57f65a0a226d604ad124f189afe5752974a) feat: automatically append talos.config to the Environment
* [`0e7f8a6`](https://github.com/talos-systems/sidero/commit/0e7f8a6820dc77a28d0d264b7e2dd216dab54eb0) docs: metadata server -> sidero
* [`a826281`](https://github.com/talos-systems/sidero/commit/a82628186d84a4f2f49c51b1cb219cf482a3653e) fix: hide password from ipmitool args
* [`ef270df`](https://github.com/talos-systems/sidero/commit/ef270df5b7af6e143bf1449e34a0a577441ab03f) docs: fork docs for 0.5
* [`a0bf382`](https://github.com/talos-systems/sidero/commit/a0bf3828623c3c4ba425601767960f64fc9d85e6) docs: dhcp option-93
* [`bbbe814`](https://github.com/talos-systems/sidero/commit/bbbe814eb98884c525301e49db09b65bc3f0a7b3) chore: conformance check
* [`701d01b`](https://github.com/talos-systems/sidero/commit/701d01ba370c0bd6baa173207575bdaf53c72970) fix: drop into the agent for non-allocated servers
* [`b0e3611`](https://github.com/talos-systems/sidero/commit/b0e3611d2383061533e10d3bc1e642c99e9f70f9) docs: update help message for extra kernel args
* [`bb87567`](https://github.com/talos-systems/sidero/commit/bb87567e242c684977a3e688fc9553f6a77c81e6) chore: bump dependencies
* [`62ab9be`](https://github.com/talos-systems/sidero/commit/62ab9beacfa33e8fb348b52582e63d6234af9146) fix: update extension for controlplane.yam l talosctl generates YAML files with the .yaml extension, update to the apply-config command to reflect this
* [`0511d78`](https://github.com/talos-systems/sidero/commit/0511d78ef29651bce4d1b0ce8d24592582cfdb2e) feat: allow defining http server and api server ports separately
* [`432ca2a`](https://github.com/talos-systems/sidero/commit/432ca2a467a7fb4ee3cc448b7508872aa2674654) docs: create `v0.4` docs and set them as latest
</p>
</details>

### Changes since v0.5.0-alpha.0
<details><summary>6 commits</summary>
<p>

* [`1f7588f`](https://github.com/talos-systems/sidero/commit/1f7588f1fa9f99db7ddf0b7ce1e5bb74b117940e) docs: update office hours
* [`fe41335`](https://github.com/talos-systems/sidero/commit/fe41335457a76f429177a24627044758773fe557) feat: introduce new conditions in the `metalmachine`
* [`6454dee`](https://github.com/talos-systems/sidero/commit/6454dee29c994b57e09fe9a7673dd95274366820) feat: populate addresses and hostname in ServerBindings
* [`d69264f`](https://github.com/talos-systems/sidero/commit/d69264f7b7ee0c4d32e8de054a6cf70017a77f60) docs: fix patching examples
* [`04d90fd`](https://github.com/talos-systems/sidero/commit/04d90fdc616c3a60302523e326f978b216d5e062) docs: add patching examples
* [`41b7451`](https://github.com/talos-systems/sidero/commit/41b7451b50d13ac979e1df4a16299bf424b2be2e) docs: update docs for upcoming Sidero 0.4.1 release
</p>
</details>

### Changes from talos-systems/cluster-api-bootstrap-provider-talos
<details><summary>9 commits</summary>
<p>

* [`1776117`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/17761177f0cfe157498acb6440d07feac8f9a6f1) release(v0.5.1): prepare release
* [`1b88f9f`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/1b88f9f8a37c3c3fefe3d84fc310c44d1dcb0ded) feat: update Talos to 0.14.0
* [`6d27c57`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/6d27c57584b99ac9aab5313191c701ccd780bc5d) release(v0.5.0): prepare release
* [`f6dc0a3`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/f6dc0a3372dba82306a4abc9b2a064f1e337421c) fix: regenerate manifests
* [`2a4115f`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/2a4115f1211a20e5058a7b0430c4dc4081acfcfe) release(v0.5.0-alpha.0): prepare release
* [`d124c07`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/d124c072c9db8d402b353a73646d2d197bae76a4) docs: update README with usage and compatibility matrix
* [`20792f3`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/20792f345b7ff3c8ffa9d65c9ca8dcab1932f49e) feat: generate talosconfig as a secret with proper endpoints
* [`abd206f`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/abd206fd8a98f5478f8ffd0f8686e32be3b7defe) feat: update to CAPI v1.0.x contract (v1beta1)
* [`b7faf9e`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/b7faf9e730b7c9f50ffa94be194ddcf908708a2c) feat: update Talos machinery to 0.13.0
</p>
</details>

### Changes from talos-systems/cluster-api-control-plane-provider-talos
<details><summary>15 commits</summary>
<p>

* [`ea7842f`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/ea7842f6fabefa1775c0a2c7bd43e1a7e0615fe7) release(v0.4.1): prepare release
* [`7f63ad0`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/7f63ad0a391dcf0901edf9009717fb57f503f652) fix: avoid long backoff when trying to bootstrap the cluster
* [`8fc9a6c`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/8fc9a6c3bf7e9de5d074c50c7e7e70d025a5369c) release(v0.4.0): prepare release
* [`b63f1d2`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/b63f1d211002b6ee048de304987151c7deda5db1) release(v0.4.0-beta.2): prepare release
* [`f5f5b2d`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/f5f5b2d441ccfb939e2573ef780a69af371775da) fix: patch the status and use APIReader to get resource
* [`d606d32`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/d606d32a79a1ba567748c782c6ffc4155ce0c81a) release(v0.4.0-beta.1): prepare release
* [`333fc02`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/333fc0221a68e978c5d31399e9befe7d6b29aebe) fix: ensure that bootstrap is called only a single time
* [`77b0bba`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/77b0bba8d7c24026d458acb14c9d2793c0450aa8) test: update templates to v1beta1
* [`a5af5e4`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/a5af5e4b09f8450b73b825d95e49b75b045cb47b) release(v0.4.0-beta.0): prepare release
* [`80b24a0`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/80b24a0abb4a9176a9f1635891a236f299d7dc64) fix: introduce a separate step for release builds
* [`a24dad3`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/a24dad3328d52a3744f685ddde913d66dd17b176) fix: do not allow scaling down controlplane to zero
* [`8a73e6a`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/8a73e6a26e52151b1dd4604c4d0737036f119c30) feat: get rid of init nodes and use bootstrap API to setup cluster
* [`205f4be`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/205f4be2057b3ea81c4dcf47004db6864ff31801) release(v0.4.0-alpha.0): prepare release
* [`b8db449`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/b8db4492d55f910e8a7d2a3b69ab08740963683e) fix: properly pick talos client configuration
* [`61fb582`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/61fb5826391e4434b64619f0590683f7fa7b82b6) feat: support clusterapi v1beta1
</p>
</details>

### Changes from talos-systems/go-blockdevice
<details><summary>2 commits</summary>
<p>

* [`15b182d`](https://github.com/talos-systems/go-blockdevice/commit/15b182db0cd233b163ed83d1724c7e28cf29d71a) fix: return partition table not exist when trying to read an empty dev
* [`b9517d5`](https://github.com/talos-systems/go-blockdevice/commit/b9517d51120d385f97b0026f99ce3c4782940c37) fix: resize partition
</p>
</details>

### Changes from talos-systems/go-smbios
<details><summary>1 commit</summary>
<p>

* [`fd5ec8c`](https://github.com/talos-systems/go-smbios/commit/fd5ec8ce4873790b7fbd46dba9d7f49c9de7176a) fix: remove useless (?) goroutines leading to data race error
</p>
</details>

### Changes from talos-systems/grpc-proxy
<details><summary>44 commits</summary>
<p>

* [`ca3bc61`](https://github.com/talos-systems/grpc-proxy/commit/ca3bc6131f052aa000517339211335aaa4ebb640) fix: ignore some errors so that we don't spam the logs
* [`5c579a7`](https://github.com/talos-systems/grpc-proxy/commit/5c579a7a61475bde3ec9c1efe000d2a55e2a3cb2) feat: allow different formats for messages streaming/unary
* [`6c9f7b3`](https://github.com/talos-systems/grpc-proxy/commit/6c9f7b399173dd5769dbc4e8e366e78f05cead85) fix: allow mode to be set for each request being proxied
* [`cc91c09`](https://github.com/talos-systems/grpc-proxy/commit/cc91c09782824e261bf1c861961a272aedb2b123) refactor: provide better public API, enforce proxying mode
* [`d8d3a75`](https://github.com/talos-systems/grpc-proxy/commit/d8d3a751d1e71d006ba90379eed388c487bbb246) chore: update import paths after repo move
* [`dbf07a4`](https://github.com/talos-systems/grpc-proxy/commit/dbf07a4d9e16fe3cf7407b9921c1746aa24ffaf6) Merge pull request  [#7](https://github.com/talos-systems/grpc-proxy/pull/7) from smira/one2many-4
* [`fc0d27d`](https://github.com/talos-systems/grpc-proxy/commit/fc0d27dc6b5b9db35173f3e78778784a9e7c95bf) More tests, small code fixes, updated README.
* [`d9ce0b1`](https://github.com/talos-systems/grpc-proxy/commit/d9ce0b1053a7f15ea65bf46e94cfe4154493bad7) Merge pull request  [#6](https://github.com/talos-systems/grpc-proxy/pull/6) from smira/one2many-3
* [`2d37ba4`](https://github.com/talos-systems/grpc-proxy/commit/2d37ba444528a00f988671f3a01666e692739a37) Support for one2many streaming calls, tests.
* [`817b035`](https://github.com/talos-systems/grpc-proxy/commit/817b03553ed7d97bd0da09283776d54592d7b5d4) Merge pull request  [#5](https://github.com/talos-systems/grpc-proxy/pull/5) from smira/one2many-2
* [`436b338`](https://github.com/talos-systems/grpc-proxy/commit/436b3383a39fd860f3b2379ffab80a44ae1809f7) More unary one-2-many tests, error propagation.
* [`1f0cb46`](https://github.com/talos-systems/grpc-proxy/commit/1f0cb466268f046e8e9fb78b1902411ac3a753ba) Merge pull request  [#4](https://github.com/talos-systems/grpc-proxy/pull/4) from smira/one2many-1
* [`992a975`](https://github.com/talos-systems/grpc-proxy/commit/992a975ccf0b97e4be329c84bd3018652e8e50ae) Proxying one to many: first iteration
* [`a0988ff`](https://github.com/talos-systems/grpc-proxy/commit/a0988ff2b29839892a7913acd76f26f4e7edcc3a) Merge pull request  [#3](https://github.com/talos-systems/grpc-proxy/pull/3) from smira/small-fixups
* [`e3111ef`](https://github.com/talos-systems/grpc-proxy/commit/e3111ef2c16f0ee4bba597a2ab1ab6a2818c2734) Small fixups in preparation to add one-to-many proxying.
* [`6d76ffc`](https://github.com/talos-systems/grpc-proxy/commit/6d76ffcff89f6636d3689ed1c9b0eebe87722114) Merge pull request  [#2](https://github.com/talos-systems/grpc-proxy/pull/2) from smira/backend-concept
* [`2aad63a`](https://github.com/talos-systems/grpc-proxy/commit/2aad63ac5bae09232ea5ac80b42338e9e3af67c4) Add concept of a 'Backend', but still one to one proxying
* [`7cc4610`](https://github.com/talos-systems/grpc-proxy/commit/7cc46101114a2779d6393e0e8f841bf3febb2753) Merge pull request  [#1](https://github.com/talos-systems/grpc-proxy/pull/1) from smira/build
* [`37f01f3`](https://github.com/talos-systems/grpc-proxy/commit/37f01f3aab3b978a8fecb428fca4d4c722141229) Rework build to use GitHub Actions, linting updates.
* [`0f1106e`](https://github.com/talos-systems/grpc-proxy/commit/0f1106ef9c766333b9acb4b81e705da4bade7215) Move error checking further up (#34)
* [`d5b35f6`](https://github.com/talos-systems/grpc-proxy/commit/d5b35f634383bf8931f8798797daaf9c1a59235e) Update gRPC and fix tests (#27)
* [`67591eb`](https://github.com/talos-systems/grpc-proxy/commit/67591eb23c48346a480470e462289835d96f70da) Break StreamDirector interface, fix metadata propagation for gRPC-Go>1.5. (#20)
* [`97396d9`](https://github.com/talos-systems/grpc-proxy/commit/97396d94749c00db659393ba5123f707062f829f) Merge pull request  [#11](https://github.com/talos-systems/grpc-proxy/pull/11) from mwitkow/fix-close-bug
* [`3fcbd37`](https://github.com/talos-systems/grpc-proxy/commit/3fcbd3737ec6baff505795417e48f162a7a3183c) fixup closing conns
* [`a8f5f87`](https://github.com/talos-systems/grpc-proxy/commit/a8f5f87a2f5e6bc3643b78d64594195b2395a238) fixup tests, extend readme
* [`428fa1c`](https://github.com/talos-systems/grpc-proxy/commit/428fa1c450320041e0ad8e251d6aed435401174e) Fix a channel closing bug
* [`af55d61`](https://github.com/talos-systems/grpc-proxy/commit/af55d612de6c5723a5a59340704db7bc771023ff) Merge pull request  [#10](https://github.com/talos-systems/grpc-proxy/pull/10) from mwitkow/bugfix/streaming-fix
* [`de4d3db`](https://github.com/talos-systems/grpc-proxy/commit/de4d3db538565636e1e977102f6f0bd1ed0ce9c2) remove spurious printfs
* [`84242c4`](https://github.com/talos-systems/grpc-proxy/commit/84242c4e690da18d16d2ab8f2fa47e45986220b6) fix the "i don't know who finished" case
* [`9b22f41`](https://github.com/talos-systems/grpc-proxy/commit/9b22f41d8535fa3e40908c78ae66066c7972b6d9) fix full duplex streaming
* [`c2f7c98`](https://github.com/talos-systems/grpc-proxy/commit/c2f7c98b0b6cd180659aed31e98cbbc18d616b1c) update readme
* [`d654141`](https://github.com/talos-systems/grpc-proxy/commit/d654141edcb92b7fa2bba9d3e690e569c72f8e9d) update README
* [`f457856`](https://github.com/talos-systems/grpc-proxy/commit/f4578565f2d34dc89774128db2bfda3a328cba40) move to proxy subdirectory
* [`4889d78`](https://github.com/talos-systems/grpc-proxy/commit/4889d78e468681601b8229c81807dcf37b00ff63) Add fixup scripts
* [`ef60a37`](https://github.com/talos-systems/grpc-proxy/commit/ef60a37547d137e52873be183f2d7a5626d7c034) version 2 of the grpc-proxy, this time with fewer grpc upstream deps
* [`07aeac1`](https://github.com/talos-systems/grpc-proxy/commit/07aeac13e988c0c0b3a886c79972e20408a765e0) Merge pull request  [#2](https://github.com/talos-systems/grpc-proxy/pull/2) from daniellowtw/master
* [`e5c3df5`](https://github.com/talos-systems/grpc-proxy/commit/e5c3df5b2f0a1ffc4cb755cbe6b30b435e35de37) Fix compatibility with latest grpc library
* [`52be0a5`](https://github.com/talos-systems/grpc-proxy/commit/52be0a559a85f0e2480bde6725f3f144396aa6ef) bugfix: fix gRPC Java deadlock, due to different dispatch logic
* [`822df7d`](https://github.com/talos-systems/grpc-proxy/commit/822df7d86b556b703fc11798a3bdcbaeb60c18a6) Fix reference to mwitkow.
* [`28341d1`](https://github.com/talos-systems/grpc-proxy/commit/28341d171dd4c1a52f46371ddfb5fd2240b79731) move out forward logic to method, allowing for use as `grpc.Server` not found handler.
* [`89e28b4`](https://github.com/talos-systems/grpc-proxy/commit/89e28b42ee9dda8e36522b77e3771d9debc645e0) add reference to upstream grpc bug
* [`00dd588`](https://github.com/talos-systems/grpc-proxy/commit/00dd588ae68adf4187a7fca87db45a73af4c834d) merge upstream `grpc.Server` changes changing the dispatch logic
* [`77edc97`](https://github.com/talos-systems/grpc-proxy/commit/77edc9715de187dcbc9969e2f0e8a04d2087fd13) move to upstream `protobuf` from `gogo`
* [`db71c3e`](https://github.com/talos-systems/grpc-proxy/commit/db71c3e7e812db8d75cb282dac38d953fcb436b3) initial commit, tested and working.
</p>
</details>

### Changes from talos-systems/net
<details><summary>2 commits</summary>
<p>

* [`b4b7181`](https://github.com/talos-systems/net/commit/b4b718179a1aa68e4f54422baf08ca3761723d2d) feat: add a way to filter list of IPs for the machine
* [`0abe5bd`](https://github.com/talos-systems/net/commit/0abe5bdae8f85e4e976bc4d90e95dcb4be8fb853) feat: implement FilterIPs function
</p>
</details>

### Changes from talos-systems/siderolink
<details><summary>7 commits</summary>
<p>

* [`9902ad2`](https://github.com/talos-systems/siderolink/commit/9902ad2774f0655e050233854b9d28dad0431f6c) feat: pass request context and node address to the events sink adapter
* [`d0612a7`](https://github.com/talos-systems/siderolink/commit/d0612a724a1b1336a2bc6a99ed3178e3e40f6d9b) refactor: pass in listener to the log receiver
* [`d86cdd5`](https://github.com/talos-systems/siderolink/commit/d86cdd59ee7a0e0504b739a913991c272c7fb3f5) feat: implement logreceiver for kernel logs
* [`f7cadbc`](https://github.com/talos-systems/siderolink/commit/f7cadbcdfbb84d367e27b5af32e89c138d72d9d7) fix: handle duplicate peer updates
* [`0755b24`](https://github.com/talos-systems/siderolink/commit/0755b24d4682410b251a2a9d662960da15153106) feat: initial implementation of SideroLink
* [`ee73ea9`](https://github.com/talos-systems/siderolink/commit/ee73ea9575a81be7685f24936b2c48a4508a159e) feat: add Talos events sink proto files and the reference implementation
* [`1e2cd9d`](https://github.com/talos-systems/siderolink/commit/1e2cd9d38621234a0a6010e33b1bab264f4d9bdf) Initial commit
</p>
</details>

### Dependency Changes

* **github.com/evanphx/json-patch**                                      v4.11.0 -> v5.6.0
* **github.com/grpc-ecosystem/go-grpc-middleware**                       v1.3.0 **_new_**
* **github.com/onsi/ginkgo**                                             v1.16.4 -> v1.16.5
* **github.com/onsi/gomega**                                             v1.16.0 -> v1.17.0
* **github.com/talos-systems/cluster-api-bootstrap-provider-talos**      v0.4.0 -> v0.5.1
* **github.com/talos-systems/cluster-api-control-plane-provider-talos**  v0.3.0 -> v0.4.1
* **github.com/talos-systems/go-blockdevice**                            v0.2.4 -> v0.2.5
* **github.com/talos-systems/go-smbios**                                 d3a32bea731a -> v0.1.1
* **github.com/talos-systems/grpc-proxy**                                v0.2.0 **_new_**
* **github.com/talos-systems/net**                                       v0.3.0 -> v0.3.1
* **github.com/talos-systems/siderolink**                                9902ad2774f0 **_new_**
* **go.uber.org/zap**                                                    v1.19.0 **_new_**
* **golang.org/x/net**                                                   853a461950ff -> 0a0e4e1bb54c
* **golang.org/x/sys**                                                   39ccf1dd6fa6 -> fe61309f8881
* **golang.zx2c4.com/wireguard/wgctrl**                                  0073765f69ba **_new_**
* **google.golang.org/grpc**                                             v1.41.0 -> v1.42.0
* **inet.af/netaddr**                                                    c74959edd3b6 **_new_**
* **k8s.io/utils**                                                       bdf08cb9a70a -> cb0fa318a74b
* **sigs.k8s.io/cluster-api**                                            v0.4.3 -> v1.0.2
* **sigs.k8s.io/controller-runtime**                                     v0.9.7 -> v0.10.3

Previous release can be found at [v0.4.0](https://github.com/talos-systems/sidero/releases/tag/v0.4.0)

## [Sidero 0.5.0-alpha.0](https://github.com/talos-systems/sidero/releases/tag/v0.5.0-alpha.0) (2021-12-16)

Welcome to the v0.5.0-alpha.0 release of Sidero!  
*This is a pre-release of Sidero*



Please try out the release binaries and report any issues at
https://github.com/talos-systems/sidero/issues.

### IPMI PXE Method

IPMI PXE method (UEFI, BIOS) can now be configured with `SIDERO_CONTROLLER_MANAGER_IPMI_PXE_METHOD` while installing Sidero.


### Siderolink

Sidero now connects to all servers using Siderolink.
This enables streaming of all dmesg logs and events back to sidero.

All server logs can now be viewed by getting logs of one of the container of the `sidero-controller-manager`:

```
kubectl logs -f -n sidero-system deployment/sidero-controller-manager serverlogs
```

Events:

```
kubectl logs -f -n sidero-system deployment/sidero-controller-manager serverevents
```


### Contributors

* Andrey Smirnov
* Michal Witkowski
* Andrey Smirnov
* Artem Chernyshev
* Andrey Smirnov
* Noel Georgi
* Charlie Haley
* Daniel Low
* Jack Wink
* Rubens Farias
* Spencer Smith

### Changes
<details><summary>30 commits</summary>
<p>

* [`f277f91`](https://github.com/talos-systems/sidero/commit/f277f9143d0d701c30be6f0fa39ec76a345f1ba3) release(v0.5.0-alpha.0): prepare release
* [`229dae1`](https://github.com/talos-systems/sidero/commit/229dae1fd6e90e3acca5eee3de05ba05dd39c809) fix: ignore BMC info if username or password are not set
* [`650743a`](https://github.com/talos-systems/sidero/commit/650743ad9a16d326bf9b327d45126a7a8f5e8c21) fix: use environment variables in the ipmitool
* [`ed39a3c`](https://github.com/talos-systems/sidero/commit/ed39a3c4018b4db9c8fa5debad636cc812e2ca55) fix: ensure we setup BMC info *after* disk wiping
* [`025ff04`](https://github.com/talos-systems/sidero/commit/025ff047d0a2d0b239d4ecff801d37d145040bae) feat: additional printer colums
* [`189793e`](https://github.com/talos-systems/sidero/commit/189793e47fb57757b8de39814d3ef831919b4bb2) fix: wait for link up in iPXE script
* [`9a27861`](https://github.com/talos-systems/sidero/commit/9a2786123a598cc3dde6e96c308b73c332f1d70a) fix: make iPXE script replace script on chain request
* [`1bbe3be`](https://github.com/talos-systems/sidero/commit/1bbe3be26fd5a925dec5f8de41ebafd1d852a5f2) feat: extend information printed in the iPXE script, add retries
* [`4cfdeda`](https://github.com/talos-systems/sidero/commit/4cfdedaf97e267e308e4641d1df28fbbbc36d922) feat: provide a way to configure IPMI PXE method
* [`2ff14c4`](https://github.com/talos-systems/sidero/commit/2ff14c4528e31bd5964073b9335791f2d683f021) docs: reflect sidero runs on ARM
* [`274ae33`](https://github.com/talos-systems/sidero/commit/274ae33fc7c3b4b8f5b517914f730a4db3a9840a) fix: check for server power state when in use
* [`d0df929`](https://github.com/talos-systems/sidero/commit/d0df929eb1d1067636caaa2a95c7126be6c02713) feat: implement events manager container
* [`ab12b81`](https://github.com/talos-systems/sidero/commit/ab12b81ef00ad1762aaa251fbaa6b022c298ae62) feat: implement Talos kernel log receiver
* [`5bf7c21`](https://github.com/talos-systems/sidero/commit/5bf7c21f8f61002fd655863cb0ed6ef8f8b649fe) docs: fix clusterctl usage
* [`e77bf54`](https://github.com/talos-systems/sidero/commit/e77bf54a31076dc70dc726f30e924e57e25a14ec) feat: support cordoning server
* [`ab29103`](https://github.com/talos-systems/sidero/commit/ab291034e5a6bafef7eaea86dcb6d594a8afc420) feat: implement SideroLink
* [`adc73b6`](https://github.com/talos-systems/sidero/commit/adc73b67c5ae3b2302be0ced3b87913de4f15b0e) fix: update CAPI resources versions to v1alpha4 in the cluster template
* [`04dbaf0`](https://github.com/talos-systems/sidero/commit/04dbaf065a20e1d5dfbe745d9126b5a174456afd) test: fix Sidero components version in sfyra docs
* [`acb5f57`](https://github.com/talos-systems/sidero/commit/acb5f57f65a0a226d604ad124f189afe5752974a) feat: automatically append talos.config to the Environment
* [`0e7f8a6`](https://github.com/talos-systems/sidero/commit/0e7f8a6820dc77a28d0d264b7e2dd216dab54eb0) docs: metadata server -> sidero
* [`a826281`](https://github.com/talos-systems/sidero/commit/a82628186d84a4f2f49c51b1cb219cf482a3653e) fix: hide password from ipmitool args
* [`ef270df`](https://github.com/talos-systems/sidero/commit/ef270df5b7af6e143bf1449e34a0a577441ab03f) docs: fork docs for 0.5
* [`a0bf382`](https://github.com/talos-systems/sidero/commit/a0bf3828623c3c4ba425601767960f64fc9d85e6) docs: dhcp option-93
* [`bbbe814`](https://github.com/talos-systems/sidero/commit/bbbe814eb98884c525301e49db09b65bc3f0a7b3) chore: conformance check
* [`701d01b`](https://github.com/talos-systems/sidero/commit/701d01ba370c0bd6baa173207575bdaf53c72970) fix: drop into the agent for non-allocated servers
* [`b0e3611`](https://github.com/talos-systems/sidero/commit/b0e3611d2383061533e10d3bc1e642c99e9f70f9) docs: update help message for extra kernel args
* [`bb87567`](https://github.com/talos-systems/sidero/commit/bb87567e242c684977a3e688fc9553f6a77c81e6) chore: bump dependencies
* [`62ab9be`](https://github.com/talos-systems/sidero/commit/62ab9beacfa33e8fb348b52582e63d6234af9146) fix: update extension for controlplane.yam l talosctl generates YAML files with the .yaml extension, update to the apply-config command to reflect this
* [`0511d78`](https://github.com/talos-systems/sidero/commit/0511d78ef29651bce4d1b0ce8d24592582cfdb2e) feat: allow defining http server and api server ports separately
* [`432ca2a`](https://github.com/talos-systems/sidero/commit/432ca2a467a7fb4ee3cc448b7508872aa2674654) docs: create `v0.4` docs and set them as latest
</p>
</details>

### Changes from talos-systems/cluster-api-control-plane-provider-talos
<details><summary>2 commits</summary>
<p>

* [`ebb7340`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/ebb73403bce3d7d7f1dc9667cede588c0cdfda6f) release(v0.3.1): prepare release
* [`8d99bfd`](https://github.com/talos-systems/cluster-api-control-plane-provider-talos/commit/8d99bfdb4af877e149e3eb609217620fea0da47c) fix: properly pick talos client configuration
</p>
</details>

### Changes from talos-systems/go-smbios
<details><summary>1 commit</summary>
<p>

* [`fd5ec8c`](https://github.com/talos-systems/go-smbios/commit/fd5ec8ce4873790b7fbd46dba9d7f49c9de7176a) fix: remove useless (?) goroutines leading to data race error
</p>
</details>

### Changes from talos-systems/grpc-proxy
<details><summary>44 commits</summary>
<p>

* [`ca3bc61`](https://github.com/talos-systems/grpc-proxy/commit/ca3bc6131f052aa000517339211335aaa4ebb640) fix: ignore some errors so that we don't spam the logs
* [`5c579a7`](https://github.com/talos-systems/grpc-proxy/commit/5c579a7a61475bde3ec9c1efe000d2a55e2a3cb2) feat: allow different formats for messages streaming/unary
* [`6c9f7b3`](https://github.com/talos-systems/grpc-proxy/commit/6c9f7b399173dd5769dbc4e8e366e78f05cead85) fix: allow mode to be set for each request being proxied
* [`cc91c09`](https://github.com/talos-systems/grpc-proxy/commit/cc91c09782824e261bf1c861961a272aedb2b123) refactor: provide better public API, enforce proxying mode
* [`d8d3a75`](https://github.com/talos-systems/grpc-proxy/commit/d8d3a751d1e71d006ba90379eed388c487bbb246) chore: update import paths after repo move
* [`dbf07a4`](https://github.com/talos-systems/grpc-proxy/commit/dbf07a4d9e16fe3cf7407b9921c1746aa24ffaf6) Merge pull request  [#7](https://github.com/talos-systems/grpc-proxy/pull/7) from smira/one2many-4
* [`fc0d27d`](https://github.com/talos-systems/grpc-proxy/commit/fc0d27dc6b5b9db35173f3e78778784a9e7c95bf) More tests, small code fixes, updated README.
* [`d9ce0b1`](https://github.com/talos-systems/grpc-proxy/commit/d9ce0b1053a7f15ea65bf46e94cfe4154493bad7) Merge pull request  [#6](https://github.com/talos-systems/grpc-proxy/pull/6) from smira/one2many-3
* [`2d37ba4`](https://github.com/talos-systems/grpc-proxy/commit/2d37ba444528a00f988671f3a01666e692739a37) Support for one2many streaming calls, tests.
* [`817b035`](https://github.com/talos-systems/grpc-proxy/commit/817b03553ed7d97bd0da09283776d54592d7b5d4) Merge pull request  [#5](https://github.com/talos-systems/grpc-proxy/pull/5) from smira/one2many-2
* [`436b338`](https://github.com/talos-systems/grpc-proxy/commit/436b3383a39fd860f3b2379ffab80a44ae1809f7) More unary one-2-many tests, error propagation.
* [`1f0cb46`](https://github.com/talos-systems/grpc-proxy/commit/1f0cb466268f046e8e9fb78b1902411ac3a753ba) Merge pull request  [#4](https://github.com/talos-systems/grpc-proxy/pull/4) from smira/one2many-1
* [`992a975`](https://github.com/talos-systems/grpc-proxy/commit/992a975ccf0b97e4be329c84bd3018652e8e50ae) Proxying one to many: first iteration
* [`a0988ff`](https://github.com/talos-systems/grpc-proxy/commit/a0988ff2b29839892a7913acd76f26f4e7edcc3a) Merge pull request  [#3](https://github.com/talos-systems/grpc-proxy/pull/3) from smira/small-fixups
* [`e3111ef`](https://github.com/talos-systems/grpc-proxy/commit/e3111ef2c16f0ee4bba597a2ab1ab6a2818c2734) Small fixups in preparation to add one-to-many proxying.
* [`6d76ffc`](https://github.com/talos-systems/grpc-proxy/commit/6d76ffcff89f6636d3689ed1c9b0eebe87722114) Merge pull request  [#2](https://github.com/talos-systems/grpc-proxy/pull/2) from smira/backend-concept
* [`2aad63a`](https://github.com/talos-systems/grpc-proxy/commit/2aad63ac5bae09232ea5ac80b42338e9e3af67c4) Add concept of a 'Backend', but still one to one proxying
* [`7cc4610`](https://github.com/talos-systems/grpc-proxy/commit/7cc46101114a2779d6393e0e8f841bf3febb2753) Merge pull request  [#1](https://github.com/talos-systems/grpc-proxy/pull/1) from smira/build
* [`37f01f3`](https://github.com/talos-systems/grpc-proxy/commit/37f01f3aab3b978a8fecb428fca4d4c722141229) Rework build to use GitHub Actions, linting updates.
* [`0f1106e`](https://github.com/talos-systems/grpc-proxy/commit/0f1106ef9c766333b9acb4b81e705da4bade7215) Move error checking further up (#34)
* [`d5b35f6`](https://github.com/talos-systems/grpc-proxy/commit/d5b35f634383bf8931f8798797daaf9c1a59235e) Update gRPC and fix tests (#27)
* [`67591eb`](https://github.com/talos-systems/grpc-proxy/commit/67591eb23c48346a480470e462289835d96f70da) Break StreamDirector interface, fix metadata propagation for gRPC-Go>1.5. (#20)
* [`97396d9`](https://github.com/talos-systems/grpc-proxy/commit/97396d94749c00db659393ba5123f707062f829f) Merge pull request  [#11](https://github.com/talos-systems/grpc-proxy/pull/11) from mwitkow/fix-close-bug
* [`3fcbd37`](https://github.com/talos-systems/grpc-proxy/commit/3fcbd3737ec6baff505795417e48f162a7a3183c) fixup closing conns
* [`a8f5f87`](https://github.com/talos-systems/grpc-proxy/commit/a8f5f87a2f5e6bc3643b78d64594195b2395a238) fixup tests, extend readme
* [`428fa1c`](https://github.com/talos-systems/grpc-proxy/commit/428fa1c450320041e0ad8e251d6aed435401174e) Fix a channel closing bug
* [`af55d61`](https://github.com/talos-systems/grpc-proxy/commit/af55d612de6c5723a5a59340704db7bc771023ff) Merge pull request  [#10](https://github.com/talos-systems/grpc-proxy/pull/10) from mwitkow/bugfix/streaming-fix
* [`de4d3db`](https://github.com/talos-systems/grpc-proxy/commit/de4d3db538565636e1e977102f6f0bd1ed0ce9c2) remove spurious printfs
* [`84242c4`](https://github.com/talos-systems/grpc-proxy/commit/84242c4e690da18d16d2ab8f2fa47e45986220b6) fix the "i don't know who finished" case
* [`9b22f41`](https://github.com/talos-systems/grpc-proxy/commit/9b22f41d8535fa3e40908c78ae66066c7972b6d9) fix full duplex streaming
* [`c2f7c98`](https://github.com/talos-systems/grpc-proxy/commit/c2f7c98b0b6cd180659aed31e98cbbc18d616b1c) update readme
* [`d654141`](https://github.com/talos-systems/grpc-proxy/commit/d654141edcb92b7fa2bba9d3e690e569c72f8e9d) update README
* [`f457856`](https://github.com/talos-systems/grpc-proxy/commit/f4578565f2d34dc89774128db2bfda3a328cba40) move to proxy subdirectory
* [`4889d78`](https://github.com/talos-systems/grpc-proxy/commit/4889d78e468681601b8229c81807dcf37b00ff63) Add fixup scripts
* [`ef60a37`](https://github.com/talos-systems/grpc-proxy/commit/ef60a37547d137e52873be183f2d7a5626d7c034) version 2 of the grpc-proxy, this time with fewer grpc upstream deps
* [`07aeac1`](https://github.com/talos-systems/grpc-proxy/commit/07aeac13e988c0c0b3a886c79972e20408a765e0) Merge pull request  [#2](https://github.com/talos-systems/grpc-proxy/pull/2) from daniellowtw/master
* [`e5c3df5`](https://github.com/talos-systems/grpc-proxy/commit/e5c3df5b2f0a1ffc4cb755cbe6b30b435e35de37) Fix compatibility with latest grpc library
* [`52be0a5`](https://github.com/talos-systems/grpc-proxy/commit/52be0a559a85f0e2480bde6725f3f144396aa6ef) bugfix: fix gRPC Java deadlock, due to different dispatch logic
* [`822df7d`](https://github.com/talos-systems/grpc-proxy/commit/822df7d86b556b703fc11798a3bdcbaeb60c18a6) Fix reference to mwitkow.
* [`28341d1`](https://github.com/talos-systems/grpc-proxy/commit/28341d171dd4c1a52f46371ddfb5fd2240b79731) move out forward logic to method, allowing for use as `grpc.Server` not found handler.
* [`89e28b4`](https://github.com/talos-systems/grpc-proxy/commit/89e28b42ee9dda8e36522b77e3771d9debc645e0) add reference to upstream grpc bug
* [`00dd588`](https://github.com/talos-systems/grpc-proxy/commit/00dd588ae68adf4187a7fca87db45a73af4c834d) merge upstream `grpc.Server` changes changing the dispatch logic
* [`77edc97`](https://github.com/talos-systems/grpc-proxy/commit/77edc9715de187dcbc9969e2f0e8a04d2087fd13) move to upstream `protobuf` from `gogo`
* [`db71c3e`](https://github.com/talos-systems/grpc-proxy/commit/db71c3e7e812db8d75cb282dac38d953fcb436b3) initial commit, tested and working.
</p>
</details>

### Changes from talos-systems/net
<details><summary>2 commits</summary>
<p>

* [`b4b7181`](https://github.com/talos-systems/net/commit/b4b718179a1aa68e4f54422baf08ca3761723d2d) feat: add a way to filter list of IPs for the machine
* [`0abe5bd`](https://github.com/talos-systems/net/commit/0abe5bdae8f85e4e976bc4d90e95dcb4be8fb853) feat: implement FilterIPs function
</p>
</details>

### Changes from talos-systems/siderolink
<details><summary>7 commits</summary>
<p>

* [`9902ad2`](https://github.com/talos-systems/siderolink/commit/9902ad2774f0655e050233854b9d28dad0431f6c) feat: pass request context and node address to the events sink adapter
* [`d0612a7`](https://github.com/talos-systems/siderolink/commit/d0612a724a1b1336a2bc6a99ed3178e3e40f6d9b) refactor: pass in listener to the log receiver
* [`d86cdd5`](https://github.com/talos-systems/siderolink/commit/d86cdd59ee7a0e0504b739a913991c272c7fb3f5) feat: implement logreceiver for kernel logs
* [`f7cadbc`](https://github.com/talos-systems/siderolink/commit/f7cadbcdfbb84d367e27b5af32e89c138d72d9d7) fix: handle duplicate peer updates
* [`0755b24`](https://github.com/talos-systems/siderolink/commit/0755b24d4682410b251a2a9d662960da15153106) feat: initial implementation of SideroLink
* [`ee73ea9`](https://github.com/talos-systems/siderolink/commit/ee73ea9575a81be7685f24936b2c48a4508a159e) feat: add Talos events sink proto files and the reference implementation
* [`1e2cd9d`](https://github.com/talos-systems/siderolink/commit/1e2cd9d38621234a0a6010e33b1bab264f4d9bdf) Initial commit
</p>
</details>

### Dependency Changes

* **github.com/evanphx/json-patch**                                      v4.11.0 -> v5.6.0
* **github.com/grpc-ecosystem/go-grpc-middleware**                       v1.3.0 **_new_**
* **github.com/onsi/ginkgo**                                             v1.16.4 -> v1.16.5
* **github.com/onsi/gomega**                                             v1.16.0 -> v1.17.0
* **github.com/talos-systems/cluster-api-control-plane-provider-talos**  v0.3.0 -> v0.3.1
* **github.com/talos-systems/go-smbios**                                 d3a32bea731a -> v0.1.1
* **github.com/talos-systems/grpc-proxy**                                v0.2.0 **_new_**
* **github.com/talos-systems/net**                                       v0.3.0 -> v0.3.1
* **github.com/talos-systems/siderolink**                                9902ad2774f0 **_new_**
* **go.uber.org/zap**                                                    v1.19.0 **_new_**
* **golang.org/x/net**                                                   853a461950ff -> 6635138e15ea
* **golang.org/x/sys**                                                   39ccf1dd6fa6 -> 51b60fd695b3
* **golang.zx2c4.com/wireguard/wgctrl**                                  0073765f69ba **_new_**
* **google.golang.org/grpc**                                             v1.41.0 -> v1.42.0
* **inet.af/netaddr**                                                    c74959edd3b6 **_new_**

Previous release can be found at [v0.4.0](https://github.com/talos-systems/sidero/releases/tag/v0.4.0)

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

