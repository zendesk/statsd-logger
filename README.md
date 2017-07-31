## StatsD Logger for Go

Simple logger for StatsD metrics, adapted from http://lee.hambley.name/2013/01/26/dirt-simple-statsd-server-for-local-development.html to make it easier to debug metrics in development.

## Usage

### CLI

```bash
go get -u github.com/catkins/statsd-logger/cmd/statsd-logger
statsd-logger
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

## Licence

The MIT License

Copyright 2017 Chris Atkins
