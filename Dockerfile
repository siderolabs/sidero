# syntax = docker/dockerfile-upstream:1.2.0-labs

# Meta args applied to stage base names.

ARG TOOLS
ARG PKGS

# Resolve package images using ${PKGS} to be used later in COPY --from=.

FROM ghcr.io/talos-systems/ca-certificates:${PKGS} AS pkg-ca-certificates
FROM ghcr.io/talos-systems/fhs:${PKGS} AS pkg-fhs
FROM ghcr.io/talos-systems/ipmitool:${PKGS} AS pkg-ipmitool
FROM --platform=amd64 ghcr.io/talos-systems/ipmitool:${PKGS} AS pkg-ipmitool-amd64
FROM --platform=arm64 ghcr.io/talos-systems/ipmitool:${PKGS} AS pkg-ipmitool-arm64
FROM ghcr.io/talos-systems/libressl:${PKGS} AS pkg-libressl
FROM --platform=amd64 ghcr.io/talos-systems/libressl:${PKGS} AS pkg-libressl-amd64
FROM --platform=arm64 ghcr.io/talos-systems/libressl:${PKGS} AS pkg-libressl-arm64
FROM --platform=amd64 ghcr.io/talos-systems/linux-firmware:${PKGS} AS pkg-linux-firmware-amd64
FROM --platform=arm64 ghcr.io/talos-systems/linux-firmware:${PKGS} AS pkg-linux-firmware-arm64
FROM ghcr.io/talos-systems/musl:${PKGS} AS pkg-musl
FROM --platform=amd64 ghcr.io/talos-systems/musl:${PKGS} AS pkg-musl-amd64
FROM --platform=arm64 ghcr.io/talos-systems/musl:${PKGS} AS pkg-musl-arm64
FROM --platform=amd64 ghcr.io/talos-systems/kernel:${PKGS} AS pkg-kernel-amd64
FROM --platform=arm64 ghcr.io/talos-systems/kernel:${PKGS} AS pkg-kernel-arm64
FROM ghcr.io/talos-systems/liblzma:${PKGS} AS pkg-liblzma
FROM ghcr.io/talos-systems/ipxe:${PKGS} AS pkg-ipxe
FROM --platform=amd64 ghcr.io/talos-systems/ipxe:${PKGS} AS pkg-ipxe-amd64
FROM --platform=arm64 ghcr.io/talos-systems/ipxe:${PKGS} AS pkg-ipxe-arm64

# The base target provides the base for running various tasks against the source
# code

FROM --platform=${BUILDPLATFORM} ${TOOLS} AS base
SHELL ["/toolchain/bin/bash", "-c"]
ENV PATH /toolchain/bin:/toolchain/go/bin:/go/bin
RUN ["/toolchain/bin/mkdir", "/bin", "/tmp"]
RUN ["/toolchain/bin/ln", "-svf", "/toolchain/bin/bash", "/bin/sh"]
RUN ["/toolchain/bin/ln", "-svf", "/toolchain/etc/ssl", "/etc/ssl"]
ENV GO111MODULE on
ENV GOPROXY https://proxy.golang.org
ARG CGO_ENABLED
ENV CGO_ENABLED ${CGO_ENABLED}
ENV GOCACHE /.cache/go-build
ENV GOMODCACHE /.cache/mod
RUN --mount=type=cache,target=/.cache go install sigs.k8s.io/controller-tools/cmd/controller-gen@v0.5.0
RUN --mount=type=cache,target=/.cache go install k8s.io/code-generator/cmd/conversion-gen@v0.21.0
RUN --mount=type=cache,target=/.cache go install mvdan.cc/gofumpt/gofumports@v0.1.1
RUN curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | bash -s -- -b /toolchain/bin v1.38.0
WORKDIR /src
COPY ./go.mod ./
COPY ./go.sum ./
RUN --mount=type=cache,target=/.cache go mod download
RUN --mount=type=cache,target=/.cache go mod verify
COPY ./app/ ./app/
COPY ./internal/ ./internal/
COPY ./hack/ ./hack/
RUN --mount=type=cache,target=/.cache go list -mod=readonly all >/dev/null
RUN --mount=type=cache,target=/.cache ! go mod tidy -v 2>&1 | grep .

FROM base AS manifests-build
RUN --mount=type=cache,target=/.cache controller-gen \
  crd:crdVersions=v1 paths="./app/caps-controller-manager/api/..." output:crd:dir="./app/caps-controller-manager/config/crd/bases" \
  rbac:roleName=manager-role paths="./app/caps-controller-manager/controllers/..." output:rbac:dir="./app/caps-controller-manager/config/rbac" \
  webhook output:webhook:dir="./app/caps-controller-manager/config/webhook"
