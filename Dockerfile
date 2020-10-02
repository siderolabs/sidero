# syntax = docker/dockerfile-upstream:1.1.4-experimental

# The base target provides the base for running various tasks against the source
# code

FROM golang:1.15 AS base
ENV GO111MODULE on
ENV GOPROXY https://proxy.golang.org
ENV CGO_ENABLED 0
WORKDIR /tmp
RUN apt-get update \
  && apt-get install -y unzip \
  && curl -L https://github.com/protocolbuffers/protobuf/releases/download/v3.7.1/protoc-3.7.1-linux-x86_64.zip -o /tmp/protoc.zip \
  && unzip -o /tmp/protoc.zip -d /usr/local bin/protoc \
  && unzip -o /tmp/protoc.zip -d /usr/local 'include/*' \
  && go get github.com/golang/protobuf/protoc-gen-go@v1.3
RUN go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.3.0
RUN go get k8s.io/code-generator/cmd/conversion-gen@v0.18.2
RUN curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | bash -s -- -b /usr/local/bin v1.28.0
RUN cd $(mktemp -d) \
  && go mod init tmp \
  && go get mvdan.cc/gofumpt/gofumports@abc0db2c416aca0f60ea33c23c76665f6e7ba0b6 \
  && mv /go/bin/gofumports /usr/local/bin/gofumports
WORKDIR /src
COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download
RUN go mod verify
COPY ./app/ ./app/
COPY ./hack/ ./hack/
COPY ./internal/ ./internal/
RUN go list -mod=readonly all >/dev/null
RUN ! go mod tidy -v 2>&1 | grep .

FROM base AS manifests-build
RUN controller-gen \
  crd:crdVersions=v1 paths="./app/cluster-api-provider-sidero/api/..." output:crd:dir="./app/cluster-api-provider-sidero/config/crd/bases" \
  rbac:roleName=manager-role paths="./app/cluster-api-provider-sidero/controllers/..." output:rbac:dir="./app/cluster-api-provider-sidero/config/rbac" \
  webhook output:webhook:dir="./app/cluster-api-provider-sidero/config/webhook"
RUN controller-gen \
  crd:crdVersions=v1 paths="./app/metal-controller-manager/api/..." output:crd:dir="./app/metal-controller-manager/config/crd/bases" \
  rbac:roleName=manager-role paths="./app/metal-controller-manager/controllers/..." output:rbac:dir="./app/metal-controller-manager/config/rbac" \
  webhook output:webhook:dir="./app/metal-controller-manager/config/webhook"

FROM scratch AS manifests
COPY --from=manifests-build /src/app/cluster-api-provider-sidero/config ./app/cluster-api-provider-sidero/config
COPY --from=manifests-build /src/app/metal-controller-manager/config ./app/metal-controller-manager/config

FROM base AS generate-build
COPY ./app/metal-controller-manager/internal/api/api.proto \
  /src/app/metal-controller-manager/internal/api/api.proto
RUN protoc -I/src/app/metal-controller-manager/internal/api \
  --go_out=plugins=grpc,paths=source_relative:/src/app/metal-controller-manager/internal/api \
  api.proto
RUN controller-gen object:headerFile="./hack/boilerplate.go.txt" paths="./..."
RUN	conversion-gen --input-dirs="./app/cluster-api-provider-sidero/api/v1alpha2" --output-base ./ --output-file-base="zz_generated.conversion" --go-header-file="./hack/boilerplate.go.txt"
FROM scratch AS generate
COPY --from=generate-build /src/app/cluster-api-provider-sidero/api ./app/cluster-api-provider-sidero/api
COPY --from=generate-build /src/app/metal-controller-manager/api ./app/metal-controller-manager/api
COPY --from=generate-build /src/app/metal-controller-manager/internal/api ./app/metal-controller-manager/internal/api

FROM k8s.gcr.io/hyperkube:v1.17.0 AS release-build
RUN apt update -y \
  && apt install -y curl \
  && curl -LO https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize%2Fv3.5.4/kustomize_v3.5.4_linux_amd64.tar.gz \
  && tar -xf kustomize_v3.5.4_linux_amd64.tar.gz -C /usr/local/bin \
  && rm kustomize_v3.5.4_linux_amd64.tar.gz
COPY ./config ./config
COPY ./templates ./templates
COPY ./app/cluster-api-provider-sidero/config ./app/cluster-api-provider-sidero/config
COPY ./app/metal-controller-manager/config ./app/metal-controller-manager/config
COPY ./app/metal-metadata-server/config ./app/metal-metadata-server/config
ARG REGISTRY_AND_USERNAME
ARG TAG
RUN cd ./app/cluster-api-provider-sidero/config/manager \
  && kustomize edit set image controller=${REGISTRY_AND_USERNAME}/cluster-api-provider-sidero:${TAG}
RUN cd ./app/metal-controller-manager/config/manager \
  && kustomize edit set image controller=${REGISTRY_AND_USERNAME}/metal-controller-manager:${TAG}
