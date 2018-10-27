#!/usr/bin/env bash

if ! which gocritic > /dev/null; then
    echo "==> Installing gocritic..."
    go get -u github.com/go-critic/go-critic/...
fi

echo "==> Checking gocritic..."
gocritic check-project --enable=all --disable importShadow,longChain,structLitKeyOrder -withExperimental -withOpinionated .
