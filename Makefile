override RELEASE="$(git tag -l --points-at HEAD)"
override COMMIT="$(shell git rev-parse --short HEAD)"
override BUILD_TIME="$(shell date -u '+%Y-%m-%dT%H:%M:%S')"
override VERSION="v3"

all:

clean:
	rm bin/musicmash || true
	rm bin/musicmashctl || true

build: clean
	go build -ldflags="-s -w" -v -o bin/musicmash ./cmd/musicmash/...
	go build -ldflags="-s -w" -v -o bin/musicmashctl ./cmd/musicmashctl/...

install:
	go install -v ./cmd/...

t tests: install
	go test -v ./internal/...

add-ssh-key:
	openssl aes-256-cbc -K $(encrypted_a4311917bb34_key) -iv $(encrypted_a4311917bb34_iv) -in travis_key.enc -out /tmp/travis_key -d
	chmod 600 /tmp/travis_key
	ssh-add /tmp/travis_key

docker-login:
	docker login -u $(REGISTRY_USER) -p $(REGISTRY_PASS)

docker-build:
	docker build \
		--build-arg RELEASE=${RELEASE} \
		--build-arg COMMIT=${COMMIT} \
		--build-arg BUILD_TIME=${BUILD_TIME} \
		-t $(REGISTRY_REPO):$(VERSION) .

docker-push: docker-login
	docker push $(REGISTRY_REPO):$(VERSION)

deploy:
	ssh -o "StrictHostKeyChecking no" $(HOST_USER)@$(HOST) make run-music

deploy-staging:
	ssh -o "StrictHostKeyChecking no" $(STAGING_USER)@$(STAGING_HOST) make run-music

lint-all l:
	bash ./scripts/revive.sh
	bash ./scripts/golangci-lint.sh

install-arm7-deps iarm7:
	apt update && apt install -y gcc-arm-linux-gnueabi/stable

build-arm7: clean
	# you must have gcc-arm-linux-gnueabi/stable package installed to build musicmash for arm7:
	# make iarm
	env GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=1 CC=arm-linux-gnueabi-gcc go build -ldflags="-s -w" -v -o ./bin/musicmash ./cmd/musicmash/...
	env GOOS=linux GOARCH=arm GOARM=7 go build -ldflags="-s -w" -v -o ./bin/musicmashctl ./cmd/musicmashctl/...

api:
	musicmash --fetcher-enabled=false --notifier-enabled=false
