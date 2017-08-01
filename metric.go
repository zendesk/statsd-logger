package statsdLogger

import "strings"

// Metric is an intermediate representation of a raw statsd metric for easier presentation
type Metric struct {
	Name  string
	Value string
	Tags  string
}

// ParseMetric takes a raw statsd metric and returns a populated Metric
func ParseMetric(rawMetric string) Metric {
	metricNameAndRest := strings.SplitN(rawMetric, ":", 2)
	name := metricNameAndRest[0]
	valueAndTags := strings.SplitN(metricNameAndRest[1], "#", 2)
	value := valueAndTags[0]
	tags := ""

	if len(valueAndTags) > 1 {
		tags = valueAndTags[1]
		tags = strings.Replace(tags, ",", " ", -1)
	}

	return Metric{Name: name, Value: value, Tags: tags}
}
