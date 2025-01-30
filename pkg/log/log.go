// Package log provides a simple logging package with support for hierarchical logging.
//
// The package supports four built-in log levels with their own prefix symbols:
//   - Debug:   » (only shown when verbose mode is enabled)
//   - Info:    +
//   - Success: ✓
//   - Error:   ×
//
// The package also supports hierarchical logging with indentation levels and custom prefixes.
// You can use L(level) to specify the indentation level (each level adds 2 spaces),
// and P(prefix) to specify a custom prefix.
//
// Basic usage:
//
//	// Simple logging with built-in levels
//	log.Info("Processing item %d", 1)
//	// Output: + Processing item 1
//
//	// Hierarchical logging with custom prefix
//	log.Info("Found 2 items").
//	    L(1).P("→").Log("Processing item 1").
//	    L(1).Success("Item 1 processed").
//	    L(1).P("→").Log("Processing item 2").
//	    L(1).Error("Failed to process item 2")
//	// Output:
//	// + Found 2 items
//	//   → Processing item 1
//	//   ✓ Item 1 processed
//	//   → Processing item 2
//	//   × Failed to process item 2
//
//	// Debug logging (only shown when verbose mode is enabled)
//	log.SetVerbose(true)
//	log.Debug("Debug message")
//	// Output: » Debug message
//
// All logging functions return a new Logger pointer, allowing for method chaining:
//
//	// L(level) sets the indentation level
//	// P(prefix) sets a custom prefix
//	// Log() outputs message with current level and prefix
//	log.L(1).P("→").Log("Message 1").Log("Message 2")
//	// Output:
//	//   → Message 1
//	//   → Message 2
//
// Each method (L, P, Log, etc.) returns a new Logger instance with the updated settings,
// making it safe for concurrent use and allowing for flexible logging patterns.
package log

import (
	"fmt"
)

var (
	verbose bool
)

