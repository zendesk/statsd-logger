FROM golang:1.9-alpine

WORKDIR /go/src/github.com/catkins/statsd-logger

COPY . .

RUN cd cmd/statsd-logger && go install

CMD ["statsd-logger"]
