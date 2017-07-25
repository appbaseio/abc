#
# abc Dockerfile
# docker build --build-arg ABC_BUILD=oss -t abc .
# docker volume create --name abc
# docker run -i --rm -v abc:/root abc login google
# root is $HOME, -i for stdin, --rm to remove container
#

# Pull the base image
FROM golang:1.8-alpine 
MAINTAINER Avi Aryan <avi.aryan123@gmail.com>

# Set GOPATH
ENV GOPATH /go

# Make directories for the code
RUN mkdir -p /go/src/github.com/appbaseio/abc

# Add abc files
ADD . /go/src/github.com/appbaseio/abc

# Define working directory
WORKDIR /go/src/github.com/appbaseio/abc

# Get build variant
ARG ABC_BUILD=oss
ENV ABC_BUILD ${ABC_BUILD}

# Run build
RUN cd /go/src/github.com/appbaseio/abc && \
	go build -tags $ABC_BUILD ./cmd/abc/...

# Define default entrypoint
# Entrypoint gets extra parameters from docker run
ENTRYPOINT ["./abc"]
