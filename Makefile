all:

clean:
	rm bin/musicmash || true

build: clean
	GOOS=linux GOARCH=amd64 go build -v -a -installsuffix cgo -gcflags "all=-trimpath=$(GOPATH)" -o bin/musicmash cmd/musicmash.go

rgo:
	go get -u github.com/kyoh86/richgo

install:
	go install -v ./cmd

t tests: install
	go test -v ./internal/...

add-ssh-key:
	openssl aes-256-cbc -K $(encrypted_a4311917bb34_key) -iv $(encrypted_a4311917bb34_iv) -in travis_key.enc -out /tmp/travis_key -d
	chmod 600 /tmp/travis_key
	ssh-add /tmp/travis_key

docker-login:
	docker login -u $(REGISTRY_USER) -p $(REGISTRY_PASS)

docker-build:
	docker build -t $(REGISTRY_REPO):$(VERSION) .

docker-push: docker-login
	docker push $(REGISTRY_REPO):$(VERSION)

deploy:
	ssh -o "StrictHostKeyChecking no" $(HOST_USER)@$(HOST) make run-music-$(VERSION)

lint-all l:
	bash ./scripts/golangci-lint.sh
	bash ./scripts/consistent.sh

rigo:
	make install & make rgo
