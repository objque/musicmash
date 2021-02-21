override RELEASE="$(git tag -l --points-at HEAD)"
override COMMIT="$(shell git rev-parse --short HEAD)"
override BUILD_TIME="$(shell date -u '+%Y-%m-%dT%H:%M:%S')"
override VERSION=v3

all:

clean:
	rm dist/musicmash || true
	rm dist/musicmashctl || true

build: clean
	go build -ldflags="-s -w" -v -o dist/musicmash ./cmd/musicmash/...
	go build -ldflags="-s -w" -v -o dist/musicmashctl ./cmd/musicmashctl/...

build-arm7: clean
	if ! which arm-linux-gnueabi-gcc > /dev/null; then \
		echo "you must have gcc-arm-linux-gnueabi/stable package installed to build musicmash for arm7:"; \
		echo "\n  apt update && apt install -y gcc-arm-linux-gnueabi/stable\n"; \
		exit 1; \
	fi

	env GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=1 CC=arm-linux-gnueabi-gcc go build -ldflags="-s -w" -v -o ./dist/musicmash ./cmd/musicmash/...
	env GOOS=linux GOARCH=arm GOARM=7 go build -ldflags="-s -w" -v -o ./dist/musicmashctl ./cmd/musicmashctl/...

install:
	go install -v ./cmd/...

test t:
	go test -v -short ./...

test-full tf:
	go test -v -p 1 ./...

update-deps:
	go get -u ./...
	go mod vendor

image:
	docker build \
		--build-arg RELEASE=${RELEASE} \
		--build-arg COMMIT=${COMMIT} \
		--build-arg BUILD_TIME=${BUILD_TIME} \
		-t $(REGISTRY_REPO):$(VERSION) .

compose:
	docker-compose up --build -d

lint l:
	bash ./scripts/revive.sh
	bash ./scripts/golangci-lint.sh

run: install
	musicmash --db-auto-migrate=true --db-migrations-dir=./migrations --config musicmash.example.yaml

ensure-go-migrate-installed:
	bash ./scripts/install-go-migrate.sh

# show latest applied migration
db-status: ensure-go-migrate-installed
	migrate -path migrations -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -verbose version

# apply migration up
db-up: ensure-go-migrate-installed
	migrate -path migrations -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -verbose up

# apply migration down
db-down: ensure-go-migrate-installed
	migrate -path migrations -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -verbose down
