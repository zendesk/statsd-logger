package metrics

import (
	"fmt"
	"strings"
)

// Metric is an intermediate representation of a raw statsd metric for easier presentation
type Metric struct {
	Name  string
	Value string
	Tags  string
}

// Parse takes a raw statsd metric and returns a populated Metric
func Parse(rawMetric string) (metric Metric) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("error parsing metric: %+v\n", r)
		}
	}()

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
