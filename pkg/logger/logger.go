package logger

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/fatih/color"
)

var defaultLogger = color.New(color.FgWhite)

// SetLogger configures the global logger
func SetLogger() {
	defaultLogger = color.New(color.FgWhite)
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
	defaultLogger.Add(color.FgWhite, color.BgRed).Printf("[ERROR] %s | %s\n", caller, msg)
}

// Info logs an info message
func Info(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	caller := getCaller()
	defaultLogger.Add(color.FgWhite, color.BgGreen).Printf("[INFO] %s | %s\n", caller, msg)
}

// Debug logs a debug message
func Debug(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	caller := getCaller()
	defaultLogger.Add(color.FgWhite, color.BgBlue).Printf("[DEBUG] %s | %s\n", caller, msg)
}

// Warn logs a warning message
func Warn(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	caller := getCaller()
	defaultLogger.Add(color.FgWhite, color.BgYellow).Printf("[WARN] %s | %s\n", caller, msg)
}
