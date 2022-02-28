## StatsD Logger for Go

[![GoDoc](https://godoc.org/github.com/zendesk/statsd-logger?status.svg)](https://godoc.org/github.com/zendesk/statsd-logger) [![Go Report Card](https://goreportcard.com/badge/github.com/zendesk/statsd-logger)](https://goreportcard.com/report/github.com/zendesk/statsd-logger)
[![GitHub tag](https://img.shields.io/github/tag/catkins/statsd-logger.svg)]()

Simple logger for StatsD metrics, adapted from http://lee.hambley.name/2013/01/26/dirt-simple-statsd-server-for-local-development.html to make it easier to debug metrics in development. Beyond converting it from Ruby to Go, also adds colour output and rendering of DogStatsd tags.

It will also listen for Datadog APM traces and log them out.

## Usage

### CLI

```bash
go get -u github.com/zendesk/statsd-logger/cmd/statsd-logger
statsd-logger

# send it some metrics using a library or (low-tech) netcat
echo -n "my.awesome_counter:1|c#cool:tags,another_tag:with_value" | nc -u -w0 localhost 8125
```

## Docker

`statsd-logger` can also log StatsD metrics being being sent to udp:8125 on running `docker` containers without modification.

```sh
# log metrics for container named "myapp"
docker run --rm -it --net="container:myapp" catkins/statsd-logger
```

## Licence

The MIT License

Copyright 2017 Zendesk
