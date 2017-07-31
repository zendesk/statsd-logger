package statsdLogger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseMetric(t *testing.T) {
	cases := []struct {
		input    string
		expected metric
	}{
		{
			input:    "catkins.cool.metric:1|c#city:melbourne,country:au",
			expected: metric{name: "catkins.cool.metric", value: "1|c", tags: "city:melbourne country:au"},
		},
		{
			input:    "another_rad.metric:142|ms|@0.5",
			expected: metric{name: "another_rad.metric", value: "142|ms|@0.5", tags: ""},
		},
	}

	for _, tc := range cases {
		assert.Equal(t, tc.expected, parseMetric(tc.input))
	}
}
