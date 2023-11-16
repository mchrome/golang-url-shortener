FROM golang:1.21.1-alpine

ENV GOPATH=/

COPY . .

RUN go mod download
RUN go build -v -o url-compression-api ./cmd/main.go

EXPOSE 8000
RUN adduser --disabled-password default-user
USER default-user