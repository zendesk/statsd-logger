package statsdLogger

import "strings"

type metric struct {
	name  string
	value string
	tags  string
}

func parseMetric(rawMetric string) metric {
	metricNameAndRest := strings.SplitN(rawMetric, ":", 2)
	name := metricNameAndRest[0]
	valueAndTags := strings.SplitN(metricNameAndRest[1], "#", 2)
	value := valueAndTags[0]
	tags := ""

	if len(valueAndTags) > 1 {
		tags = valueAndTags[1]
		tags = strings.Replace(tags, ",", " ", -1)
	}

	return metric{name: name, value: value, tags: tags}
}
