FROM golang:1.16 as binary-builder

WORKDIR /statsd-logger
COPY . .

ENV CGO_ENABLED=0
RUN cd cmd/statsd-logger && go install

# use multistage builds for smaller final image
# https://docs.docker.com/engine/userguide/eng-image/multistage-build/#use-multi-stage-builds
FROM alpine:latest

# https://medium.com/microscaling-systems/labelling-automated-builds-on-docker-hub-f3d073fb8e1
ARG BUILD_DATE
ARG VCS_REF
LABEL org.label-schema.build-date=$BUILD_DATE \
  org.label-schema.vcs-url="https://github.com/catkins/statsd-logger.git" \
  org.label-schema.vcs-ref=$VCS_REF \
  org.label-schema.schema-version="1.0.0-rc1"

EXPOSE 8125/udp

RUN mkdir /app
WORKDIR /app

COPY --from=binary-builder /go/bin/statsd-logger .

CMD ["/app/statsd-logger"]
