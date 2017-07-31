// Package statsdLogger provides a simple dummy StatsD logging server for local development
//
// Adapted from http://lee.hambley.name/2013/01/26/dirt-simple-statsd-server-for-local-development.html
package statsdLogger

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/fatih/color"
)

// DefaultAddress to listen for metrics on
var DefaultAddress = ":8125"

// DefaultOutput where to output the logs
var DefaultOutput io.Writer = os.Stdout

// Server listens for statsd metrics, and logs them to the console for development
type Server interface {
	// Listen starts server listening on provided UDP address
	Listen() error

	// Close stops the server from listening for metrics
	Close() error
}

// New returns a local statsd logging server which logs to statsdLogger.DefaultOutput
func New(address string) (Server, error) {
	return NewWithWriter(address, DefaultOutput)
}

// NewWithWriter returns a local statsd logging server which logs to provided output
func NewWithWriter(address string, output io.Writer) (Server, error) {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, errors.New("[StatsD] Invalid address")
	}

	conn, _ := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, errors.New("[StatsD] Unable to listen to udp stream")
	}

	server := server{address: address, connection: conn, output: output, closed: false}

	return &server, nil
}

type server struct {
	address    string
	connection *net.UDPConn
	closed     bool
	output     io.Writer
}

func (l *server) Listen() error {
	if l.closed {
		return errors.New("[StatsD] Server already closed")
	}

	fmt.Fprintf(l.output, "[StatsD] Listening at %s\n", l.address)

	buffer := make([]byte, 1024)

	for !l.closed {
		_, err := l.connection.Read(buffer)

		if err != nil {
			fmt.Printf("[StatsD] Read error %s: \n", err)
		}

		metric := parseMetric(string(buffer))
		l.logMetric(metric)
	}

	return nil
}

func (l *server) Close() error {
	l.closed = true
	fmt.Fprintf(l.output, "[StatsD] Shutting down\n")
	l.connection.Close()
	return nil
}

func (l *server) logMetric(metric metric) {
	fmt.Fprintf(
		l.output,
		"[StatsD] %s %s %s\n",
		color.BlueString(metric.name),
		color.YellowString(metric.value),
		color.CyanString(metric.tags))
}
