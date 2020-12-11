#!/usr/bin/env bash

if ! which golangci-lint > /dev/null; then
    echo "==> Installing golangci-lint"
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.25.1
fi

echo "==> Checking golangci-ling..."
$(go env GOPATH)/bin/golangci-lint -v run
