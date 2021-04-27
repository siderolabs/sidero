REGISTRY ?= ghcr.io
USERNAME ?= talos-systems
SHA ?= $(shell git describe --match=none --always --abbrev=8 --dirty)
TAG ?= $(shell git describe --tag --always --dirty)
BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD)
REGISTRY_AND_USERNAME := $(REGISTRY)/$(USERNAME)
IMAGE := $(REGISTRY_AND_USERNAME)/$(NAME)
MODULE := $(shell head -1 go.mod | cut -d' ' -f2)

ARTIFACTS := _out
TEST_PKGS ?= ./...
TALOS_RELEASE ?= v0.9.3

TOOLS ?= ghcr.io/talos-systems/tools:v0.5.0
PKGS ?= v0.5.0

SFYRA_CLUSTERCTL_CONFIG ?= $(HOME)/.cluster-api/clusterctl.sfyra.yaml

BUILD := docker buildx build
PLATFORM ?= linux/amd64
PROGRESS ?= auto
PUSH ?= false
COMMON_ARGS := --file=Dockerfile
COMMON_ARGS += --progress=$(PROGRESS)
COMMON_ARGS += --platform=$(PLATFORM)
COMMON_ARGS += --build-arg=REGISTRY_AND_USERNAME=$(REGISTRY_AND_USERNAME)
COMMON_ARGS += --build-arg=TAG=$(TAG)
COMMON_ARGS += --build-arg=MODULE=$(MODULE)
COMMON_ARGS += --build-arg=TEST_PKGS=$(TEST_PKGS)
COMMON_ARGS += --build-arg=PKGS=$(PKGS)
COMMON_ARGS += --build-arg=TOOLS=$(TOOLS)
COMMON_ARGS += --build-arg=TALOS_RELEASE=$(TALOS_RELEASE)

all: manifests generate cluster-api-provider-sidero metal-controller-manager metal-metadata-server sfyra


# Help Menu

define HELP_MENU_HEADER
# Getting Started

To build this project, you must have the following installed:

