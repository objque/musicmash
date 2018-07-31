all:

clean:
	rm bin/musicmash || true

build: clean
	CGO_ENABLED=0 GOOS=linux go build -v -a -installsuffix cgo -o bin/musicmash cmd/musicmash.go

prepare-tests:
	go get -u github.com/kyoh86/richgo

t tests:
	richgo test -v ./internal/...
