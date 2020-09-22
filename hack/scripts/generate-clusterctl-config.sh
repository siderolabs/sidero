#!/bin/bash

set -eou pipefail

mkdir -p ~/.cluster-api

cat > ~/.cluster-api/clusterctl.yaml <<EOF
providers:
  - name: "sidero"
    url: "file://${COMPONENTS_YAML}"
    type: "InfrastructureProvider"
EOF
