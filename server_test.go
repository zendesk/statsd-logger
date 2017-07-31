package statsdLogger

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestListen(t *testing.T) {
	server, err := New(":0")

	assert.Nil(t, err)
	assert.NotNil(t, server)

	go func() { server.Listen() }()
	timer := time.NewTimer(time.Second)
	<-timer.C
	defer server.Close()
}

func TestShutdown(t *testing.T) {

}
