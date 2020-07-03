#!/usr/bin/env bash

if ! which golangci-lint > /dev/null; then
    echo "==> Installing golangci-lint"
    # do not forget to bump version in gh/action
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $GOPATH/bin v1.25.0
fi

echo "==> Checking golangci-ling..."
golangci-lint -v run
