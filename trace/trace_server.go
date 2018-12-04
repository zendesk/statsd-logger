package trace

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// DefaultAddress is the default address to listen to traces
const DefaultAddress = ":8126"

// Server mimics the HTTP API exposed by the datadog APM agent
type Server struct {
	// tcp address to listen on
	Address   string
	server    *http.Server
	output    io.Writer
	formatter Formatter
}

// ServerOption provides optional configuration for initializing a new server
type ServerOption func(*serverConfig)

type serverConfig struct {
	output    io.Writer
	formatter Formatter
}

// WithOutput allows for overriding the output destination for traces
func WithOutput(output io.Writer) ServerOption {
	return func(config *serverConfig) {
		config.output = output
	}
}

// NewServer creates a new trace server
func NewServer(address string, options ...ServerOption) *Server {
	config := serverConfig{
		output:    os.Stdout,
		formatter: ColourFormatter{},
	}

	for _, opt := range options {
		opt(&config)
	}

	server := Server{
		Address:   address,
		output:    config.output,
		formatter: config.formatter,
	}

	mux := http.NewServeMux()

	// v0.3 & v0.4 formats are messagepack based
	mux.HandleFunc("/v0.3/traces", server.HandleTrace)
	mux.HandleFunc("/v0.4/traces", server.HandleTrace)

	httpServer := &http.Server{
		Addr:         address,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
	server.server = httpServer

	return &server
}

// Listen starts the trace server
func (server Server) Listen() error {
	fmt.Printf("[Trace] Listening at %s\n", server.Address)

	return server.server.ListenAndServe()
}

// HandleTrace receives tracing requests from datadog clients and prints them to the console
func (server Server) HandleTrace(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Printf("error reading body: %+v", err)
		w.WriteHeader(422)
		return
	}
	defer req.Body.Close()

	spans, err := DecodeSpans(body)
	if err != nil {
		log.Printf("error decoding spans: %+v", err)
		w.WriteHeader(422)
		return
	}

	for _, span := range spans {
		fmt.Fprintf(server.output, server.formatter.Format(span))
	}

	w.WriteHeader(200)
}

// Close shuts down the trace server
func (server Server) Close() error {
	fmt.Println("[Trace] Shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return server.server.Shutdown(ctx)
}
