REGISTRY ?= ghcr.io
USERNAME ?= siderolabs
SHA ?= $(shell git describe --match=none --always --abbrev=8 --dirty)
TAG ?= $(shell git describe --tag --always --dirty)
BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD)
REGISTRY_AND_USERNAME := $(REGISTRY)/$(USERNAME)
IMAGE := $(REGISTRY_AND_USERNAME)/$(NAME)
MODULE := $(shell head -1 go.mod | cut -d' ' -f2)

ARTIFACTS := _out
TEST_PKGS ?= ./...
TALOS_RELEASE ?= v1.11.0-beta.2
DEFAULT_K8S_VERSION ?= v1.33.3

KRES_IMAGE ?= ghcr.io/siderolabs/kres:latest

TOOLS ?= ghcr.io/siderolabs/tools:v1.11.0
PKGS ?= v1.11.0

SFYRA_CLUSTERCTL_CONFIG ?= $(HOME)/.cluster-api/clusterctl.sfyra.yaml

CGO_ENABLED ?= 0
GO_BUILDFLAGS ?=
GO_LDFLAGS ?=

WITH_RACE ?= false
WITH_DEBUG ?= false

ifneq (, $(filter $(WITH_RACE), t true TRUE y yes 1))
CGO_ENABLED = 1
GO_BUILDFLAGS += -race
GO_LDFLAGS += -linkmode=external -extldflags '-static'
endif

ifneq (, $(filter $(WITH_DEBUG), t true TRUE y yes 1))
GO_BUILDFLAGS += -tags sidero.debug
else
GO_LDFLAGS += -s -w
endif

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
COMMON_ARGS += --build-arg=DEFAULT_K8S_VERSION=$(DEFAULT_K8S_VERSION)
COMMON_ARGS += --build-arg=CGO_ENABLED=$(CGO_ENABLED)
COMMON_ARGS += --build-arg=GO_BUILDFLAGS="$(GO_BUILDFLAGS)"
COMMON_ARGS += --build-arg=GO_LDFLAGS="$(GO_LDFLAGS)"

all: manifests generate caps-controller-manager sidero-controller-manager sfyra


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
docker buildx create --driver docker-container --name local --use
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
	@grep -E '^[a-zA-Z0-9%_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

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

.PHONY: check-dirty
check-dirty: ## Verifies that source tree is not dirty
	@if test -n "`git status --porcelain`"; then echo "Source tree is dirty"; git status; exit 1 ; fi

# Artifacts

.PHONY: release
release: manifests ## Create the release YAML. The build result will be output to the specified local destination.
	@$(MAKE) local-$@ DEST=./$(ARTIFACTS)

.PHONY: caps-controller-manager
caps-controller-manager: ## Build the CAPI provider container image.
	@$(MAKE) docker-$@ TARGET_ARGS="--push=$(PUSH)" NAME="$@"

.PHONY: sidero-controller-manager
sidero-controller-manager: ## Build the CAPI provider container image.
	@$(MAKE) docker-$@ TARGET_ARGS="--push=$(PUSH)" NAME="$@"

.PHONY: release-notes
release-notes:
	@mkdir -p $(ARTIFACTS)
	@ARTIFACTS=$(ARTIFACTS) ./hack/release.sh $@ $(ARTIFACTS)/RELEASE_NOTES.md $(TAG)

# Sfyra

$(ARTIFACTS)/$(TALOS_RELEASE)/%:
	@mkdir -p $(ARTIFACTS)/$(TALOS_RELEASE)/
	@curl -L -o "$(ARTIFACTS)/$(TALOS_RELEASE)/$*" "https://github.com/siderolabs/talos/releases/download/$(TALOS_RELEASE)/$*"

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
run-sfyra: talos-artifacts clusterctl-release ## Run Sfyra integration test.
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

.PHONY: conformance
conformance: ## Performs policy checks against the commit and source code.
	docker run --rm -it -v $(PWD):/src -w /src ghcr.io/siderolabs/conform:v0.1.0-alpha.25-2-g625a1c5 enforce

.PHONY: clean
clean:
	@rm -rf $(ARTIFACTS)

.PHONY: docs-preview
docs-preview: ## Starts a local preview of the documentation using Hugo in docker
	@docker run --rm --interactive --tty \
        --volume $(PWD):/src --workdir /src/website \
        --publish 1313:1313 \
        klakegg/hugo:0.95.0-ext-alpine \
        server

.PHONY: rekres
rekres:
	@docker pull $(KRES_IMAGE)
	@docker run --rm --net=host --user $(shell id -u):$(shell id -g) -v $(PWD):/src -w /src -e GITHUB_TOKEN $(KRES_IMAGE)
