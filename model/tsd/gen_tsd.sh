#!/bin/sh
PROTOC_GEN_TS_PATH="./node_modules/.bin/protoc-gen-ts"
OUT_DIR="."
 
protoc \
    -I.. \
    --plugin="protoc-gen-ts=${PROTOC_GEN_TS_PATH}" \
    --js_out="import_style=commonjs,binary:${OUT_DIR}" \
    --ts_out="${OUT_DIR}" \
    syzoj.model.proto syzoj.api.proto
