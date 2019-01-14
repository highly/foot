package log

import (
	"strings"
)

type Level int

const (
	_ Level = iota
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
)

func ToLevel(level string) Level {
	switch strings.ToLower(level) {
	case "error":
		return ErrorLevel
	case "warn":
		return WarnLevel
	case "debug":
		return DebugLevel
	default:
		return InfoLevel
	}
}
