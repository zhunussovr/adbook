// ...existing code...
# "adbook" - Lightweight Active Directory people finder tool

Small project to write lightweight web service which you can deploy and run it with few clicks.
With this service anyone in your company(with respective permissions) will be able to search people profiles and make small reports.

This is only starting point and more details will be available soon.

## Small information about techologies
Service will be written on GoLang and will use MongoDB as a data storage.

## Instructions

Quick, minimal steps to prepare, build and run this project on macOS.

Prerequisites
(Mac OS)
- Install Homebrew if needed: `/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"`
- Install Go and Docker (Docker Desktop includes docker and compose):
  - `brew install go`
  - `brew install --cask docker` (then open Docker.app once)
  - `brew install docker-compose`
- Verify:
  - `go version`
  - `docker --version`
  - `docker-compose --version`

Local build & run
1. Fetch deps:
   - `go mod download`
2. Format & vet:
   - `gofmt -w .`
   - `go vet ./...`
3. Run tests:
   - `go test ./... -v -race -coverprofile=coverage.out`
   - `go tool cover -html=coverage.out` (view coverage)
4. Build binary:
   - `go build -o ./bin/adbook ./cmd/adbook`
5. Run:
   - `./bin/adbook --conf=config.toml`
   - Check logs and endpoints (default port 8080 unless changed in config)

Docker
- Build image:
  - `docker build -t adbook:local .`
- Run container:
  - `docker run --rm -p 8080:8080 adbook:local`
- Or use docker-compose (if provided):
  - `docker compose up --build`

Troubleshooting
- If Docker commands fail, ensure Docker Desktop is running.
- If tests fail, run the failing package with `go test ./pkg/name -v` to inspect errors.
- Check `config.toml` for required connection strings (LDAP/Mongo) and credentials.

If you want, I can:
- add a small bootstrap script (bootstrap.sh) to automate install/build steps,
- add a health endpoint and graceful shutdown,
- create a GitHub Actions CI workflow to run these checks automatically.

// ...existing code...