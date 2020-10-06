#!/bin/bash

set -eou pipefail

INTEGRATION_TEST="${ARTIFACTS}/sfyra"

TALOSCTL="${ARTIFACTS}/${TALOS_RELEASE}/talosctl-linux-amd64"

chmod +x "${TALOSCTL}"

function build_registry_mirrors {
  if [[ "${CI:-false}" == "true" ]]; then
    REGISTRY_MIRROR_FLAGS=

    for registry in docker.io k8s.gcr.io quay.io gcr.io registry.dev.talos-systems.io; do
      local service="registry-${registry//./-}.ci.svc"
      local addr=`python3 -c "import socket; print(socket.gethostbyname('${service}'))"`

      REGISTRY_MIRROR_FLAGS="${REGISTRY_MIRROR_FLAGS} --registry-mirrors ${registry}=http://${addr}:5000"
    done

    local addr=`python3 -c "import socket; print(socket.gethostbyname('registry.ci.svc'))"`
    REGISTRY_MIRROR_FLAGS="${REGISTRY_MIRROR_FLAGS} --registry-mirrors registry.ci.svc:5000=http://${addr}:5000"
  else
    REGISTRY_MIRROR_FLAGS=${REGISTRY_MIRROR_FLAGS:-}
  fi
}

build_registry_mirrors

if [ "$EUID" -ne 0 ]; then
    PREFIX="sudo -E"
else
    PREFIX=
fi

${PREFIX} "${INTEGRATION_TEST}" test integration \
    --talosctl-path "${TALOSCTL}" \
    --clusterctl-config "${CLUSTERCTL_CONFIG}" \
    ${REGISTRY_MIRROR_FLAGS}
