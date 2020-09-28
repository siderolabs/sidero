#!/bin/bash

set -eou pipefail

mkdir -p `dirname "${CLUSTERCTL_CONFIG}"`

cat > "${CLUSTERCTL_CONFIG}" <<EOF
providers:
  - name: "sidero"
    url: "file://${COMPONENTS_YAML}"
    type: "InfrastructureProvider"
EOF
