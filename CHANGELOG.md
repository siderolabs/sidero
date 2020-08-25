
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

