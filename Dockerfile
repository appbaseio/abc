#
# abc Dockerfile
# docker build --build-arg ABC_BUILD=oss -t abc .
# private: docker build --build-arg ABC_BUILD=noss -t abc .
# docker volume create --name abc
# docker run -i --rm -v abc:/root abc login google
# root is $HOME, -i for stdin, --rm to remove container
#

# Pull the base image
FROM golang:1.12.5 AS builder
MAINTAINER Siddharth Kothari <siddharth@appbase.io>

# Get build variant
ARG ABC_BUILD=oss
ENV ABC_BUILD ${ABC_BUILD}

RUN apt-get update && \
	apt-get install -y libssl-dev && \
	mkdir -p $GOPATH/github.com/src/appbaseio/abc && \
	mkdir -p /abc && \
	curl -LO https://github.com/neo4j-drivers/seabolt/releases/download/v1.7.4/seabolt-1.7.4-Linux-ubuntu-18.04.deb && \
	dpkg -i seabolt-1.7.4-Linux-ubuntu-18.04.deb && \
	go get github.com/neo4j/neo4j-go-driver/neo4j && \
	go get gopkg.in/olivere/elastic.v7

WORKDIR $GOPATH/src/github.com/appbaseio/abc

COPY . .

RUN go build -tags "seabolt_static $ABC_BUILD" -o /abc/abc ./cmd/abc/...

FROM ubuntu:bionic
MAINTAINER Siddharth Kothari <siddharth@appbase.io>

# certs
RUN apt-get update && \
	apt-get install -y ca-certificates && \
	update-ca-certificates

COPY --from=builder /abc/abc /abc/abc

# Define default entrypoint
# Entrypoint gets extra parameters from docker run
ENTRYPOINT ["/abc/abc"]
