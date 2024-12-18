package logger

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
)

var defaultLogger = log.New(os.Stdout, "", log.LstdFlags)

// SetLogger configures the global logger
func SetLogger() {
	defaultLogger.SetFlags(log.LstdFlags)
}

func getCaller() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return "???"
	}
	// Get only the file name without the full path
	parts := strings.Split(file, "/")
	file = parts[len(parts)-1]
	return fmt.Sprintf("%s:%d", file, line)
}

// Error logs an error message
func Error(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	caller := getCaller()
	defaultLogger.Printf("[ERROR] %s | %s", caller, msg)
}

// Info logs an info message
func Info(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	caller := getCaller()
	defaultLogger.Printf("[INFO] %s | %s", caller, msg)
}

// Debug logs a debug message
func Debug(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	caller := getCaller()
	defaultLogger.Printf("[DEBUG] %s | %s", caller, msg)
}

// Warn logs a warning message
func Warn(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	caller := getCaller()
	defaultLogger.Printf("[WARN] %s | %s", caller, msg)
}
