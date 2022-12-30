# syntax=docker/dockerfile:1

## Build
FROM golang:1.18-alpine AS build

WORKDIR /app

# if you build it in China, add this
#ENV GOPROXY=https://goproxy.cn,direct

COPY ./go.mod ./
COPY ./go.sum ./
COPY ./Makefile ./
RUN go mod download && \
    go install github.com/swaggo/swag/cmd/swag@v1.8.7 && \
    apk add --no-cache --update make gcc g++

COPY . .
RUN make

## Deploy
FROM alpine:3

RUN apk add --no-cache libstdc++

RUN adduser -D sibyl
USER sibyl
WORKDIR /home/sibyl

COPY --from=build /app/sibyl /home/sibyl/
EXPOSE 9876
ENV GIN_MODE=release

RUN mkdir sibyl2-badger-storage
VOLUME /home/sibyl/sibyl2-badger-storage

ENTRYPOINT ["./sibyl", "server"]