RUN --mount=type=cache,target=/.cache controller-gen \
  crd:crdVersions=v1 paths="./app/sidero-controller-manager/api/..." output:crd:dir="./app/sidero-controller-manager/config/crd/bases" \
  rbac:roleName=manager-role paths="./app/sidero-controller-manager/controllers/..." output:rbac:dir="./app/sidero-controller-manager/config/rbac" \
  webhook output:webhook:dir="./app/sidero-controller-manager/config/webhook"

FROM scratch AS manifests
COPY --from=manifests-build /src/app/caps-controller-manager/config ./app/caps-controller-manager/config
COPY --from=manifests-build /src/app/sidero-controller-manager/config ./app/sidero-controller-manager/config

FROM base AS generate-build
COPY ./app/sidero-controller-manager/internal/api/api.proto \
  /src/app/sidero-controller-manager/internal/api/api.proto
RUN protoc -I/src/app/sidero-controller-manager/internal/api \
  --go_out=paths=source_relative:/src/app/sidero-controller-manager/internal/api --go-grpc_out=paths=source_relative:/src/app/sidero-controller-manager/internal/api \
  api.proto
RUN --mount=type=cache,target=/.cache controller-gen object:headerFile="./hack/boilerplate.go.txt" paths="./..."
RUN --mount=type=cache,target=/.cache conversion-gen --input-dirs="./app/caps-controller-manager/api/v1alpha2" --output-base ./ --output-file-base="zz_generated.conversion" --go-header-file="./hack/boilerplate.go.txt"
ARG MODULE
RUN --mount=type=cache,target=/.cache gofumports -w -local ${MODULE} .

FROM scratch AS generate
COPY --from=generate-build /src/app/caps-controller-manager/api ./app/caps-controller-manager/api
COPY --from=generate-build /src/app/sidero-controller-manager/api ./app/sidero-controller-manager/api
COPY --from=generate-build /src/app/sidero-controller-manager/internal/api ./app/sidero-controller-manager/internal/api

FROM --platform=${BUILDPLATFORM} alpine:3.14 AS release-build
ADD https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize%2Fv4.1.0/kustomize_v4.1.0_linux_amd64.tar.gz .
RUN  tar -xf kustomize_v4.1.0_linux_amd64.tar.gz -C /usr/local/bin && rm kustomize_v4.1.0_linux_amd64.tar.gz
COPY ./config ./config
COPY ./templates ./templates
COPY ./app/caps-controller-manager/config ./app/caps-controller-manager/config
COPY ./app/sidero-controller-manager/config ./app/sidero-controller-manager/config
ARG REGISTRY_AND_USERNAME
ARG TAG
RUN cd ./app/caps-controller-manager/config/manager \
  && kustomize edit set image controller=${REGISTRY_AND_USERNAME}/caps-controller-manager:${TAG}
RUN cd ./app/sidero-controller-manager/config/manager \
  && kustomize edit set image controller=${REGISTRY_AND_USERNAME}/sidero-controller-manager:${TAG}
RUN kustomize build config > /infrastructure-components.yaml \
  && cp ./config/metadata/metadata.yaml /metadata.yaml \
  && cp ./templates/cluster-template.yaml /cluster-template.yaml

FROM scratch AS release
ARG TAG
COPY --from=release-build /infrastructure-components.yaml /infrastructure-sidero/${TAG}/infrastructure-components.yaml
COPY --from=release-build /metadata.yaml /infrastructure-sidero/${TAG}/metadata.yaml
COPY --from=release-build /cluster-template.yaml /infrastructure-sidero/${TAG}/cluster-template.yaml

FROM base AS build-caps-controller-manager
ARG TARGETARCH
ARG GO_BUILDFLAGS
ARG GO_LDFLAGS
RUN --mount=type=cache,target=/.cache GOOS=linux GOARCH=${TARGETARCH} go build ${GO_BUILDFLAGS} -ldflags "${GO_LDFLAGS}" -o /manager ./app/caps-controller-manager
RUN chmod +x /manager

## TODO(rsmitty): make bmc pkg and move to talos-systems image
FROM scratch AS caps-controller-manager
COPY --from=pkg-ca-certificates / /
COPY --from=pkg-fhs / /
COPY --from=pkg-musl / /
COPY --from=pkg-libressl / /
COPY --from=pkg-ipmitool / /
COPY --from=build-caps-controller-manager /manager /manager
LABEL org.opencontainers.image.source https://github.com/talos-systems/sidero
ENTRYPOINT [ "/manager" ]

