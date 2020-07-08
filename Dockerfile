# syntax = docker/dockerfile-upstream:1.1.4-experimental

# The base target provides the base for running various tasks against the source
# code

FROM golang:1.13 AS base
ENV GO111MODULE on
ENV GOPROXY https://proxy.golang.org
ENV CGO_ENABLED 0
WORKDIR /tmp
RUN go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.3.0
RUN go get k8s.io/code-generator/cmd/conversion-gen@v0.18.2
WORKDIR /src
COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download
RUN go mod verify
COPY ./ ./
RUN go list -mod=readonly all >/dev/null
RUN ! go mod tidy -v 2>&1 | grep .

FROM base AS manifests-build
RUN controller-gen \
  crd:crdVersions=v1 paths="./internal/app/cluster-api-provider/api/..." output:crd:dir="./internal/app/cluster-api-provider/config/crd/bases" \
  rbac:roleName=manager-role paths="./internal/app/cluster-api-provider/controllers/..." output:rbac:dir="./internal/app/cluster-api-provider/config/rbac" \
  webhook output:webhook:dir="./internal/app/cluster-api-provider/config/webhook"
RUN controller-gen \
  crd:crdVersions=v1 paths="./internal/app/metal-controller-manager/api/..." output:crd:dir="./internal/app/metal-controller-manager/config/crd/bases" \
  rbac:roleName=manager-role paths="./internal/app/metal-controller-manager/controllers/..." output:rbac:dir="./internal/app/metal-controller-manager/config/rbac" \
  webhook output:webhook:dir="./internal/app/metal-controller-manager/config/webhook"
FROM scratch AS manifests
COPY --from=manifests-build /src/internal/app/cluster-api-provider/config ./internal/app/cluster-api-provider/config

FROM base AS generate-build
RUN controller-gen object:headerFile="./hack/boilerplate.go.txt" paths="./..."
RUN	conversion-gen --input-dirs="./internal/app/cluster-api-provider/api/v1alpha2" --output-base ./ --output-file-base="zz_generated.conversion" --go-header-file="./hack/boilerplate.go.txt"
FROM scratch AS generate
COPY --from=generate-build /src/internal/app/cluster-api-provider/api ./internal/app/cluster-api-provider/api
COPY --from=generate-build /src/internal/app/metal-controller-manager/api ./internal/app/metal-controller-manager/api

FROM k8s.gcr.io/hyperkube:v1.17.0 AS release-build
RUN apt update -y \
  && apt install -y curl \
  && curl -LO https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize%2Fv3.5.4/kustomize_v3.5.4_linux_amd64.tar.gz \
  && tar -xf kustomize_v3.5.4_linux_amd64.tar.gz -C /usr/local/bin \
  && rm kustomize_v3.5.4_linux_amd64.tar.gz
COPY ./config ./config
COPY ./internal/app/cluster-api-provider/config ./internal/app/cluster-api-provider/config
COPY ./internal/app/metal-controller-manager/config ./internal/app/metal-controller-manager/config
COPY ./internal/app/metal-metadata-server/config ./internal/app/metal-metadata-server/config
ARG REGISTRY_AND_USERNAME
ARG TAG
RUN cd ./internal/app/cluster-api-provider/config/manager \
  && kustomize edit set image controller=${REGISTRY_AND_USERNAME}/cluster-api-provider:${TAG} \
  && cd - \
  && kustomize build config > /infrastructure-components.yaml \
  && cp ./internal/app/cluster-api-provider/config/metadata/metadata.yaml /metadata.yaml

FROM scratch AS release
COPY --from=release-build /infrastructure-components.yaml /infrastructure-components.yaml
COPY --from=release-build /metadata.yaml /metadata.yaml

FROM base AS build-cluster-api-provider
RUN --mount=type=cache,target=/root/.cache/go-build GOOS=linux go build -ldflags "-s -w" -o /manager ./internal/app/cluster-api-provider
RUN chmod +x /manager

## TODO(rsmitty): make bmc pkg and move to autonomy image
FROM scratch AS cluster-api-provider
COPY --from=docker.io/autonomy/ca-certificates:ffdacf0 / /
COPY --from=docker.io/autonomy/fhs:ffdacf0 / /
COPY --from=docker.io/autonomy/musl:ffdacf0 / /
COPY --from=docker.io/autonomy/libressl:ffdacf0 / /
COPY --from=docker.io/autonomy/ipmitool:ffdacf0 / /
COPY --from=build-cluster-api-provider /manager /manager
ENTRYPOINT [ "/manager" ]

FROM base AS build-metal-controller-manager
RUN --mount=type=cache,target=/root/.cache/go-build GOOS=linux go build -ldflags "-s -w" -o /manager ./internal/app/metal-controller-manager
RUN chmod +x /manager

