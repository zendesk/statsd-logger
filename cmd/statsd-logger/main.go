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
