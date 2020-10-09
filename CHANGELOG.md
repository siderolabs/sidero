
<a name="v0.1.0-alpha.3"></a>
## [v0.1.0-alpha.3](https://github.com/talos-systems/talos/compare/v0.1.0-alpha.2...v0.1.0-alpha.3) (2020-10-08)

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

