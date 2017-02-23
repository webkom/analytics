FROM golang:1.8-alpine
MAINTAINER Abakus Webkom <webkom@abakus.no>

RUN mkdir -p /go/src/app
WORKDIR /go/src/app
COPY . /go/src/app

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

RUN set -e \
    && apk update \
    && apk upgrade \
    && apk add --no-cache bash git openssh \
    && go get golang.org/x/sys/unix \
    && go get github.com/tools/godep \
    && godep get \
    && go build

CMD ["go-wrapper", "run"]
