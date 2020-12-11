FROM golang:1-alpine as builder

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

ARG RELEASE=unset
ARG COMMIT=unset
ARG BUILD_TIME=unset
ENV PROJECT=github.com/musicmash/musicmash

WORKDIR /go/src/github.com/musicmash
COPY migrations /var/musicmash/migrations
COPY go.mod go.mod
COPY go.sum go.sum
COPY cmd cmd
COPY pkg pkg
COPY internal internal

RUN go build -v \
    -gcflags "all=-trimpath=${WORKDIR}" \
    -ldflags "-w -s \
       -X ${PROJECT}/internal/version.Release=${RELEASE} \
       -X ${PROJECT}/internal/version.Commit=${COMMIT} \
       -X ${PROJECT}/internal/version.BuildTime=${BUILD_TIME}" \
    -o /usr/local/bin/musicmash ./cmd/musicmash/...

RUN go build -v \
    -gcflags "all=-trimpath=${WORKDIR}" \
    -ldflags "-w -s \
       -X ${PROJECT}/internal/version.Release=${RELEASE} \
       -X ${PROJECT}/internal/version.Commit=${COMMIT} \
       -X ${PROJECT}/internal/version.BuildTime=${BUILD_TIME}" \
    -o /usr/local/bin/musicmashctl ./cmd/musicmashctl/...

FROM alpine:latest

RUN addgroup -S musicmash && adduser -S musicmash -G musicmash
USER musicmash
WORKDIR /home/musicmash

COPY --from=builder --chown=musicmash:musicmash /var/musicmash/migrations /var/musicmash/migrations
COPY --from=builder --chown=musicmash:musicmash /usr/local/bin/musicmash /usr/local/bin/musicmash
COPY --from=builder --chown=musicmash:musicmash /usr/local/bin/musicmashctl /usr/local/bin/musicmashctl

ENTRYPOINT ["/usr/local/bin/musicmash"]
CMD ["-db-auto-migrate=true", "-db-migrations-dir=/var/musicmash/migrations"]
