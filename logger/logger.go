package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorGray   = "\033[37m"
	ColorBold   = "\033[1m"
)

// Logger levels
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

var levelColors = map[LogLevel]string{
	DEBUG: ColorGray,
	INFO:  ColorGreen,
	WARN:  ColorYellow,
	ERROR: ColorRed,
}

var levelNames = map[LogLevel]string{
	DEBUG: "DEBUG",
	INFO:  "INFO",
	WARN:  "WARN",
	ERROR: "ERROR",
}

// Custom logger instance
type PrettyLogger struct {
	showColors bool
	showCaller bool
	minLevel   LogLevel
}

var DefaultLogger = &PrettyLogger{
	showColors: true,
	showCaller: true,
	minLevel:   INFO,
}

// getCaller returns the file and line number of the caller
func (l *PrettyLogger) getCaller(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown"
	}

	// Get just the filename, not the full path
	filename := filepath.Base(file)
	return fmt.Sprintf("%s:%d", filename, line)
}

// formatLog formats the log message with colors and caller info
func (l *PrettyLogger) formatLog(level LogLevel, msg string) string {
	timestamp := time.Now().Format("2006/01/02 15:04:05")

	var levelStr string
	if l.showColors {
		color := levelColors[level]
		levelStr = fmt.Sprintf("%s%s%s%s", ColorBold, color, levelNames[level], ColorReset)
	} else {
		levelStr = levelNames[level]
	}

	var callerStr string
	if l.showCaller {
		caller := l.getCaller(4) // Skip formatLog -> logAtLevel -> Info/Warn/Error -> actual caller
		if l.showColors {
			callerStr = fmt.Sprintf(" %s[%s]%s", ColorCyan, caller, ColorReset)
		} else {
			callerStr = fmt.Sprintf(" [%s]", caller)
		}
	}

	if l.showColors {
		return fmt.Sprintf("%s%s%s %s%s: %s", ColorGray, timestamp, ColorReset, levelStr, callerStr, msg)
	} else {
		return fmt.Sprintf("%s %s%s: %s", timestamp, levelStr, callerStr, msg)
	}
}

// logAtLevel logs a message at the specified level
func (l *PrettyLogger) logAtLevel(level LogLevel, format string, args ...interface{}) {
	if level < l.minLevel {
		return
	}

	msg := fmt.Sprintf(format, args...)
	formatted := l.formatLog(level, msg)

	if level >= ERROR {
		fmt.Fprintln(os.Stderr, formatted)
	} else {
		fmt.Println(formatted)
	}
}

// Public logging functions
func Debug(format string, args ...interface{}) {
	DefaultLogger.logAtLevel(DEBUG, format, args...)
}

func Info(format string, args ...interface{}) {
	DefaultLogger.logAtLevel(INFO, format, args...)
}

func Warn(format string, args ...interface{}) {
	DefaultLogger.logAtLevel(WARN, format, args...)
}

func Error(format string, args ...interface{}) {
	DefaultLogger.logAtLevel(ERROR, format, args...)
}

// Convenience functions that match the standard log package
func Printf(format string, args ...interface{}) {
	Info(format, args...)
}

func Println(args ...interface{}) {
	msg := fmt.Sprintln(args...)
	// Remove trailing newline since our formatter adds it
	msg = strings.TrimSuffix(msg, "\n")
	Info("%s", msg)
}

// SetLogLevel sets the minimum log level
func SetLogLevel(level LogLevel) {
	DefaultLogger.minLevel = level
}

// DisableColors disables colored output
func DisableColors() {
	DefaultLogger.showColors = false
}

// DisableCaller disables showing caller information
func DisableCaller() {
	DefaultLogger.showCaller = false
}
