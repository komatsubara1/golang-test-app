FROM golang:1.20-bullseye

ENV GO111MODULE=
ENV GOPATH=/go

RUN apt-get update && apt-get install -y git

COPY ./ /go

WORKDIR /go/app

RUN go mod tidy

RUN go run main.go
