FROM golang:1-alpine3.18 as builder

RUN mkdir -p /src
RUN apk add git curl

WORKDIR /src
COPY ../.. .

RUN curl -L -o grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v0.3.1/grpc_health_probe-linux-amd64 && chmod +x grpc_health_probe && mv grpc_health_probe /

RUN CGO_ENABLED=0  GOOS=linux GOARCH=amd64 go build -a -tags netgo -o app .

# second stage
FROM scratch
ADD res /res
COPY --from=0 /src/app /app
COPY --from=0 /grpc_health_probe /bin/grpc_health_probe
CMD ["/app"]