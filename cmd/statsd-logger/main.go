package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/catkins/statsd-logger/metrics"
)

func main() {
	shutdownChan := make(chan os.Signal)
	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	metricsServer, err := metrics.NewServer(metrics.DefaultAddress)
	if err != nil {
		panic(err)
	}

	go func() {
		metricsServer.Listen()
	}()

	<-shutdownChan
	metricsServer.Close()
}
