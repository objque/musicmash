FROM golang:latest as builder

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /go/src/github.com/musicmash/musicmash
COPY cmd cmd
COPY internal internal
COPY pkg pkg
COPY vendor vendor

RUN go build -v -a -installsuffix cgo -gcflags "all=-trimpath=$(GOPATH)" -o bin/musicmash ./cmd/...

FROM alpine:latest

RUN apk update && apk upgrade && \
    apk add --no-cache \
    ca-certificates vim curl && \
    rm -rf /var/cache/apk/*

WORKDIR /root/
COPY --from=builder /go/src/github.com/musicmash/musicmash/bin .

ENTRYPOINT ["./musicmash"]