// ANSI color codes
const (
	// Colors
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[90m"

	// Styles
	styleBold = "\033[1m"
)

// Logger represents a logger with a specific indentation level and prefix
type Logger struct {
	level  int    // indentation level
	prefix string // prefix symbol
	color  string // ANSI color code
}

// getIndent returns the current indentation string
func (l *Logger) getIndent() string {
	indentStr := ""
	for i := 0; i < l.level; i++ {
		indentStr += "  " // Two spaces per level
	}
	return indentStr
}

// P sets a custom prefix for the logger and returns the logger.
// The prefix will be used by subsequent Log calls.
//
// Example:
//
//	log.L(1).P("→").Log("Processing item")
//	// Output:
//	//   → Processing item
func (l *Logger) P(prefix string) *Logger {
	newLogger := *l
	newLogger.prefix = prefix + " "
	return &newLogger
}

// L sets the indentation level and returns a new logger.
// Each level adds 2 spaces of indentation.
//
// Example:
//
//	log.P("→").L(1).Log("Child message")
//	// Output:
//	//   → Child message
func (l *Logger) L(level int) *Logger {
	newLogger := *l
	newLogger.level = level
	return &newLogger
}

// C sets the color of the logger and returns a new logger.
// The color will be used by subsequent Log calls.
//
// Example:
//
//	log.L(1).C(colorRed).Log("Error message")
//	// Output:
//	//   Error message (in red)
func (l *Logger) C(color string) *Logger {
	newLogger := *l
	newLogger.color = color
	return &newLogger
}

// Log prints message with current level and prefix, then returns a new logger.
//
// Example:
//
//	log.L(1).P("→").Log("Message 1").Log("Message 2")
//	// Output:
//	//   → Message 1
//	//   → Message 2
func (l *Logger) Log(format string, args ...interface{}) *Logger {
	if l.color != "" {
		fmt.Printf(l.getIndent()+l.color+l.prefix+format+colorReset+"\n", args...)
	} else {
		fmt.Printf(l.getIndent()+l.prefix+format+"\n", args...)
	}
	newLogger := *l
	return &newLogger
}

// Debug prints debug message if verbose is true and returns a new logger.
// Debug messages are prefixed with "»" and are only shown when verbose mode is enabled.
//
// Example:
//
//	log.SetVerbose(true)
//	log.Debug("Processing data")
//	// Output: » Processing data
func (l *Logger) Debug(format string, args ...interface{}) *Logger {
	newLogger := *l
	newLogger.prefix = "» "
	newLogger.color = colorGray
	if verbose {
		fmt.Printf(newLogger.getIndent()+newLogger.color+newLogger.prefix+format+colorReset+"\n", args...)
	}
	return &newLogger
}

// Info prints info message and returns a new logger.
// Info messages are prefixed with "+".
//
// Example:
//
//	log.Info("Processing %d items", 5)
//	// Output: + Processing 5 items
func (l *Logger) Info(format string, args ...interface{}) *Logger {
	newLogger := *l
	newLogger.prefix = "+ "
	newLogger.color = colorBlue
	fmt.Printf(newLogger.getIndent()+newLogger.color+newLogger.prefix+format+colorReset+"\n", args...)
	return &newLogger
}

// Success prints success message and returns a new logger.
// Success messages are prefixed with "✓".
//
// Example:
//
//	log.Success("All items processed")
//	// Output: ✓ All items processed
func (l *Logger) Success(format string, args ...interface{}) *Logger {
	newLogger := *l
	newLogger.prefix = "✓ "
	newLogger.color = colorGreen
	fmt.Printf(newLogger.getIndent()+newLogger.color+newLogger.prefix+format+colorReset+"\n", args...)
	return &newLogger
}

// Error prints error message and returns a new logger.
// Error messages are prefixed with "×".
//
// Example:
//
//	log.Error("Failed to process item: %v", err)
//	// Output: × Failed to process item: connection refused
func (l *Logger) Error(format string, args ...interface{}) *Logger {
	newLogger := *l
	newLogger.prefix = "× "
	newLogger.color = colorRed
	fmt.Printf(newLogger.getIndent()+newLogger.color+newLogger.prefix+format+colorReset+"\n", args...)
	return &newLogger
}

// SetVerbose sets the verbose flag.
// When verbose is true, Debug messages will be printed.
// When verbose is false, Debug messages will be suppressed.
func SetVerbose(v bool) {
	verbose = v
}

// Global functions that return a new logger

// New creates a new logger with default settings (level 0, no prefix)
func New() *Logger {
	return &Logger{}
}

// L sets the indentation level and returns a new logger.
// Each level adds 2 spaces of indentation.
//
// Example:
//
//	log.L(1).P("→").Log("Child message")
//	// Output:
//	//   → Child message
func L(level int) *Logger {
	return &Logger{level: level}
}

// P sets a custom prefix for the logger and returns a new logger.
// The prefix will be used by subsequent Log calls.
//
// Example:
//
//	log.L(1).P("→").Log("Child message")
//	// Output:
//	//   → Child message
func P(prefix string) *Logger {
	return &Logger{prefix: prefix}
}

// C sets the color of the logger and returns a new logger.
// The color will be used by subsequent Log calls.
//
// Example:
//
//	log.L(1).C(colorRed).Log("Error message")
//	// Output:
//	//   Error message (in red)
func C(color string) *Logger {
	return &Logger{color: color}
}

// Log is a convenience function that creates a new logger and calls Log.
func Log(format string, args ...interface{}) *Logger {
	return New().Log(format, args...)
}

// Debug is a convenience function that creates a new logger and calls Debug.
func Debug(format string, args ...interface{}) *Logger {
	return New().Debug(format, args...)
}

// Info is a convenience function that creates a new logger and calls Info.
func Info(format string, args ...interface{}) *Logger {
	return New().Info(format, args...)
}

// Success is a convenience function that creates a new logger and calls Success.
func Success(format string, args ...interface{}) *Logger {
	return New().Success(format, args...)
}

// Error is a convenience function that creates a new logger and calls Error.
func Error(format string, args ...interface{}) *Logger {
	return New().Error(format, args...)
}

// Color constants for use with C() method
var (
	ColorReset  = colorReset
	ColorRed    = colorRed
	ColorGreen  = colorGreen
	ColorYellow = colorYellow
	ColorBlue   = colorBlue
	ColorPurple = colorPurple
	ColorCyan   = colorCyan
	ColorGray   = colorGray
	StyleBold   = styleBold
)
