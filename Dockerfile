FROM golang:1.9 as binary-builder

WORKDIR /go/src/github.com/catkins/statsd-logger

# install dep
RUN curl -L https://github.com/golang/dep/releases/download/v0.3.2/dep-linux-amd64 > $GOPATH/bin/dep \
  && chmod +x $GOPATH/bin/dep

# install runtime dependencies
COPY Gopkg.lock .
COPY Gopkg.toml .
RUN dep ensure -v --vendor-only

COPY . .

ENV CGO_ENABLED=0
RUN cd cmd/statsd-logger && go install

# use multistage builds for smaller final image
# https://docs.docker.com/engine/userguide/eng-image/multistage-build/#use-multi-stage-builds
FROM alpine:latest

RUN mkdir /app
WORKDIR /app

COPY --from=binary-builder /go/bin/statsd-logger .

CMD ["/app/statsd-logger"]
