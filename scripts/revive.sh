#!/usr/bin/env bash

if ! which revive > /dev/null; then
    echo "==> Installing revive"
    go get -u github.com/mgechev/revive@master
fi

echo "==> Checking revive..."
$(go env GOPATH)/bin/revive --config ./.revive.toml --formatter stylish ./cmd/...
$(go env GOPATH)/bin/revive --config ./.revive.toml --formatter stylish ./internal/...
$(go env GOPATH)/bin/revive --config ./.revive.toml --formatter stylish ./pkg/...
