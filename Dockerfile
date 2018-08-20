#
# abc Dockerfile
# docker build --build-arg ABC_BUILD=oss -t abc .
# private: docker build --build-arg ABC_BUILD=noss -t abc .
# docker volume create --name abc
# docker run -i --rm -v abc:/root abc login google
# root is $HOME, -i for stdin, --rm to remove container
#

# Pull the base image
FROM alpine:3.8
MAINTAINER Siddharth Kothari <siddharth@appbase.io>

# Set GOPATH
ENV GOPATH /go

# certs
RUN apk --update add --no-cache ca-certificates
RUN update-ca-certificates

# Get build variant
ARG ABC_BUILD=oss
ENV ABC_BUILD ${ABC_BUILD}

# Make directories for the code
RUN mkdir -p /go/src/github.com/appbaseio/abc
RUN mkdir -p /abc

# Add abc files
ADD . /go/src/github.com/appbaseio/abc

# install
RUN apk add --no-cache --virtual build-dependencies go libc-dev && \
	cd /go/src/github.com/appbaseio/abc && \
	go build -tags $ABC_BUILD ./cmd/abc/... && \
	mv ./abc /abc/ && \
	apk del build-dependencies && \
	rm -rf /go && \
	rm -rf /usr/local/go

WORKDIR /abc

# Define default entrypoint
# Entrypoint gets extra parameters from docker run
ENTRYPOINT ["./abc"]
