FROM golang:alpine

RUN go mod download

RUN go test -bench . -benchmem ./...