- git
- make
- docker (19.03 or higher)
- buildx (https://github.com/docker/buildx)

## Creating a Builder Instance

The build process makes use of features not currently supported by the default
builder instance (docker driver). To create a compatible builder instance, run:

```
docker buildx create --driver docker-container --name local --buildkitd-flags --use
```

If you already have a compatible builder instance, you may use that instead.

## Artifacts

All artifacts will be output to ./$(ARTIFACTS). Images will be tagged with the
registry "$(REGISTRY)", username "$(USERNAME)", and a dynamic tag (e.g. $(IMAGE):$(TAG)).
The registry and username can be overriden by exporting REGISTRY, and USERNAME
respectively.

endef

export HELP_MENU_HEADER

help: ## This help menu.
	@echo "$$HELP_MENU_HEADER"
	@grep -E '^[a-zA-Z%_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

# Build Abstractions

target-%: ## Builds the specified target defined in the Dockerfile. The build result will remain only in the build cache.
	@$(BUILD) \
		--target=$* \
		$(COMMON_ARGS) \
		$(TARGET_ARGS) .

local-%: ## Builds the specified target defined in the Dockerfile using the local output type. The build result will be output to the specified local destination.
	@$(MAKE) target-$* TARGET_ARGS="--output=type=local,dest=$(DEST) $(TARGET_ARGS)"

docker-%: ## Builds the specified target defined in the Dockerfile using the docker output type. The build result will be loaded into docker.
	@$(MAKE) target-$* TARGET_ARGS="--tag $(IMAGE):$(TAG) $(TARGET_ARGS)"

# Code Generation

.PHONY: generate
generate: ## Generate source code.
	@$(MAKE) local-$@ DEST=./ PLATFORM=linux/amd64

.PHONY: manifests
manifests: ## Generate manifests (e.g. CRD, RBAC, etc.).
	@$(MAKE) local-$@ DEST=./ PLATFORM=linux/amd64

# Artifacts

.PHONY: release
release: manifests ## Create the release YAML. The build result will be ouput to the specified local destination.
	@$(MAKE) local-$@ DEST=./$(ARTIFACTS)

.PHONY: cluster-api-provider-sidero
cluster-api-provider-sidero: ## Build the CAPI provider container image.
	@$(MAKE) docker-$@ TARGET_ARGS="--push=$(PUSH)" NAME="$@"

.PHONY: metal-controller-manager
metal-controller-manager: ## Build the CAPI provider container image.
	@$(MAKE) docker-$@ TARGET_ARGS="--push=$(PUSH)" NAME="$@"

.PHONY: metal-metadata-server
metal-metadata-server: ## Build the CAPI provider container image.
	@$(MAKE) docker-$@ TARGET_ARGS="--push=$(PUSH)" NAME="$@"

# Sfyra

$(ARTIFACTS)/$(TALOS_RELEASE)/%:
	@mkdir -p $(ARTIFACTS)/$(TALOS_RELEASE)/
	@curl -L -o "$(ARTIFACTS)/$(TALOS_RELEASE)/$*" "https://github.com/talos-systems/talos/releases/download/$(TALOS_RELEASE)/$*"

.PHONY: $(ARTIFACTS)/$(TALOS_RELEASE)
$(ARTIFACTS)/$(TALOS_RELEASE): $(ARTIFACTS)/$(TALOS_RELEASE)/vmlinuz-amd64 $(ARTIFACTS)/$(TALOS_RELEASE)/initramfs-amd64.xz $(ARTIFACTS)/$(TALOS_RELEASE)/talosctl-linux-amd64

.PHONY: talos-artifacts
talos-artifacts: $(ARTIFACTS)/$(TALOS_RELEASE)
	@chmod +x $(ARTIFACTS)/$(TALOS_RELEASE)/talosctl-linux-amd64

.PHONY: sfyra
sfyra: ## Build the Sfyra test binary.
	@$(MAKE) local-$@ DEST=./$(ARTIFACTS) PLATFORM=linux/amd64

.PHONY: clusterctl-release
clusterctl-release: release
	@COMPONENTS_YAML="$(abspath $(ARTIFACTS)/infrastructure-sidero/$(TAG)/infrastructure-components.yaml)" \
		CLUSTERCTL_CONFIG=$(SFYRA_CLUSTERCTL_CONFIG) \
		./hack/scripts/generate-clusterctl-config.sh

.PHONY: run-sfyra
run-sfyra: talos-artifacts clusterctl-release
	@ARTIFACTS=$(ARTIFACTS) \
		CLUSTERCTL_CONFIG=$(SFYRA_CLUSTERCTL_CONFIG) \
		TALOS_RELEASE=$(TALOS_RELEASE) \
		./hack/scripts/integration-test.sh

# Development

.PHONY: deploy
deploy: manifests ## Deploy to a cluster. This is for testing purposes only.
	kubectl apply -k config/default

.PHONY: destroy
destroy: ## Remove from a cluster. This is for testing purposes only.
	kubectl delete -k config/default

.PHONY: install
install: manifests ## Install CRDs into a cluster.
	kubectl apply -k config/crd

.PHONY: uninstall
uninstall: manifests ## Uninstall CRDs from a cluster.
	kubectl delete -k config/crd

.PHONY: run
run: install ## Run the controller locally. This is for testing purposes only.
	@$(MAKE) docker-container TARGET_ARGS="--load"
	@docker run --rm -it --net host -v $(PWD):/src -v $(KUBECONFIG):/root/.kube/config -e KUBECONFIG=/root/.kube/config $(IMAGE):$(TAG)

# Code Quality

.PHONY: fmt
fmt: ## Formats the source code.
	@$(MAKE) local-fmt DEST=./ PLATFORM=linux/amd64

lint-%: ## Runs the specified linter. Valid options are go, protobuf, and markdown (e.g. lint-go).
	@$(MAKE) target-lint-$* PLATFORM=linux/amd64

lint: ## Runs linters on go, protobuf, and markdown file types.
	@$(MAKE) lint-go lint-markdown lint-sfyra

# Tests

.PHONY: unit-tests
unit-tests: ## Performs unit tests.
	@$(MAKE) local-$@ DEST=$(ARTIFACTS) PLATFORM=linux/amd64

.PHONY: unit-tests-race
unit-tests-race: ## Performs unit tests with race detection enabled.
	@$(MAKE) target-$@ PLATFORM=linux/amd64

# Utilities

.PHONY: clean
clean:
	@rm -rf $(ARTIFACTS)
