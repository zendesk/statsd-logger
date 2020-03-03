package metrics

import (
	"fmt"
	"strings"

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
		c.formatTags(metric.Tags))
}

func (c colorFormatter) formatTags(rawTags string) string {
	tags := strings.Split(rawTags, " ")

	formattedTags := []string{}
	for _, tag := range tags {
		tagParts := strings.SplitN(tag, ":", 2)
		formattedTag := fmt.Sprintf("%s%s", color.CyanString(tagParts[0]+":"), color.WhiteString(tagParts[1]))
		formattedTags = append(formattedTags, formattedTag)
	}

	return strings.Join(formattedTags, " ")
}
