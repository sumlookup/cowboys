FROM golang:1-alpine3.18 as builder

RUN mkdir -p /src
RUN apk add git curl

WORKDIR /src
COPY ../.. .

ENV GOPRIVATE=github.com/sumlookup
ENV GO111MODULE=on

RUN curl -L -o grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v0.3.1/grpc_health_probe-linux-amd64 && chmod +x grpc_health_probe && mv grpc_health_probe /

ADD ../../res /res

CMD go run .