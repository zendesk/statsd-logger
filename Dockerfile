FROM golang:1.17 as binary-builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

COPY . .

ENV CGO_ENABLED=0
RUN go build -o statsd-logger cmd/statsd-logger/main.go

# use multistage builds for smaller final image
# https://docs.docker.com/engine/userguide/eng-image/multistage-build/#use-multi-stage-builds
FROM alpine:latest

EXPOSE 8125/udp
EXPOSE 8126

RUN mkdir /app
WORKDIR /app

COPY --from=binary-builder /app/statsd-logger .

CMD ["/app/statsd-logger"]