FROM base AS build-sidero-controller-manager
ARG TALOS_RELEASE
ARG TARGETARCH
ARG GO_BUILDFLAGS
ARG GO_LDFLAGS
RUN --mount=type=cache,target=/.cache GOOS=linux GOARCH=${TARGETARCH} go build ${GO_BUILDFLAGS} -ldflags "${GO_LDFLAGS} -X main.TalosRelease=${TALOS_RELEASE}" -o /manager ./app/sidero-controller-manager
RUN chmod +x /manager

FROM base AS agent-build-amd64
ARG GO_BUILDFLAGS
ARG GO_LDFLAGS
RUN --mount=type=cache,target=/.cache GOOS=linux GOARCH=amd64 go build ${GO_BUILDFLAGS} -ldflags "${GO_LDFLAGS}" -o /agent ./app/sidero-controller-manager/cmd/agent
RUN chmod +x /agent

FROM base AS agent-build-arm64
ARG GO_BUILDFLAGS
ARG GO_LDFLAGS
RUN --mount=type=cache,target=/.cache GOOS=linux GOARCH=arm64 go build ${GO_BUILDFLAGS} -ldflags "${GO_LDFLAGS}" -o /agent ./app/sidero-controller-manager/cmd/agent
RUN chmod +x /agent

FROM base AS initramfs-archive-amd64
WORKDIR /initramfs
COPY --from=pkg-ca-certificates / .
COPY --from=pkg-musl-amd64 / .
COPY --from=pkg-libressl-amd64 / .
COPY --from=pkg-ipmitool-amd64 / .
COPY --from=agent-build-amd64 /agent ./init
COPY --from=pkg-linux-firmware-amd64 /lib/firmware/bnx2 ./lib/firmware/bnx2
COPY --from=pkg-linux-firmware-amd64 /lib/firmware/bnx2x ./lib/firmware/bnx2x
RUN set -o pipefail && find . 2>/dev/null | cpio -H newc -o | xz -v -C crc32 -0 -e -T 0 -z >/initramfs.xz

FROM base AS initramfs-archive-arm64
WORKDIR /initramfs
COPY --from=pkg-ca-certificates / .
COPY --from=pkg-musl-arm64 / .
COPY --from=pkg-libressl-arm64 / .
COPY --from=pkg-ipmitool-arm64 / .
COPY --from=agent-build-arm64 /agent ./init
COPY --from=pkg-linux-firmware-arm64 /lib/firmware/bnx2 ./lib/firmware/bnx2
COPY --from=pkg-linux-firmware-arm64 /lib/firmware/bnx2x ./lib/firmware/bnx2x
RUN set -o pipefail && find . 2>/dev/null | cpio -H newc -o | xz -v -C crc32 -0 -e -T 0 -z >/initramfs.xz

FROM scratch AS sidero-controller-manager-image
COPY --from=pkg-ca-certificates / /
COPY --from=pkg-fhs / /
COPY --from=pkg-musl / /
COPY --from=pkg-libressl / /
COPY --from=pkg-liblzma / /
COPY --from=pkg-ipmitool / /
COPY --from=pkg-ipxe-amd64 /usr/libexec/ /var/lib/sidero/ipxe/amd64
COPY --from=pkg-ipxe-arm64 /usr/libexec/ /var/lib/sidero/ipxe/arm64
COPY --from=pkg-ipxe /usr/libexec/zbin /bin/zbin
COPY --from=initramfs-archive-amd64 /initramfs.xz /var/lib/sidero/env/agent-amd64/initramfs.xz
COPY --from=initramfs-archive-arm64 /initramfs.xz /var/lib/sidero/env/agent-arm64/initramfs.xz
COPY --from=pkg-kernel-amd64 /boot/vmlinuz /var/lib/sidero/env/agent-amd64/vmlinuz
COPY --from=pkg-kernel-arm64 /boot/vmlinuz /var/lib/sidero/env/agent-arm64/vmlinuz
COPY --from=build-sidero-controller-manager /manager /manager

FROM sidero-controller-manager-image AS sidero-controller-manager
LABEL org.opencontainers.image.source https://github.com/talos-systems/sidero
ENTRYPOINT [ "/manager" ]

