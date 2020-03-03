// Package metrics provides a simple dummy StatsD logging server for local development
//
// Adapted from http://lee.hambley.name/2013/01/26/dirt-simple-statsd-server-for-local-development.html
package metrics

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"sync/atomic"
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

// NewServer returns a local statsd logging server which logs to statsdLogger.DefaultOutput and is formatted with statsdLogger.DefaultFormatter
func NewServer(address string, options ...func(*server)) (Server, error) {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, errors.New("[StatsD] Invalid address")
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, fmt.Errorf("[StatsD] Unable to listen at udp: %s", address)
	}

	server := server{
		port:       addr.Port,
		connection: conn,
		output:     DefaultOutput,
		formatter:  DefaultFormatter,
		outputLock: &sync.Mutex{},
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
	closed     int32
	output     io.Writer
	formatter  MetricFormatter
	outputLock *sync.Mutex
}

func (s *server) Listen() error {
	if s.isClosed() {
		return errors.New("[StatsD] Server already closed")
	}

	fmt.Fprintf(s.output, "[StatsD] Listening on port %d\n", s.port)

	buffer := make([]byte, 10000)

	for !s.isClosed() {
		numBytes, err := s.connection.Read(buffer)

		// blocking read returns an error when connection is closed
		if err != nil && s.isClosed() {
			break
		}

		if err != nil {
			fmt.Printf("[StatsD] Read error %s: \n", err)
			continue
		}

		rawMetrics := buffer[0:numBytes]
		splitMetrics := bytes.Split(rawMetrics, []byte("\n"))

		for _, rawMetric := range splitMetrics {
			rawMetric = bytes.TrimSpace(rawMetric)
			if len(rawMetric) == 0 {
				continue
			}
			metric := Parse(string(rawMetric))
			s.logMetric(metric)
		}
	}

	return nil
}

func (s *server) Close() error {
	atomic.StoreInt32(&s.closed, 1)

	s.outputLock.Lock()
	defer s.outputLock.Unlock()

	fmt.Fprintf(s.output, "[StatsD] Shutting down\n")

	return s.connection.Close()
}

func (s *server) logMetric(metric Metric) {
	s.outputLock.Lock()
	defer s.outputLock.Unlock()

	fmt.Fprintf(s.output, s.formatter.Format(metric))
}

func (s *server) isClosed() bool {
	return atomic.LoadInt32(&s.closed) != 0
}
