package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/catkins/statsd-logger/metrics"
	"github.com/catkins/statsd-logger/trace"
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

	traceServer := trace.NewServer(trace.DefaultAddress)

	go func() {
		traceServer.Listen()
	}()

	<-shutdownChan
	metricsServer.Close()
	traceServer.Close()
}
