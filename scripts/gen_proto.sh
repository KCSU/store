#!/usr/bin/env bash
# Regenerate the Go protobuf code in backend/proto from the proto submodule.
set -euo pipefail

cd "$(dirname "$0")/.."

${PROTOC:-protoc} \
    -I proto \
    --go_out=backend/proto \
    --go_opt=paths=source_relative \
    --go_opt=Mkcsu/alerts/v1/alert.proto=github.com/kcsu/store/proto/kcsu/alerts/v1\;alertsv1 \
    proto/kcsu/alerts/v1/alert.proto
