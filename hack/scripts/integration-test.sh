#!/bin/bash

set -eou pipefail

# extra flags from environment; for example
#   export SFYRA_EXTRA_FLAGS="--skip-teardown"
SFYRA_EXTRA_FLAGS="${SFYRA_EXTRA_FLAGS:-}"

INTEGRATION_TEST="${ARTIFACTS}/sfyra"

TALOSCTL="${ARTIFACTS}/${TALOS_RELEASE}/talosctl-linux-amd64"

chmod +x "${TALOSCTL}"

function build_registry_mirrors {
  if [[ "${REGISTRY_MIRROR_FLAGS:-yes}" == "no" ]]; then
    REGISTRY_MIRROR_FLAGS=()

    return
  fi

  if [[ "${CI:-false}" == "true" ]]; then
    REGISTRY_MIRROR_FLAGS=()

    for registry in docker.io registry.k8s.io quay.io gcr.io ghcr.io; do
      local service="registry-${registry//./-}.ci.svc"
      addr=$(python3 -c "import socket; print(socket.gethostbyname('${service}'))")

      REGISTRY_MIRROR_FLAGS+=("--registry-mirror=${registry}=http://${addr}:5000")
    done
  fi
}

build_registry_mirrors

"${INTEGRATION_TEST}" test integration \
    --talosctl-path "${TALOSCTL}" \
    --clusterctl-config "${CLUSTERCTL_CONFIG}" \
    --power-simulated-explicit-failure-prob=0.1 \
    --power-simulated-silent-failure-prob=0.0 \
    "${REGISTRY_MIRROR_FLAGS[@]}" \
    ${SFYRA_EXTRA_FLAGS}
