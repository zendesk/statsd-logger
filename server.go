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
func New(address string) (Server, error) {
	return NewWithWriterAndFormatter(address, DefaultOutput, DefaultFormatter)
}

// NewWithWriter returns a local statsd logging server which logs to provided output
func NewWithWriter(address string, output io.Writer) (Server, error) {
	return NewWithWriterAndFormatter(address, output, DefaultFormatter)
}

// NewWithFormatter returns a local statsd logging server which logs with the provided formatter to statsdLogger.DefaultFormatter
func NewWithFormatter(address string, formatter MetricFormatter) (Server, error) {
	return NewWithWriterAndFormatter(address, DefaultOutput, formatter)
}

// NewWithWriterAndFormatter returns a local statsd logging server which logs to provided output and formats output with provided formatter
func NewWithWriterAndFormatter(address string, output io.Writer, formatter MetricFormatter) (Server, error) {
	addr, err := net.ResolveUDPAddr("udp", address)
	port := addr.Port

	if err != nil {
		return nil, errors.New("[StatsD] Invalid address")
	}

	conn, _ := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, errors.New("[StatsD] Unable to listen to udp stream")
	}

	server := server{
		port:       port,
		connection: conn,
		output:     output,
		closed:     false,
		formatter:  formatter,
	}

	return &server, nil
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
