### Go Build stage

FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git

WORKDIR /go/src/app

COPY . .

RUN go get -d -v

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/app ./cmd/adbook/main.go

### Image Build stage

FROM alpine

RUN apk update && apk add --no-cache ca-certificates

COPY --from=builder /go/bin/app /go/bin/app
COPY --from=builder /go/src/app/config.toml /go/bin/
COPY --from=builder /go/src/app/web /go/bin/web/

WORKDIR /go/bin

EXPOSE 8080/tcp

ENTRYPOINT ["./app"]