#!/usr/bin/env bash

if ! which golangci-lint > /dev/null; then
    version="v.1.12.3"
    echo "==> Installing golangci-lint $(version)"
    curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $GOPATH/bin $(version)
fi

echo "==> Checking golangci-ling..."
golangci-lint -v run
