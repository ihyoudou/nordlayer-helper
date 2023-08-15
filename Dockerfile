FROM golang:1.18.10-bullseye

WORKDIR /app

COPY . .

RUN apt-get update && apt-get install -y make gcc libgtk-3-dev libayatana-appindicator3-dev protobuf-compiler

RUN go mod download && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

RUN make build