FROM alpine:3.11 AS assets
RUN apk add --no-cache curl
RUN curl -s -o /undionly.kpxe http://boot.ipxe.org/undionly.kpxe
RUN curl -s -o /ipxe.efi http://boot.ipxe.org/ipxe.efi

FROM base AS agent-build
RUN --mount=type=cache,target=/root/.cache/go-build GOOS=linux go build -ldflags "-s -w" -o /agent ./internal/app/metal-controller-manager/cmd/agent
RUN chmod +x /agent

FROM scratch AS agent
COPY --from=docker.io/autonomy/ca-certificates:v0.2.0 / /
COPY --from=docker.io/autonomy/fhs:v0.2.0 / /
COPY --from=agent-build /agent /agent
ENTRYPOINT [ "/agent" ]

FROM autonomy/tools:v0.2.0 AS initramfs-archive
ENV PATH /toolchain/bin
RUN [ "/toolchain/bin/mkdir", "/bin" ]
RUN [ "ln", "-s", "/toolchain/bin/bash", "/bin/sh" ]
WORKDIR /initramfs
COPY --from=agent /agent ./init
COPY --from=docker.io/autonomy/linux-firmware:v0.2.0 /lib/firmware/bnx2 ./lib/firmware/bnx2
COPY --from=docker.io/autonomy/linux-firmware:v0.2.0 /lib/firmware/bnx2x ./lib/firmware/bnx2x
RUN set -o pipefail && find . 2>/dev/null | cpio -H newc -o | xz -v -C crc32 -0 -e -T 0 -z >/initramfs.xz

FROM scratch AS initramfs
COPY --from=initramfs-archive /initramfs.xz /initramfs.xz

FROM scratch AS metal-controller-manager
COPY --from=docker.io/autonomy/ca-certificates:v0.2.0 / /
COPY --from=docker.io/autonomy/fhs:v0.2.0 / /
COPY --from=assets /undionly.kpxe /var/lib/sidero/tftp/undionly.kpxe
COPY --from=assets /undionly.kpxe /var/lib/sidero/tftp/undionly.kpxe.0
COPY --from=assets /ipxe.efi /var/lib/sidero/tftp/ipxe.efi
COPY --from=initramfs /initramfs.xz /var/lib/sidero/env/discovery/initramfs.xz
ADD https://github.com/talos-systems/talos/releases/download/v0.4.1/vmlinuz /var/lib/sidero/env/discovery/vmlinuz
COPY --from=build-metal-controller-manager /manager /manager
ENTRYPOINT [ "/manager" ]

FROM base AS build-metal-metadata-server
RUN --mount=type=cache,target=/root/.cache/go-build GOOS=linux go build -ldflags "-s -w" -o /metal-metadata-server ./internal/app/metal-metadata-server
RUN chmod +x /metal-metadata-server

FROM scratch AS metal-metadata-server
COPY --from=docker.io/autonomy/ca-certificates:v0.1.0 / /
COPY --from=docker.io/autonomy/fhs:v0.1.0 / /
COPY --from=build-metal-metadata-server /metal-metadata-server /metal-metadata-server
ENTRYPOINT [ "/metal-metadata-server" ]

FROM base AS unit-tests-runner
ARG PKGS
RUN --mount=type=cache,id=testspace,target=/tmp --mount=type=cache,target=/root/.cache/go-build go test -v -covermode=atomic -coverprofile=coverage.txt -count 1 ${PKGS}
#
FROM scratch AS unit-tests
COPY --from=unit-tests-runner /src/coverage.txt /coverage.txt
#
# The unit-tests-race target performs tests with race detector.
#
FROM base AS unit-tests-race
ENV CGO_ENABLED 1
ARG PKGS
RUN --mount=type=cache,target=/root/.cache/go-build go test -v -count 1 -race ${PKGS}
#
# The lint target performs linting on the source code.
#
FROM base AS lint-go
ENV GOGC=50
RUN --mount=type=cache,target=/root/.cache/go-build /go/bin/golangci-lint run
ARG MODULE
RUN FILES="$(gofumports -l -local ${MODULE} .)" && test -z "${FILES}" || (echo -e "Source code is not formatted with 'gofumports -w -local ${MODULE} .':\n${FILES}"; exit 1)
#
# The fmt target formats the source code.
#
FROM base AS fmt-build
ARG MODULE
RUN gofumports -w -local ${MODULE} .
#
FROM scratch AS fmt
COPY --from=fmt-build /src /
#
# The markdownlint target performs linting on Markdown files.
#
FROM node:8.16.1-alpine AS lint-markdown
RUN npm install -g markdownlint-cli
RUN npm i sentences-per-line
WORKDIR /src
COPY --from=base /src .
RUN markdownlint --rules /node_modules/sentences-per-line/index.js .
