package statsdLogger

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"testing"
	"time"

	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	testServerAddress = ":5678"
	testSErverPort    = 5678
)

func TestListen(t *testing.T) {
	color.NoColor = true

	output := bytes.NewBuffer(make([]byte, 1000))
	server, err := New(testServerAddress, WithWriter(output))

	go func() { server.Listen() }()

	assert.Nil(t, err)
	assert.NotNil(t, server)

	sendMetric(Metric{Name: "hello", Value: "1|c", Tags: "country:australia"})

	<-time.NewTimer(5 * time.Millisecond).C

	assert.Contains(t, output.String(), "[StatsD] Listening on port 5678")
	assert.Contains(t, output.String(), "[StatsD] hello 1|c country:australia")
	assert.NotContains(t, output.String(), "[StatsD] Shutting down")

	server.Close()

	assert.Contains(t, output.String(), "[StatsD] Shutting down")
}

func TestInvalidAddress(t *testing.T) {
	server, err := New("abcd", WithWriter(ioutil.Discard))

	assert.Error(t, err, "[StatsD] Invalid address")

	if err == nil {
		server.Close()
	}
}

func TestClose(t *testing.T) {
	server, err := New(testServerAddress, WithWriter(ioutil.Discard))
	assert.Nil(t, err)
	assert.NotNil(t, server)

	go func() { server.Listen() }()
	sendMetric(Metric{Name: "hello", Value: "1|c", Tags: "country:australia"})
	<-time.NewTimer(5 * time.Millisecond).C

	err = server.Close()
	assert.Nil(t, err)
}

func TestWithFormatter(t *testing.T) {
	formatter := &mockFormatter{}
	metric := Metric{Name: "cool_metric"}
	formatter.On("Format", metric).Return("ate 2 burgers")

	server, err := New(testServerAddress, WithFormatter(formatter), WithWriter(ioutil.Discard))
	defer server.Close()
	assert.Nil(t, err)

	go func() { server.Listen() }()

	sendMetric(metric)
	<-time.NewTimer(5 * time.Millisecond).C

	formatter.AssertCalled(t, "Format", metric)
}

func sendMetric(metric Metric) {
	address, _ := net.ResolveUDPAddr("udp", testServerAddress)
	conn, _ := net.DialUDP("udp", nil, address)
	defer conn.Close()

	rawMetric := fmt.Sprintf("%s:%s#%s", metric.Name, metric.Value, metric.Tags)
	conn.Write([]byte(rawMetric))
}

type mockFormatter struct {
	mock.Mock
}

var _ MetricFormatter = &mockFormatter{}

func (m *mockFormatter) Format(metric Metric) string {
	args := m.MethodCalled("Format", metric)
	return args.String(0)
}
