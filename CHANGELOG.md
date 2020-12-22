
<a name="v0.1.0-beta.0"></a>
## [v0.1.0-beta.0](https://github.com/talos-systems/talos/compare/v0.1.0-alpha.12...v0.1.0-beta.0) (2020-12-21)

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

