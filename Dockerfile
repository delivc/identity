FROM golang:1.13-alpine as build
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux

RUN apk add --no-cache make git

WORKDIR /go/src/github.com/delivc/identity
COPY . /go/src/github.com/delivc/identity

RUN make deps build

FROM alpine:3.7
RUN adduser -D -u 1000 delivc

RUN apk add --no-cache ca-certificates
COPY --from=build /go/src/github.com/delivc/identity /usr/local/bin/identity
COPY --from=build /go/src/github.com/delivc/identity/migrations /usr/local/etc/identity/migrations/

ENV GOTRUE_DB_MIGRATIONS_PATH /usr/local/etc/identity/migrations

USER delivc
CMD ["oauth"]