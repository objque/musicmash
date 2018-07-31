all:

clean:
	rm bin/musicmash || true

build: clean
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a -installsuffix cgo -o bin/musicmash cmd/musicmash.go

prepare-tests:
	go get -u github.com/kyoh86/richgo

install:
	go install -v ./cmd

t tests: install
	go test -v ./internal/...

add-ssh-key:
	openssl aes-256-cbc -K $(encrypted_a4311917bb34_key) -iv $(encrypted_a4311917bb34_iv) -in travis_key.enc -out /tmp/travis_key -d
	chmod 600 /tmp/travis_key
	ssh-add /tmp/travis_key

docker-build:
	docker build -t $(REGISTRY_REPO) .

docker-login:
	docker login -u $(REGISTRY_USER) -p $(REGISTRY_PASS)

docker-push: docker-login
	docker push $(REGISTRY_REPO)

docker-full:
	make docker-build
	make docker-push

deploy:
	ssh -o "StrictHostKeyChecking no" $(HOST_USER)@$(HOST) "(docker stop $(CONTAINER_NAME) || true) && (docker rm $(CONTAINER_NAME) || true)"
	ssh -o "StrictHostKeyChecking no" $(HOST_USER)@$(HOST) docker pull $(REGISTRY_REPO)
	ssh -o "StrictHostKeyChecking no" $(HOST_USER)@$(HOST) docker run -d --link mariadb -e TG_TOKEN=$(TG_TOKEN) --name $(CONTAINER_NAME) -v /etc/musicmash:/etc/musicmash -v /etc/ssl:/etc/ssl $(REGISTRY_REPO)
