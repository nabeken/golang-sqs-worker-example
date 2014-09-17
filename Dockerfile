FROM google/golang
MAINTAINER TANABE Ken-ichi <nabeken@tknetworks.org>

WORKDIR /gopath/src/github.com/nabeken/golang-sqs-worker-example
ADD . /gopath/src/github.com/nabeken/golang-sqs-worker-example
RUN go get -v ./...
