#!/bin/bash

set -eou pipefail

mkdir -p "$(dirname "${CLUSTERCTL_CONFIG}")"

cat > "${CLUSTERCTL_CONFIG}" <<EOF
providers:
  - name: "sidero"
    url: "file://${COMPONENTS_YAML}"
    type: "InfrastructureProvider"
# temporary, see https://github.com/kubernetes-sigs/cluster-api/issues/6051
cert-manager:
  url: "https://github.com/cert-manager/cert-manager/releases/latest/cert-manager.yaml"
EOF
