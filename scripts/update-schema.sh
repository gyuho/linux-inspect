#!/usr/bin/env bash

set -e

if ! [[ "$0" =~ "scripts/update-schema.sh" ]]; then
  echo "must be run from repository root"
  exit 255
fi

go run ./cmd/generate-df/main.go
go run ./cmd/generate-etc/main.go
go run ./cmd/generate-proc/main.go
go run ./cmd/generate-top/main.go