RUN cd ./app/metal-metadata-server/config/server \
  && kustomize edit set image server=${REGISTRY_AND_USERNAME}/metal-metadata-server:${TAG}
RUN kustomize build config > /infrastructure-components.yaml \
  && cp ./config/metadata/metadata.yaml /metadata.yaml \
  && cp ./templates/cluster-template.yaml /cluster-template.yaml

FROM scratch AS release
ARG TAG
COPY --from=release-build /infrastructure-components.yaml /infrastructure-sidero/${TAG}/infrastructure-components.yaml
COPY --from=release-build /metadata.yaml /infrastructure-sidero/${TAG}/metadata.yaml
COPY --from=release-build /cluster-template.yaml /infrastructure-sidero/${TAG}/cluster-template.yaml

FROM base AS build-cluster-api-provider-sidero
RUN --mount=type=cache,target=/root/.cache/go-build GOOS=linux go build -ldflags "-s -w" -o /manager ./app/cluster-api-provider-sidero
RUN chmod +x /manager

## TODO(rsmitty): make bmc pkg and move to autonomy image
FROM scratch AS cluster-api-provider-sidero
COPY --from=docker.io/autonomy/ca-certificates:ffdacf0 / /
COPY --from=docker.io/autonomy/fhs:ffdacf0 / /
COPY --from=docker.io/autonomy/musl:ffdacf0 / /
COPY --from=docker.io/autonomy/libressl:ffdacf0 / /
COPY --from=docker.io/autonomy/ipmitool:ffdacf0 / /
COPY --from=build-cluster-api-provider-sidero /manager /manager
ENTRYPOINT [ "/manager" ]

FROM base AS build-metal-controller-manager
RUN --mount=type=cache,target=/root/.cache/go-build GOOS=linux go build -ldflags "-s -w" -o /manager ./app/metal-controller-manager
RUN chmod +x /manager

FROM alpine:3.11 AS assets
RUN apk add --no-cache curl
RUN curl -s -o /undionly.kpxe http://boot.ipxe.org/undionly.kpxe
RUN curl -s -o /ipxe.efi http://boot.ipxe.org/ipxe.efi

FROM base AS agent-build
RUN --mount=type=cache,target=/root/.cache/go-build GOOS=linux go build -ldflags "-s -w" -o /agent ./app/metal-controller-manager/cmd/agent
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
COPY --from=initramfs /initramfs.xz /var/lib/sidero/env/agent/initramfs.xz
ADD https://github.com/talos-systems/talos/releases/download/v0.4.1/vmlinuz /var/lib/sidero/env/agent/vmlinuz
COPY --from=build-metal-controller-manager /manager /manager
ENTRYPOINT [ "/manager" ]

FROM base AS build-metal-metadata-server
RUN --mount=type=cache,target=/root/.cache/go-build GOOS=linux go build -ldflags "-s -w" -o /metal-metadata-server ./app/metal-metadata-server
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
RUN --mount=type=cache,target=/root/.cache/go-build --mount=type=cache,target=/root/.cache/golangci-lint /usr/local/bin/golangci-lint run --enable-all --disable gochecknoglobals,gochecknoinits,lll,goerr113,funlen,nestif,maligned,gomnd,gocognit,gocyclo
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
RUN npm install -g markdownlint-cli@0.23.2
RUN npm i sentences-per-line
WORKDIR /src
COPY --from=base /src .
RUN markdownlint --ignore '**/hack/chglog/**' --rules /node_modules/sentences-per-line/index.js .
#
# The sfyra-build target builds the Sfyra source.
#
FROM base AS sfyra-base
WORKDIR /src/sfyra
COPY ./sfyra/go.mod ./
COPY ./sfyra/go.sum ./
RUN go mod download
RUN go mod verify
COPY ./sfyra/ ./
RUN go list -mod=readonly all >/dev/null
RUN ! go mod tidy -v 2>&1 | grep .

FROM sfyra-base AS lint-sfyra
ENV GOGC=50
RUN --mount=type=cache,target=/root/.cache/go-build --mount=type=cache,target=/root/.cache/golangci-lint /usr/local/bin/golangci-lint run --enable-all --disable gochecknoglobals,gochecknoinits,lll,goerr113,funlen,nestif,maligned,gomnd,gocognit,gocyclo,godox
ARG MODULE
RUN FILES="$(gofumports -l -local ${MODULE} .)" && test -z "${FILES}" || (echo -e "Source code is not formatted with 'gofumports -w -local ${MODULE} .':\n${FILES}"; exit 1)

FROM sfyra-base AS sfyra-build
WORKDIR /src/sfyra/cmd/sfyra
ARG TALOS_RELEASE
ARG SFYRA_CMD_PKG=github.com/talos-systems/sidero/sfyra/cmd/sfyra/cmd
RUN --mount=type=cache,target=/root/.cache/go-build GOOS=linux go build -ldflags "-s -w -X ${SFYRA_CMD_PKG}.TalosRelease=${TALOS_RELEASE}" -o /sfyra
RUN chmod +x /sfyra

FROM scratch AS sfyra
COPY --from=sfyra-build /sfyra /sfyra
