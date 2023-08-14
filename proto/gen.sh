#!/usr/bin/env bash

set -eo pipefail

echo "Generating gogo proto code"
cd proto
#buf mod update
for src in "cosmos" "cosmos_proto" "ethermint" "google" "tendermint" "osmosis" "fx" "ibc" "erc20"; do
  proto_dirs=$(find "./$src" -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
  for dir in $proto_dirs; do
    proto_files=$(find "${dir}" -maxdepth 1 -name '*.proto')
    for file in $proto_files; do
      if grep "option go_package" "$file" &>/dev/null; then
        buf generate --template buf.gen.gogo.yaml "$file"
      fi
    done
  done
done
cd ..
cp -r github.com/functionx/go-sdk/* .
rm -r github.com
