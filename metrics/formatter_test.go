package metrics

import (
	"testing"

	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
)

func TestDefaultFormatterWithColorDisabled(t *testing.T) {
	formatter := DefaultFormatter
	color.NoColor = true

	output := formatter.Format(Metric{
		Name:  "exciting_metric",
		Value: "1|c",
		Tags:  "a_tag:123",
	})

	assert.Equal(t, output, "[StatsD] exciting_metric 1|c a_tag:123\n")
}

func TestDefaultFormatterWithColorEnabled(t *testing.T) {
	formatter := DefaultFormatter
	color.NoColor = false

	output := formatter.Format(Metric{
		Name:  "exciting_metric",
		Value: "1|c",
		Tags:  "a_tag:123",
	})

	assert.Equal(t, "[StatsD] \x1b[34mexciting_metric\x1b[0m \x1b[33m1|c\x1b[0m \x1b[36ma_tag:\x1b[0m\x1b[37m123\x1b[0m\n", output)
}
