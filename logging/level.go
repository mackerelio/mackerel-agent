package logging

//go:generate stringer -type=level level.go
type level uint8

// loglevels
const (
	UNKNOWN level = iota
	TRACE
	DEBUG
	INFO
	WARNING
	ERROR
	CRITICAL
)
