### Go Build stage

FROM golang:1.16.4-alpine3.13 AS builder

RUN apk update && apk add --no-cache git

WORKDIR /go/src/app

COPY . .

RUN go get -d -v && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/app ./cmd/adbook/main.go

### Image Build stage

FROM alpine:3.13

RUN apk update && \
    adduser -D -g '' localUser1 && \
    addgroup localUser1 wheel \
USER localUser1

WORKDIR /go/bin

COPY --from=builder /go/bin/app .
COPY --from=builder /go/src/app/config.toml .
COPY --from=builder /go/src/app/web web/

EXPOSE 8080/tcp

HEALTHCHECK CMD wget -q -O /dev/null http://localhost:8080/health || exit 1

ENTRYPOINT ["./app"]
