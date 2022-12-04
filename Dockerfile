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
    apk add --update make gcc g++

COPY . .
RUN make

## Deploy
FROM alpine

RUN apk add --no-cache libstdc++

WORKDIR /

COPY --from=build /app/sibyl /app/sibyl

EXPOSE 9876

RUN adduser -D sibyl
USER sibyl
ENV GIN_MODE=release

ENTRYPOINT ["/app/sibyl", "server"]
