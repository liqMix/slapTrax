package logger

import (
	"fmt"
	"image/color"
	"time"
)

var l = new()
var (
	Debug          = l.Debug
	Info           = l.Info
	Warn           = l.Warn
	Error          = l.Error
	Fatal          = l.Fatal
	GetMessages    = l.GetMessages
	IsDebugEnabled = l.IsDebugEnabled
	ToggleDebug    = l.ToggleDebug
)

type logLevel string

const (
	logLevelDebug logLevel = "DEBUG"
	logLevelWarn  logLevel = "WARN"
	logLevelError logLevel = "ERROR"
	logLevelInfo  logLevel = "INFO"
	logLevelFatal logLevel = "FATAL"
)

func (l logLevel) terminalColor() string {
	switch l {
	case logLevelDebug:
		return "\033[1;34m"
	case logLevelWarn:
		return "\033[1;33m"
	case logLevelError:
		return "\033[1;31m"
	case logLevelInfo:
		return "\033[1;32m"
	case logLevelFatal:
		return "\033[1;31m"
	default:
		return "\033[0m"
	}
}

func (l logLevel) string() string {
	return string(l)
}

func (l logLevel) uiColor() color.RGBA {
	switch l {
	case logLevelDebug:
		return color.RGBA{0, 0, 255, 255}
	case logLevelWarn:
		return color.RGBA{255, 255, 0, 255}
	case logLevelError:
		return color.RGBA{255, 0, 0, 255}
	default:
		return color.RGBA{0, 0, 0, 0}
	}
}

// Message represents a log entry with display properties
type Message struct {
	Level   logLevel
	Message string
	Color   color.RGBA
	Time    time.Time
}

// logger handles message logging with different severity levels
type logger struct {
	enableDebug bool
	messages    []*Message
}

// New creates a new Logger instance
func new() *logger {
	return &logger{
		enableDebug: true,
	}
}

func (l *logger) emit(m *Message) {
	m.Time = time.Now()
	m.Color = m.Level.uiColor()
	fmt.Printf("%s %s \033[0m %s\n", m.Level.terminalColor(), m.Level.string(), m.Message)
	l.messages = append(l.messages, m)
}

// Debug logs debug level messages
func (l *logger) Debug(s string, args ...interface{}) {
	if !l.enableDebug {
		return
	}
	message := fmt.Sprintf(s, args...)
	l.emit(&Message{Level: logLevelDebug, Message: message})
}

// Info logs informational messages
func (l *logger) Info(s string, args ...interface{}) {
	message := fmt.Sprintf(s, args...)
	l.emit(&Message{Level: logLevelInfo, Message: message})
}

// Warn logs warning messages
func (l *logger) Warn(s string, args ...interface{}) {
	message := fmt.Sprintf(s, args...)
	l.emit(&Message{Level: logLevelWarn, Message: message})
}

// UserMessage logs messages intended for the user
func (l *logger) UserMessage(s string, args ...interface{}) {
	message := fmt.Sprintf(s, args...)
	l.emit(&Message{Level: logLevelInfo, Message: message})
}

// Error logs error messages
func (l *logger) Error(s string, args ...interface{}) {
	message := fmt.Sprintf(s, args...)
	l.emit(&Message{Level: logLevelError, Message: message})
}

// Fatal logs fatal messages and exits the program
func (l *logger) Fatal(s string, args ...interface{}) {
	message := fmt.Sprintf(s, args...)
	l.emit(&Message{Level: logLevelFatal, Message: message})
	panic(message)
}

// SetDebug enables or disables debug logging
func (l *logger) ToggleDebug() {
	l.enableDebug = !l.enableDebug
}

func (l *logger) IsDebugEnabled() bool {
	return l.enableDebug
}

// GetMessages returns all logged messages
func (l *logger) GetMessages() []*Message {
	messages := l.messages
	l.messages = nil
	return messages
}
