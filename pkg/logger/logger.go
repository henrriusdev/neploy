package logger

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/fatih/color"
)

func getCaller() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return "???"
	}
	parts := strings.Split(file, "/")
	file = parts[len(parts)-1]
	return fmt.Sprintf("%s:%d", file, line)
}

// Error logs an error message (🔴 Rojo claro)
func Error(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	caller := getCaller()
	color.New(color.FgHiRed).Printf("[ERROR] %s | %s\n", caller, msg)
}

// Info logs an info message (🔵 Azul cielo)
func Info(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	caller := getCaller()
	color.New(color.FgHiCyan).Printf("[INFO] %s | %s\n", caller, msg)
}

// Debug logs a debug message with detailed variable output (🟣 Púrpura suave)
func Debug(format string, args ...interface{}) {
	msg := spew.Sprintf(format, args...) // Usa spew para imprimir structs, maps, slices, etc.
	caller := getCaller()
	color.New(color.FgHiMagenta).Printf("[DEBUG] %s | %s\n", caller, msg)
}

// Warn logs a warning message (🟡 Amarillo dorado)
func Warn(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	caller := getCaller()
	color.New(color.FgHiYellow).Printf("[WARN] %s | %s\n", caller, msg)
}
