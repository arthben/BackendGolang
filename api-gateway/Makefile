build:
	go build -o api-gateway ./cmd/api-gateway

swag:
	swag init -g ./cmd/api-gateway/main.go --markdownFiles swagger-markdown --parseDependency true

run:
	go run cmd/api-gateway/main.go

.PHONY: swag run
