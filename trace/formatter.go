package trace

import (
	"fmt"
	"strconv"
	"time"

	"github.com/fatih/color"
)

// Formatter takes a span and formats the output
type Formatter interface {
	Format(span Span) string
}

// ColourFormatter formats spans in ANSI colour
type ColourFormatter struct{}

// Format taks a span and formats it in ANSI colour
func (ColourFormatter) Format(span Span) string {
	duration, _ := time.ParseDuration(fmt.Sprintf("%dns", span.Duration))

	tags := ""

	for k, v := range span.Meta {
		tags = fmt.Sprintf("%s %s:%s", tags, color.CyanString(k), strconv.Quote(v))
	}

	return fmt.Sprintf(
		"[Trace] %s %s %s %s %s %s%s\n",
		color.HiCyanString(span.Service),
		color.GreenString(span.Operation),
		color.MagentaString(span.Resource),
		color.WhiteString(span.Type),
		color.YellowString("%s", duration),
		color.WhiteString("%d / %d", span.ParentID, span.SpanID),
		tags,
	)
}
