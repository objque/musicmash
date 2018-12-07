#!/usr/bin/env bash

if ! which go-consistent > /dev/null; then
    echo "==> Installing go-consistent..."
    go get -u github.com/Quasilyte/go-consistent
fi

echo "==> Checking go-consistent..."
go-consistent -v ./...
