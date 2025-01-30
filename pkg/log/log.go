package log

import (
	"fmt"
)

var verbose bool

// SetVerbose sets the verbose flag
func SetVerbose(v bool) {
	verbose = v
}

// Debug prints debug message if verbose is true
func Debug(format string, args ...interface{}) {
	if verbose {
		fmt.Printf("» "+format+"\n", args...)
	}
}

// Info prints info message
func Info(format string, args ...interface{}) {
	fmt.Printf("+ "+format+"\n", args...)
}

// Success prints success message
func Success(format string, args ...interface{}) {
	fmt.Printf("✓ "+format+"\n", args...)
}

// Error prints error message
func Error(format string, args ...interface{}) {
	fmt.Printf("× "+format+"\n", args...)
}
