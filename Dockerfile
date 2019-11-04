FROM golang:latest as builder

ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /go/src/github.com/musicmash/musicmash
COPY cmd cmd
COPY internal internal
COPY pkg pkg
COPY vendor vendor

RUN go build -v -a \
    -installsuffix cgo \
    -gcflags "all=-trimpath=$(GOPATH)" \
    -ldflags '-linkmode external -extldflags "-static"' \
    -o /usr/local/bin/musicmash ./cmd/...

FROM alpine:latest

RUN addgroup -S musicmash && adduser -S musicmash -G musicmash
USER musicmash
WORKDIR /home/musicmash

COPY --from=builder --chown=musicmash:musicmash /usr/local/bin/musicmash /usr/local/bin/musicmash

ENTRYPOINT ["musicmash"]
