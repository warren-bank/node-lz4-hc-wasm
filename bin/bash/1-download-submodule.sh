#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

cd "${DIR}/../.."

git submodule add -b 'v4' --name 'lz4' --depth '1' 'https://github.com/pierrec/lz4.git' 'src/lz4'
git submodule update --init --checkout --recursive 'src/lz4'

cd "src"
go mod edit -replace "github.com/pierrec/lz4/v4=./lz4"
go mod tidy
