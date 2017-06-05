#
# ABC Dockerfile
# docker build -t abc .
# docker run --rm --name abc abc
# docker run --rm --name abc abc login google   
#

# Pull the base image
FROM golang:1.8
MAINTAINER Avi Aryan <avi.aryan123@gmail.com>

# Set GOPATH
ENV GOPATH /go

# Make directories for api_frontend
RUN mkdir -p /go/src/github.com/appbaseio/abc

# Add api_frontend files
ADD . /go/src/github.com/appbaseio/abc

# Define working directory
WORKDIR /go/src/github.com/appbaseio/abc

RUN cd /go/src/github.com/appbaseio/abc && \
	go build ./cmd/abc/...

# Define default entrypoint
# Entrypoint gets extra parameters from docker run
ENTRYPOINT ["./abc"]
