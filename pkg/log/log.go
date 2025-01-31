// Package log provides a simple logging package with support for hierarchical logging and colors.
//
// The package supports five built-in log levels with their own prefix symbols and colors:
//   - Debug:   + (light gray, only shown when verbose mode is enabled)
//   - Info:    » (blue)
//   - Warn:    ! (yellow)
//   - Success: ✓ (green)
//   - Error:   × (red)
//
// The package also supports hierarchical logging with indentation levels and custom prefixes.
// You can use:
//   - L(level) to specify the indentation level (each level adds 2 spaces)
//   - P(prefix) to specify a custom prefix
//   - C(color) to specify a custom color
//   - N() to disable newline at the end of the message
//   - B() to set the text to bold
//
// Basic usage:
//
//	// Simple logging with built-in levels (with default colors)
//	log.Info("Processing item %d", 1)
//	// Output: » Processing item 1 (in blue)
//
//	// Hierarchical logging with custom prefix and colors
//	log.Info("Found 2 items").
//	    L(1).P("→").C(log.ColorCyan).Log("Processing item 1").
//	    L(1).Success("Item 1 processed").
//	    L(1).P("→").C(log.ColorCyan).Log("Processing item 2").
//	    L(1).Error("Failed to process item 2")
//	// Output:
//	// » Found 2 items (in blue)
//	//   → Processing item 1 (in cyan)
//	//   ✓ Item 1 processed (in green)
//	//   → Processing item 2 (in cyan)
//	//   × Failed to process item 2 (in red)
//
//	// Debug logging (only shown when verbose mode is enabled)
//	log.SetVerbose(true)
//	log.Debug("Debug message")
//	// Output: + Debug message (in light gray)
//
//	// Disable newline at the end of the message
//	log.N().Info("Enter your name: ")
//	// Output: » Enter your name: (without newline)
//
//	// Using bold text
//	// Method 1: Make entire message bold
//	log.B().Info("This entire message is bold")
//	// Output: » This entire message is bold (in bold)
//
//	// Method 2: Make part of the message bold
//	log.Info("Current user: %s", log.Bold("elliotxx"))
//	// Output: » Current user: elliotxx (with "elliotxx" in bold)
//
//	// Method 3: Multiple bold parts
//	log.Info("Found %s issues, %s are critical", log.Bold("10"), log.Bold("5"))
//	// Output: » Found 10 issues, 5 are critical (with "10" and "5" in bold)
//
//	// Method 4: Combine bold with colors
//	log.B().C(log.ColorRed).Error("Critical error: %s", log.Bold("permission denied"))
//	// Output: × Critical error: permission denied (in red, with entire message and "permission denied" in bold)
//
// All logging functions return a new Logger pointer, allowing for method chaining:
//
//	// L(level) sets the indentation level
//	// P(prefix) sets a custom prefix
//	// C(color) sets a custom color
//	// B() sets the text to bold
//	// Log() outputs message with current level, prefix and color
//	log.L(1).P("→").C(log.ColorYellow).B().Log("Message 1").Log("Message 2")
//	// Output:
//	//   → Message 1 (in yellow and bold)
//	//   → Message 2 (in yellow and bold)
//
// Each method (L, P, C, B, Log, etc.) returns a new Logger instance with the updated settings,
// making it safe for concurrent use and allowing for flexible logging patterns.
//
// Available colors for use with C():
//   - log.ColorReset  (reset to default color)
//   - log.ColorRed    (red)
//   - log.ColorGreen  (green)
//   - log.ColorYellow (yellow)
//   - log.ColorBlue   (blue)
//   - log.ColorPurple (purple)
//   - log.ColorCyan   (cyan)
//   - log.ColorGray   (light gray)
//   - log.StyleBold   (bold style)
package log

import (
	"fmt"
	"strings"
)

