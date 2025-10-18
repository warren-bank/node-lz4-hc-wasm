#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

PRJ_DIR="${DIR}/../.."

[ -d "${PRJ_DIR}/dist" ] || mkdir "${PRJ_DIR}/dist"
[ -f "${PRJ_DIR}/dist/lz4.wasm" ] && rm -f "${PRJ_DIR}/dist/lz4.wasm"
[ -f "${PRJ_DIR}/dist/wasm_exec.js" ] && rm -f "${PRJ_DIR}/dist/wasm_exec.js"

# --------------------------------------------------------------------

export GOOS='js'
export GOARCH='wasm'
export GO111MODULE='on'

cd "${PRJ_DIR}/src"
go build -o "${PRJ_DIR}/dist/lz4.wasm" "main.go"

# --------------------------------------------------------------------

if [ -z "$GOROOT" ];then
  GR="$(dirname $(which go))"
  if [ -f "${GR}/../lib/wasm/wasm_exec.js" ];then
    GOROOT="${GR}/.."
  elif [ -f "${GR}/../golang/env.sh" ];then
    source "${GR}/../golang/env.sh"
  fi
fi

if [ -n "$GOROOT" -a -f "${GOROOT}/lib/wasm/wasm_exec.js" ];then
  cp "${GOROOT}/lib/wasm/wasm_exec.js" "${PRJ_DIR}/dist/wasm_exec.js"
fi
