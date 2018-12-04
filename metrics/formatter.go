package metrics

import (
	"fmt"

	"github.com/fatih/color"
)

// DefaultFormatter provides output format for metrics
var DefaultFormatter MetricFormatter = colorFormatter{}

// MetricFormatter formats metrics as a string for logging
type MetricFormatter interface {
	// Format a metric
	Format(metric Metric) string
}

type colorFormatter struct{}

var _ MetricFormatter = colorFormatter{}

func (c colorFormatter) Format(metric Metric) string {
	return fmt.Sprintf(
		"[StatsD] %s %s %s\n",
		color.BlueString(metric.Name),
		color.YellowString(metric.Value),
		color.CyanString(metric.Tags))
}
