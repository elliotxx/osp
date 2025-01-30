// Package log provides a simple logging package with support for hierarchical logging.
//
// The package supports four log levels: Debug, Info, Success, and Error.
// Each log level has its own prefix symbol:
//   - Debug:   » (only shown when verbose mode is enabled)
//   - Info:    +
//   - Success: ✓
//   - Error:   ×
//
// The package also supports hierarchical logging with indentation levels.
// You can use L(level) to specify the indentation level (each level adds 2 spaces).
//
// Basic usage:
//
//	// Simple logging
//	log.Info("Processing item %d", 1)
//	// Output: + Processing item 1
//
//	// Hierarchical logging
//	log.Info("Found 2 items").
//	    L(1).Info("Processing item 1").
//	    L(1).Success("Item 1 processed").
//	    L(1).Info("Processing item 2").
//	    L(1).Error("Failed to process item 2")
//	// Output:
//	// + Found 2 items
//	//   + Processing item 1
//	//   ✓ Item 1 processed
//	//   + Processing item 2
//	//   × Failed to process item 2
//
//	// Debug logging (only shown when verbose mode is enabled)
//	log.SetVerbose(true)
//	log.Debug("Debug message")
//	// Output: » Debug message
//
// All logging functions return a Logger pointer, allowing for method chaining:
//
//	log.Info("Starting process").
//	    L(1).Info("Step 1").
//	    L(2).Debug("Detail 1").
//	    L(2).Debug("Detail 2").
//	    L(1).Success("Step 1 completed")
//
// This is particularly useful for maintaining consistent indentation levels
// throughout a process while keeping the code readable.
package log

import (
	"fmt"
)

var (
	verbose bool
)

// Logger represents a logger with a specific indentation level
type Logger struct {
	level int
}

// getIndent returns the current indentation string
func (l *Logger) getIndent() string {
	indentStr := ""
	for i := 0; i < l.level; i++ {
		indentStr += "  " // Two spaces per level
	}
	return indentStr
}

// L creates a new Logger with the specified indentation level.
// Each level adds 2 spaces of indentation.
//
// Example:
//
//	log.Info("Parent").L(1).Info("Child")
//	// Output:
//	// + Parent
//	//   + Child
func L(level int) *Logger {
	return &Logger{level: level}
}

// Debug prints debug message if verbose is true and returns the logger.
// Debug messages are prefixed with "»" and are only shown when verbose mode is enabled.
//
// Example:
//
//	log.SetVerbose(true)
//	log.Debug("Processing data")
//	// Output: » Processing data
func (l *Logger) Debug(format string, args ...interface{}) *Logger {
	if verbose {
		fmt.Printf(l.getIndent()+"» "+format+"\n", args...)
	}
	return l
}

// Info prints info message and returns the logger.
// Info messages are prefixed with "+".
//
// Example:
//
//	log.Info("Processing %d items", 5)
//	// Output: + Processing 5 items
func (l *Logger) Info(format string, args ...interface{}) *Logger {
	fmt.Printf(l.getIndent()+"+ "+format+"\n", args...)
	return l
}

// Success prints success message and returns the logger.
// Success messages are prefixed with "✓".
//
// Example:
//
//	log.Success("All items processed")
//	// Output: ✓ All items processed
func (l *Logger) Success(format string, args ...interface{}) *Logger {
	fmt.Printf(l.getIndent()+"✓ "+format+"\n", args...)
	return l
}

// Error prints error message and returns the logger.
// Error messages are prefixed with "×".
//
// Example:
//
//	log.Error("Failed to process item: %v", err)
//	// Output: × Failed to process item: connection refused
func (l *Logger) Error(format string, args ...interface{}) *Logger {
	fmt.Printf(l.getIndent()+"× "+format+"\n", args...)
	return l
}

// SetVerbose sets the verbose flag.
// When verbose is true, Debug messages will be printed.
// When verbose is false, Debug messages will be suppressed.
func SetVerbose(v bool) {
	verbose = v
}

// Global functions that return a logger with level 0

// Debug is a convenience function that creates a level 0 logger and calls Debug.
func Debug(format string, args ...interface{}) *Logger {
	return L(0).Debug(format, args...)
}

// Info is a convenience function that creates a level 0 logger and calls Info.
func Info(format string, args ...interface{}) *Logger {
	return L(0).Info(format, args...)
}

// Success is a convenience function that creates a level 0 logger and calls Success.
func Success(format string, args ...interface{}) *Logger {
	return L(0).Success(format, args...)
}

// Error is a convenience function that creates a level 0 logger and calls Error.
func Error(format string, args ...interface{}) *Logger {
	return L(0).Error(format, args...)
}
