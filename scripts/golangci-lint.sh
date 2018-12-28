#!/usr/bin/env bash

if ! which golangci-lint > /dev/null; then
    echo "==> Installing golangci-lint"
    curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $GOPATH/bin latest
fi

echo "==> Checking golangci-ling..."
golangci-lint -v run
