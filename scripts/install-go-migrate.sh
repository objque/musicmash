#!/usr/bin/env bash

if ! which migrate > /dev/null; then
    echo "==> Installing go-migrate"
		curl -L https://github.com/golang-migrate/migrate/releases/download/v4.13.0/migrate.linux-amd64.tar.gz | tar xvz
		mv migrate.linux-amd64 $GOPATH/bin/migrate
fi
