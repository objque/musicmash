FROM golang:latest as builder

ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64

ARG RELEASE=unset
ARG COMMIT=unset
ARG BUILD_TIME=unset
ENV PROJECT=/go/src/github.com/musicmash/musicmash/internal

WORKDIR /go/src/github.com/musicmash/musicmash
COPY cmd cmd
COPY internal internal
COPY pkg pkg
COPY vendor vendor

RUN go build -v -a \
    -mod=vendor \
    -installsuffix cgo \
    -gcflags "all=-trimpath=$(GOPATH)" \
    -ldflags '-linkmode external -extldflags "-static" -s -w \
       -X ${PROJECT}/version.Release=${RELEASE} \
       -X ${PROJECT}/version.Commit=${COMMIT} \
       -X ${PROJECT}/version.BuildTime=${BUILD_TIME}"' \
    -o /usr/local/bin/musicmash ./cmd/...

FROM alpine:latest

RUN addgroup -S musicmash && adduser -S musicmash -G musicmash
USER musicmash
WORKDIR /home/musicmash

COPY --from=builder --chown=musicmash:musicmash /usr/local/bin/musicmash /usr/local/bin/musicmash

ENTRYPOINT ["musicmash"]
