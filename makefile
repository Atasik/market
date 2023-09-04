.PHONY:
.SILENT:

build:
    go mod download && CGO_ENABLED=0 GOOS=linux go build -o ./.bin/market ./cmd/market/main.go 

run: build
    docker compose up market

test:
    go test -v ./...

swag:
    swag init -g internal/app/app.go