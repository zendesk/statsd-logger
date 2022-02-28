package trace

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vmihailenco/msgpack"
)

const testServerAddr = ":56789"

func TestServer_Listen(t *testing.T) {
	output := bytes.Buffer{}
	server := NewServer(testServerAddr, WithOutput(&output))

	go func() {
		server.Listen()
	}()
	defer server.Close()

	time.Sleep(50 * time.Millisecond)

	traces := Traces{
		Trace{
			&Span{
				Service:   "a_service",
				Operation: "an_operation",
				Resource:  "a_resource",
				Type:      "web",
				Duration:  1000000, // 1 second,
				Meta: map[string]string{
					"a_key": "a_value",
				},
				SpanID:   987,
				TraceID:  654,
				ParentID: 321,
			},
			&Span{
				Service:   "another_service",
				Operation: "another_operation",
				Resource:  "another_resource",
				Type:      "db",
				Duration:  1000000, // 1 second,
				Meta: map[string]string{
					"another_key": "another_value",
				},
				SpanID:   789,
				TraceID:  456,
				ParentID: 124,
			},
		},
	}

	encodedTrace, err := msgpack.Marshal(traces)
	requestBody := bytes.NewBuffer(encodedTrace)
	assert.NoError(t, err)

	response, err := http.Post(fmt.Sprintf("http://%s/v0.3/traces", testServerAddr), "application/msgpack", requestBody)

	assert.NotNil(t, response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotEmpty(t, output)
}
