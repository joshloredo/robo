package logger

import (
	"fmt"
	"time"
)

// Color codes for terminal output
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
)

// Log levels
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

// Logger structure
type Logger struct {
	MinLevel LogLevel
}

// New creates a new Logger instance
func New(minLevel LogLevel) *Logger {
	return &Logger{
		MinLevel: minLevel,
	}
}

// formatMessage creates a formatted log message with timestamp and level
func (l *Logger) formatMessage(level LogLevel, message string, args ...interface{}) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	levelStr := ""
	color := ""

	switch level {
	case DEBUG:
		levelStr = "DEBUG"
		color = Purple
	case INFO:
		levelStr = "INFO"
		color = Blue
	case WARN:
		levelStr = "WARN"
		color = Yellow
	case ERROR:
		levelStr = "ERROR"
		color = Red
	case FATAL:
		levelStr = "FATAL"
		color = Red
	}

	// Format the message with any provided arguments
	formattedMsg := fmt.Sprintf(message, args...)

	return fmt.Sprintf("%s[%s] %s%-5s%s %s",
		color, timestamp, color, levelStr, Reset, formattedMsg)
}

// Debug logs a debug message
func (l *Logger) Debug(message string, args ...interface{}) {
	if l.MinLevel <= DEBUG {
		fmt.Println(l.formatMessage(DEBUG, message, args...))
	}
}

// Info logs an info message
func (l *Logger) Info(message string, args ...interface{}) {
	if l.MinLevel <= INFO {
		fmt.Println(l.formatMessage(INFO, message, args...))
	}
}

// Warn logs a warning message
func (l *Logger) Warn(message string, args ...interface{}) {
	if l.MinLevel <= WARN {
		fmt.Println(l.formatMessage(WARN, message, args...))
	}
}

// Error logs an error message
func (l *Logger) Error(message string, args ...interface{}) {
	if l.MinLevel <= ERROR {
		fmt.Println(l.formatMessage(ERROR, message, args...))
	}
}

// Fatal logs a fatal message and exits the program
func (l *Logger) Fatal(message string, args ...interface{}) {
	if l.MinLevel <= FATAL {
		fmt.Println(l.formatMessage(FATAL, message, args...))
		panic(message)
	}
}