var (
	verbose bool
	noColor bool // If true, disable color output
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

// getColor returns the color code if color output is enabled, otherwise returns empty string
func getColor(color string) string {
	if noColor {
		return ""
	}
	return color
}

// getColorReset returns the color reset code if color output is enabled, otherwise returns empty string
func getColorReset() string {
	if noColor {
		return ""
	}
	return colorReset
}

// SetNoColor sets the global color output setting
func SetNoColor(disable bool) {
	noColor = disable
}

// Logger represents a logger with a specific indentation level and prefix
type Logger struct {
	level     int    // indentation level
	prefix    string // prefix symbol
	color     string // ANSI color code
	noNewline bool   // control whether to output newline at the end
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
	newLogger.prefix = prefix
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

// B sets the text to bold and returns a new logger.
//
// Example:
//
//	log.B().Info("This is bold")
//	// Output: » This is bold (in bold)
func (l *Logger) B() *Logger {
	newLogger := *l
	if newLogger.color == "" {
		newLogger.color = styleBold
	} else {
		newLogger.color += styleBold
	}
	return &newLogger
}

// N disables the newline at the end of the log message
func (l *Logger) N() *Logger {
	newLogger := *l
	newLogger.noNewline = true
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
	indent := l.getIndent()
	msg := fmt.Sprintf(format, args...)
	color := getColor(l.color)
	reset := getColorReset()

	// Handle newlines in the message
	lines := strings.Split(msg, "\n")
	for i, line := range lines {
		if line == "" {
			continue
		}
		fmt.Print(indent)
		fmt.Print(color)
		if l.prefix != "" {
			fmt.Printf("%s ", l.prefix)
		}
		fmt.Print(line)
		fmt.Print(reset)
		if !l.noNewline || i < len(lines)-1 {
			fmt.Println()
		}
	}

	return New().L(l.level).P(l.prefix).C(l.color)
}

// Debug prints debug message if verbose is true and returns a new logger.
// Debug messages are prefixed with "+" and are only shown when verbose mode is enabled.
//
// Example:
//
//	log.SetVerbose(true)
//	log.Debug("Debug message")
//	// Output: + Debug message
func (l *Logger) Debug(format string, args ...interface{}) *Logger {
	if !verbose {
		return l
	}
	return l.P("+").C(colorGray).Log(format, args...)
}

// Info prints info message and returns a new logger.
// Info messages are prefixed with "»".
//
// Example:
//
//	log.Info("Processing %d items", 5)
//	// Output: » Processing 5 items
func (l *Logger) Info(format string, args ...interface{}) *Logger {
	return l.P("»").C(colorBlue).Log(format, args...)
}

// Warn prints warning message and returns a new logger.
// Warning messages are prefixed with "!".
//
// Example:
//
//	log.Warn("Low disk space: %d%%", 10)
//	// Output: ! Low disk space: 10%
func (l *Logger) Warn(format string, args ...interface{}) *Logger {
	return l.P("!").C(colorYellow).Log(format, args...)
}

// Success prints success message and returns a new logger.
// Success messages are prefixed with "✓".
//
// Example:
//
//	log.Success("All items processed")
//	// Output: ✓ All items processed
func (l *Logger) Success(format string, args ...interface{}) *Logger {
	return l.P("✓").C(colorGreen).Log(format, args...)
}

// Error prints error message and returns a new logger.
// Error messages are prefixed with "×".
//
// Example:
//
//	log.Error("Failed to process item: %v", err)
//	// Output: × Failed to process item: connection refused
func (l *Logger) Error(format string, args ...interface{}) *Logger {
	return l.P("×").C(colorRed).Log(format, args...)
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

// B is a convenience function that creates a new logger with bold text.
func B() *Logger {
	return New().B()
}

// N is a convenience function that creates a new logger and disables the newline
func N() *Logger {
	return New().N()
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

// Warn is a convenience function that creates a new logger and calls Warn.
func Warn(format string, args ...interface{}) *Logger {
	return New().Warn(format, args...)
}

// Success is a convenience function that creates a new logger and calls Success.
func Success(format string, args ...interface{}) *Logger {
	return New().Success(format, args...)
}

// Error is a convenience function that creates a new logger and calls Error.
func Error(format string, args ...interface{}) *Logger {
	return New().Error(format, args...)
}

// Bold wraps text in bold style
//
// Example:
//
//	log.Info("Normal text %s more text", log.Bold("bold text"))
//	// Output: » Normal text bold text more text (with "bold text" in bold)
func Bold(text string) string {
	return styleBold + text + colorReset
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
