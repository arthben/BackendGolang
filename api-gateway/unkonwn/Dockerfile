FROM golang:1.23-alpine

WORKDIR /build
COPY . .
RUN go mod download
RUN go get -u github.com/swaggo/swag
RUN go build -o ./api-gateway ./cmd/api-gateway/main.go
