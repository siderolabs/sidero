#!/bin/bash

set -eou pipefail

mkdir -p "$(dirname "${CLUSTERCTL_CONFIG}")"

cat > "${CLUSTERCTL_CONFIG}" <<EOF
providers:
  - name: "talos"
    url: "https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/releases/latest/bootstrap-components.yaml"
    type: "BootstrapProvider"
  - name: "talos"
    url: "https://github.com/siderolabs/cluster-api-control-plane-provider-talos/releases/latest/control-plane-components.yaml"
    type: "ControlPlaneProvider"
  - name: "sidero"
    url: "file://${COMPONENTS_YAML}"
    type: "InfrastructureProvider"
EOF
