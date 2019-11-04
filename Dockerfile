FROM golang:latest as builder

WORKDIR /go/src/github.com/musicmash/musicmash
COPY . .

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
RUN go build -v -a -installsuffix cgo -gcflags "all=-trimpath=$(GOPATH)" -o bin/musicmash ./cmd/...

FROM alpine:latest

RUN apk update && apk upgrade && \
    apk add --no-cache \
    ca-certificates vim curl && \
    rm -rf /var/cache/apk/*

WORKDIR /root/
COPY --from=builder /go/src/github.com/musicmash/musicmash/bin .

ENTRYPOINT ["./musicmash"]
