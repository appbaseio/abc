#
# abc Dockerfile
# docker build --build-arg ABC_BUILD=oss -t abc .
# private: docker build --build-arg ABC_BUILD=noss -t abc .
# docker volume create --name abc
# docker run -i --rm -v abc:/root abc login google
# root is $HOME, -i for stdin, --rm to remove container
#

# Pull the base image
FROM golang:1.16-alpine as builder
MAINTAINER Siddharth Kothari <siddharth@appbase.io>

WORKDIR /app

COPY . .

RUN go build -tags "oss" -o abc ./cmd/abc/...

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /abc/
COPY --from=builder /app/abc abc
EXPOSE 8080

# Define default entrypoint
# Entrypoint gets extra parameters from docker run
ENTRYPOINT ["/abc/abc"]
