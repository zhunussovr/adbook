### Go Build stage
FROM golang:1.23-alpine AS builder

RUN apk update && apk add --no-cache git ca-certificates

WORKDIR /go/src/app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
ARG BINARY_NAME=app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/${BINARY_NAME} ./cmd/adbook/main.go

### Image Build stage
FROM alpine:3.20

# Install dependencies and create user
RUN apk update && \
    apk add --no-cache wget ca-certificates && \
    adduser -D -g '' localuser && \
    addgroup localuser wheel

# Set up working directory
WORKDIR /app

# Copy binary and config files
ARG BINARY_NAME=app
COPY --from=builder /go/bin/${BINARY_NAME} ./app
COPY --from=builder /go/src/app/config.toml ./
COPY --from=builder /go/src/app/web ./web/

# Change ownership to non-root user
RUN chown -R localuser:localuser /app

# Switch to non-root user
USER localuser

EXPOSE 8080

# Run the application
CMD ["./app"]