# Dockerfile References: https://docs.docker.com/engine/reference/builder/
# Start from the latest golang base image
FROM golang:1.20.2-bullseye as build-env

# Set the Current Working Directory inside the container
WORKDIR /app


# Copy the source from the current directory to the Working Directory inside the container
COPY . .

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

# RUN go mod download

# # Build the Go app
# RUN  go build -o go-pentor-bank.bin .


FROM alpine:3.10

ARG git_commit=default
ARG version=0.0.0
ARG env=dev

ENV GIT_COMMIT=$git_commit
ENV VERSION=$version
ENV ENV=$env
ENV TZ=Asia/Bangkok
LABEL GIT_COMMIT=$git_commit \
      VERSION=$version \
      ENV=$env
LABEL vendor="TODO" project="go-pentor-bank"
RUN apk add --no-cache tzdata

WORKDIR /app

RUN ls -la /app

RUN mkdir configs/
RUN mkdir docs/
# RUN mkdir templates/

# Add configs so go.viper can read config.
COPY configs ./configs
COPY docs ./docs
COPY deployments/docker/assets/zoneinfo.zip /


COPY  go-pentor-bank.bin ./go-pentor-bank.bin
COPY  csr ./csr
COPY  docker-entrypoint.sh ./docker-entrypoint.sh

# Add zoneinfo.zip so time.LoadLocation can work inside alpine image.

# COPY docs ./docs

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
ENTRYPOINT [ "/bin/sh", "docker-entrypoint.sh" ]

