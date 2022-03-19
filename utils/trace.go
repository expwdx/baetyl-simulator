package utils

import (
	"time"

	"baetyl-simulator/middleware/log"
)

// Trace print elapsed time
func Trace(f func(string, ...log.Field), msg string, fields ...log.Field) func() {
	start := time.Now()
	return func() {
		fields := append(fields, log.Any("cost", time.Since(start)))
		f(msg, fields...)
	}
}
