FROM golang:1.3
MAINTAINER TANABE Ken-ichi <nabeken@tknetworks.org>

RUN mkdir -p /go/src/app
WORKDIR /go/src/app

COPY . /go/src/app
RUN go-wrapper download ./...
RUN go install ./...
