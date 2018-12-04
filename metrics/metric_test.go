package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	cases := []struct {
		input    string
		expected Metric
	}{
		{
			input:    "catkins.cool.metric:1|c#city:melbourne,country:au",
			expected: Metric{Name: "catkins.cool.metric", Value: "1|c", Tags: "city:melbourne country:au"},
		},
		{
			input:    "another_rad.metric:142|ms|@0.5",
			expected: Metric{Name: "another_rad.metric", Value: "142|ms|@0.5", Tags: ""},
		},
	}

	for _, tc := range cases {
		assert.Equal(t, tc.expected, Parse(tc.input))
	}
}
