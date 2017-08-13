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

// New returns a local statsd logging server which logs to statsdLogger.DefaultOutput and is formatted with statsdLogger.DefaultFormatter
func New(address string, options ...func(*server)) (Server, error) {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, errors.New("[StatsD] Invalid address")
	}

	conn, _ := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, fmt.Errorf("[StatsD] Unable to listen at udp: %s", address)
	}

	server := server{
		port:       addr.Port,
		connection: conn,
		output:     DefaultOutput,
		closed:     false,
		formatter:  DefaultFormatter,
	}

	for _, optionFunc := range options {
		optionFunc(&server)
	}

	return &server, nil
}

// WithWriter is is provided as an option to New to specify a custom io.Writer to output logs
// usage:
//	statsdLogger.New("0.0.0.0:8125", WithWriter(os.Stderr))
func WithWriter(output io.Writer) func(*server) {
	return func(server *server) {
		server.output = output
	}
}

// WithFormatter is is provided as an option to New to specify a custom formatter
// usage:
//
//	statsdLogger.New("0.0.0.0:8125", WithFormatter(myCustomFormatter))
func WithFormatter(formatter MetricFormatter) func(*server) {
	return func(server *server) {
		server.formatter = formatter
	}
}

type server struct {
	port       int
	connection *net.UDPConn
	closed     bool
	output     io.Writer
	formatter  MetricFormatter
}

func (s *server) Listen() error {
	if s.closed {
		return errors.New("[StatsD] Server already closed")
	}

	fmt.Fprintf(s.output, "[StatsD] Listening on port %d\n", s.port)

	buffer := make([]byte, 1024)

	for !s.closed {
		numBytes, err := s.connection.Read(buffer)

		// blocking read returns an error when connection is closed
		if err != nil && s.closed {
			break

		}

		if err != nil {
			fmt.Printf("[StatsD] Read error %s: \n", err)
			continue
		}

		rawMetric := buffer[0:numBytes]
		metric := ParseMetric(string(rawMetric))

		s.logMetric(metric)
	}

	return nil
}

func (s *server) Close() error {
	s.closed = true
	fmt.Fprintf(s.output, "[StatsD] Shutting down\n")

	return s.connection.Close()
}

func (s *server) logMetric(metric Metric) {
	fmt.Fprintf(s.output, s.formatter.Format(metric))
}
