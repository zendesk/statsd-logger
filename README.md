## StatsD Logger for Go

[![Build Status](https://travis-ci.org/catkins/statsd-logger.svg?branch=master)](https://travis-ci.org/catkins/statsd-logger) [![GoDoc](https://godoc.org/github.com/catkins/statsd-logger?status.svg)](https://godoc.org/github.com/catkins/statsd-logger) [![Go Report Card](https://goreportcard.com/badge/github.com/catkins/statsd-logger)](https://goreportcard.com/report/github.com/catkins/statsd-logger)
[![GitHub tag](https://img.shields.io/github/tag/catkins/statsd-logger.svg)]()
[![Docker Automated build](https://img.shields.io/docker/automated/catkins/statsd-logger.svg)]()
[![ImageLayers Size](https://img.shields.io/imagelayers/image-size/catkins/statsd-logger/latest.svg)]()

Simple logger for StatsD metrics, adapted from http://lee.hambley.name/2013/01/26/dirt-simple-statsd-server-for-local-development.html to make it easier to debug metrics in development. Beyond converting it from Ruby to Go, also add colour output and rendering of DogStatsd tags.

## Usage

### CLI

```bash
go get -u github.com/catkins/statsd-logger/cmd/statsd-logger
statsd-logger

# send it some metrics using a library or (low-tech) netcat
echo -n "my.awesome_counter:1|c#cool:tags,another_tag:with_value" | nc -u -u -w0 localhost 8125
```

### Library

```bash
go get -u github.com/catkins/statsd-logger
```

Embed it into an existing application

```go
package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/catkins/statsd-logger"
)

func main() {
	shutdownChan := make(chan os.Signal)
	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	server, err := statsdLogger.New("0.0.0.0:8125")
	if err != nil {
		panic(err)
	}

	go func() {
		server.Listen()
	}()

	<-shutdownChan
	server.Close()
}
```

## Docker

This can also log StatsD metrics beibng being sent to udp:8125 on running `docker` containers.

```sh
# log metrics for container named "myapp"
docker run --rm -it --net="container:myapp" catkins/statsd-logger
```

## Licence

The MIT License

Copyright 2017 Chris Atkins
