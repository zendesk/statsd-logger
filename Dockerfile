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

ARG BUILD_DATE
ARG VCS_REF

# https://medium.com/microscaling-systems/labelling-automated-builds-on-docker-hub-f3d073fb8e1
LABEL org.label-schema.build-date=$BUILD_DATE \
  org.label-schema.vcs-url="https://github.com/catkins/statsd-logger.git" \
  org.label-schema.vcs-ref=$VCS_REF \
  org.label-schema.schema-version="1.0.0-rc1"

RUN mkdir /app
WORKDIR /app

COPY --from=binary-builder /go/bin/statsd-logger .

CMD ["/app/statsd-logger"]
