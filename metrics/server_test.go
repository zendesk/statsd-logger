package metrics

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	testServerAddress = ":5678"
	testServerPort    = 5678
)

func TestListen(t *testing.T) {
	color.NoColor = true

	output := &buffer{
		bytes.NewBuffer(make([]byte, 1000)),
		sync.Mutex{},
	}

	server, err := NewServer(testServerAddress, WithWriter(output))

	go func() { server.Listen() }()

	assert.Nil(t, err)
	assert.NotNil(t, server)

	sendMetric(Metric{Name: "hello", Value: "1|c", Tags: "country:australia"})
	sendMetrics(
		Metric{Name: "another_metric", Value: "2|c", Tags: "country:greece"},
		Metric{Name: "another_metric", Value: "3|c", Tags: "country:malta"},
	)

	<-time.NewTimer(5 * time.Millisecond).C

	assert.Contains(t, output.String(), "[StatsD] Listening on port 5678")
	assert.Contains(t, output.String(), "[StatsD] hello 1|c country:australia")
	assert.Contains(t, output.String(), "[StatsD] another_metric 2|c country:greece")
	assert.Contains(t, output.String(), "[StatsD] another_metric 3|c country:malta")
	assert.NotContains(t, output.String(), "[StatsD] Shutting down")

	server.Close()

	assert.Contains(t, output.String(), "[StatsD] Shutting down")
}

func TestInvalidAddress(t *testing.T) {
	server, err := NewServer("abcd", WithWriter(ioutil.Discard))

	assert.Error(t, err, "[StatsD] Invalid address")

	if err == nil {
		server.Close()
	}
}

func TestClose(t *testing.T) {
	// TODO: for some reason, having this test use the `testServerAddress` seems to cause panics
	// probably should figure out why
	server, err := NewServer(":0", WithWriter(ioutil.Discard))
	assert.Nil(t, err)
	assert.NotNil(t, server)

	go func() { server.Listen() }()

	<-time.NewTimer(5 * time.Millisecond).C

	err = server.Close()
	assert.Nil(t, err)
}

func TestWithFormatter(t *testing.T) {
	formatter := &mockFormatter{}
	metric := Metric{Name: "cool_metric"}
	formatter.On("Format", metric).Return("ate 2 burgers")

	server, err := NewServer(testServerAddress, WithFormatter(formatter), WithWriter(ioutil.Discard))
	defer server.Close()
	assert.Nil(t, err)

	go func() { server.Listen() }()

	sendMetric(metric)
	<-time.NewTimer(5 * time.Millisecond).C

	formatter.AssertCalled(t, "Format", metric)
	server.Close()
}

func sendMetric(metric Metric) {
	address, _ := net.ResolveUDPAddr("udp", testServerAddress)
	conn, _ := net.DialUDP("udp", nil, address)
	defer conn.Close()

	rawMetric := fmt.Sprintf("%s:%s#%s", metric.Name, metric.Value, metric.Tags)
	conn.Write([]byte(rawMetric))
}

func sendMetrics(metrics ...Metric) {
	address, _ := net.ResolveUDPAddr("udp", testServerAddress)
	conn, _ := net.DialUDP("udp", nil, address)
	defer conn.Close()

	message := ""
	for _, metric := range metrics {
		message = message + "\n" + fmt.Sprintf("%s:%s#%s", metric.Name, metric.Value, metric.Tags)
	}
	conn.Write([]byte(message))
}

type mockFormatter struct {
	mock.Mock
}

// type assertion
var _ MetricFormatter = &mockFormatter{}

func (m *mockFormatter) Format(metric Metric) string {
	args := m.MethodCalled("Format", metric)
	return args.String(0)
}

// goroutine safe buffer to keep the race detector happy
type buffer struct {
	*bytes.Buffer
	mutex sync.Mutex
}

func (s *buffer) Write(p []byte) (n int, err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.Buffer.Write(p)
}

func (s *buffer) String() string {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.Buffer.String()
}
