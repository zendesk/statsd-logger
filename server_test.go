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
)

func TestListen(t *testing.T) {
	color.NoColor = true

	output := bytes.NewBuffer(make([]byte, 1000))
	server, err := NewWithWriter(":5678", output)

	go func() { server.Listen() }()

	assert.Nil(t, err)
	assert.NotNil(t, server)

	sendMetric(Metric{Name: "hello", Value: "1|c", Tags: "country:australia"}, 5678)

	<-time.NewTimer(5 * time.Millisecond).C

	assert.Contains(t, output.String(), "[StatsD] Listening on port 5678")
	assert.Contains(t, output.String(), "[StatsD] hello 1|c country:australia")
	assert.NotContains(t, output.String(), "[StatsD] Shutting down")

	server.Close()

	assert.Contains(t, output.String(), "[StatsD] Shutting down")
}

func TestClose(t *testing.T) {
	server, err := NewWithWriter(":5679", ioutil.Discard)
	assert.Nil(t, err)
	assert.NotNil(t, server)

	go func() { server.Listen() }()

	<-time.NewTimer(5 * time.Millisecond).C

	err = server.Close()
	assert.Nil(t, err)
}

func sendMetric(metric Metric, port int) {
	address, _ := net.ResolveUDPAddr("udp", ":5678")
	conn, _ := net.DialUDP("udp", nil, address)
	defer conn.Close()

	rawMetric := fmt.Sprintf("%s:%s#%s", metric.Name, metric.Value, metric.Tags)
	conn.Write([]byte(rawMetric))
}
