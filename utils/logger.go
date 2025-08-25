package utils

import (
	"fmt"
	"sync"
	"time"
)

var (
	logMutex  sync.Mutex
	logPrefix string
)

// LogInfo prints an info level message
func LogInfo(format string, v ...interface{}) {
	logWithLevel("INFO", format, v...)
}

// LogWarning prints a warning level message
func LogWarning(format string, v ...interface{}) {
	logWithLevel("WARN", format, v...)
}

// LogError prints an error level message
func LogError(format string, v ...interface{}) {
	logWithLevel("ERROR", format, v...)
}

// logWithLevel prints a message with specified log level
func logWithLevel(level, format string, v ...interface{}) {
	logMutex.Lock()
	defer logMutex.Unlock()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, v...)

	if logPrefix != "" {
		fmt.Printf("[%s] %-5s %s%s\n", timestamp, level, logPrefix, message)
	} else {
		fmt.Printf("[%s] %-5s %s\n", timestamp, level, message)
	}
}

// WithLogPrefix sets a temporary prefix for log messages
func WithLogPrefix(prefix string, fn func()) {
	logMutex.Lock()
	originalPrefix := logPrefix
	logPrefix = prefix
	logMutex.Unlock()

	defer func() {
		logMutex.Lock()
		logPrefix = originalPrefix
		logMutex.Unlock()
	}()

	fn()
}