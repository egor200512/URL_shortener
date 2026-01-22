package nats

import "strings"

func join(prefix, event string) string {
	p := strings.TrimSuffix(prefix, ".")
	if p == "" {
		return event
	}
	return p + "." + event
}

func Created(prefix string) string { return join(prefix, "created") }
func Fetched(prefix string) string { return join(prefix, "fetched") }
func Deleted(prefix string) string { return join(prefix, "deleted") }