FROM base AS unit-tests-runner
ARG TEST_PKGS
RUN --mount=type=cache,target=/.cache --mount=type=cache,id=testspace,target=/tmp --mount=type=cache,target=/root/.cache/go-build go test -v -covermode=atomic -coverprofile=coverage.txt -count 1 ${TEST_PKGS}
#
FROM scratch AS unit-tests
COPY --from=unit-tests-runner /src/coverage.txt /coverage.txt
#
# The unit-tests-race target performs tests with race detector.
#
FROM base AS unit-tests-race
ENV CGO_ENABLED 1
ARG TEST_PKGS
RUN --mount=type=cache,target=/.cache --mount=type=cache,target=/root/.cache/go-build go test -v -count 1 -race ${TEST_PKGS}
#
# The lint target performs linting on the source code.
#
FROM base AS lint-go
COPY .golangci.yml .
ENV GOGC=50
ENV GOLANGCI_LINT_CACHE /.cache/lint
RUN --mount=type=cache,target=/.cache golangci-lint run --config .golangci.yml
ARG MODULE
RUN --mount=type=cache,target=/.cache FILES="$(gofumports -l -local ${MODULE} .)" && test -z "${FILES}" || (echo -e "Source code is not formatted with 'gofumports -w -local ${MODULE} .':\n${FILES}"; exit 1)
#
# The fmt target formats the source code.
#
FROM base AS fmt-build
ARG MODULE
RUN --mount=type=cache,target=/.cache gofumports -w -local ${MODULE} .
#
FROM scratch AS fmt
COPY --from=fmt-build /src /

#
# The markdownlint target performs linting on Markdown files.
#
FROM node:16.6.1-alpine AS lint-markdown
RUN apk add --no-cache findutils
RUN npm i -g markdownlint-cli@0.23.2
RUN npm i -g textlint@11.7.6
RUN npm i -g textlint-filter-rule-comments@1.2.2
RUN npm i -g textlint-rule-one-sentence-per-line@1.0.2
WORKDIR /src
COPY . .
RUN markdownlint \
    --ignore '**/LICENCE.md' \
    --ignore '**/CHANGELOG.md' \
    --ignore '**/CODE_OF_CONDUCT.md' \
    --ignore '**/node_modules/**' \
    --ignore '**/hack/chglog/**' \
    .
RUN find . \
    -name '*.md' \
    -not -path './LICENCE.md' \
    -not -path './CHANGELOG.md' \
    -not -path './CODE_OF_CONDUCT.md' \
    -not -path '*/node_modules/*' \
    -not -path './hack/chglog/**' \
    -print0 \
    | xargs -0 textlint

#
# The sfyra-build target builds the Sfyra source.
#
FROM base AS sfyra-base
WORKDIR /src/sfyra
COPY ./sfyra/go.mod ./
COPY ./sfyra/go.sum ./
RUN --mount=type=cache,target=/.cache go mod download
RUN --mount=type=cache,target=/.cache go mod verify
COPY ./sfyra/ ./
RUN --mount=type=cache,target=/.cache go list -mod=readonly all >/dev/null
RUN --mount=type=cache,target=/.cache ! go mod tidy -v 2>&1 | grep .

FROM sfyra-base AS lint-sfyra
COPY .golangci.yml .
ENV GOGC=50
ENV GOLANGCI_LINT_CACHE /.cache/lint
RUN --mount=type=cache,target=/.cache golangci-lint run --config .golangci.yml
ARG MODULE
RUN --mount=type=cache,target=/.cache FILES="$(gofumports -l -local ${MODULE} .)" && test -z "${FILES}" || (echo -e "Source code is not formatted with 'gofumports -w -local ${MODULE} .':\n${FILES}"; exit 1)

FROM sfyra-base AS sfyra-build
WORKDIR /src/sfyra/cmd/sfyra
ARG TALOS_RELEASE
ARG DEFAULT_K8S_VERSION
ARG SFYRA_CMD_PKG=github.com/talos-systems/sidero/sfyra/cmd/sfyra/cmd
ARG GO_BUILDFLAGS
ARG GO_LDFLAGS
RUN --mount=type=cache,target=/.cache GOOS=linux go build ${GO_BUILDFLAGS} -ldflags "${GO_LDFLAGS} -X ${SFYRA_CMD_PKG}.TalosRelease=${TALOS_RELEASE} -X ${SFYRA_CMD_PKG}.KubernetesVersion=${DEFAULT_K8S_VERSION}" -o /sfyra
RUN chmod +x /sfyra

FROM scratch AS sfyra
COPY --from=sfyra-build /sfyra /sfyra
