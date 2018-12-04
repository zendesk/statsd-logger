## StatsD Logger for Go

[![Build Status](https://travis-ci.org/catkins/statsd-logger.svg?branch=master)](https://travis-ci.org/catkins/statsd-logger) [![GoDoc](https://godoc.org/github.com/catkins/statsd-logger?status.svg)](https://godoc.org/github.com/catkins/statsd-logger) [![Go Report Card](https://goreportcard.com/badge/github.com/catkins/statsd-logger)](https://goreportcard.com/report/github.com/catkins/statsd-logger)
[![GitHub tag](https://img.shields.io/github/tag/catkins/statsd-logger.svg)]()
[![Docker Automated build](https://img.shields.io/docker/automated/catkins/statsd-logger.svg)]()
[![](https://images.microbadger.com/badges/image/catkins/statsd-logger.svg)](https://microbadger.com/images/catkins/statsd-logger "Get your own image badge on microbadger.com")

Simple logger for StatsD metrics, adapted from http://lee.hambley.name/2013/01/26/dirt-simple-statsd-server-for-local-development.html to make it easier to debug metrics in development. Beyond converting it from Ruby to Go, also add colour output and rendering of DogStatsd tags.

## Usage

### CLI

```bash
go get -u github.com/catkins/statsd-logger/cmd/statsd-logger
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

Copyright 2017 Chris Atkins
