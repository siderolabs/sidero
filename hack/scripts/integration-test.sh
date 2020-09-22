#!/bin/bash

set -eou pipefail

INTEGRATION_TEST="${ARTIFACTS}/sfyra"

BOOTSTRAP_VMLINUZ="${ARTIFACTS}/${TALOS_RELEASE}/vmlinuz"
BOOTSTRAP_INITRAMFS="${ARTIFACTS}/${TALOS_RELEASE}/initramfs.xz"
BOOTSTRAP_INSTALLER="docker.io/autonomy/installer:${TALOS_RELEASE}"

TALOSCTL="${ARTIFACTS}/${TALOS_RELEASE}/talosctl-linux-amd64"

chmod +x "${TALOSCTL}"

function build_registry_mirrors {
  if [[ "${CI:-false}" == "true" ]]; then
    REGISTRY_MIRROR_FLAGS=

    for registry in docker.io k8s.gcr.io quay.io gcr.io; do
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

${PREFIX} "${INTEGRATION_TEST}" \
    -bootstrap-initramfs "${BOOTSTRAP_INITRAMFS}" \
    -bootstrap-vmlinuz "${BOOTSTRAP_VMLINUZ}" \
    -bootstrap-installer "${BOOTSTRAP_INSTALLER}" \
    -talosctl-path "${TALOSCTL}" \
    ${REGISTRY_MIRROR_FLAGS} \
    -test.v
