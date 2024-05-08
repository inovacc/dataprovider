package logger

import "io"

const (
	dateFormat = "2006-01-02T15:04:05.000" // YYYY-MM-DDTHH:MM:SS.ZZZ
)

// LogLevel defines log levels.
type LogLevel uint8

// defines our own log levels, just in case we'll change logger in future
const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
)

// Logger is the interface that wraps the basic logging methods.
type Logger interface {
	// Debug logs a message at debug level.
	Debug(msg string, fields ...any)

	// Info logs a message at info level.
	Info(msg string, fields ...any)

	// Warn logs a message at warn level.
	Warn(msg string, fields ...any)

	// Error logs a message at error level.
	Error(msg string, fields ...any)

	// WithFields returns a new Logger with the given fields.
	WithFields(fields ...any) Logger

	// WithError returns a new Logger with the given error.
	WithError(err error) Logger

	// WithField returns a new Logger with the given key and value.
	WithField(key string, value any) Logger
}

type logger struct {
	level    LogLevel
	formater any
	out      io.Writer
}

// NewLogger creates a new Logger with the given name.
func NewLogger(name string) Logger {
	return newLogger(name)
}

// SetLevel sets the log level for the default logger.
func SetLevel(level LogLevel) {
	defaultLogger.SetLevel(level)
}

// SetOutput sets the output for the default logger.
func SetOutput(out io.Writer) {
	defaultLogger.SetOutput(out)
}

// SetFormatter sets the formatter for the default logger.
func SetFormatter(formatter Formatter) {
	defaultLogger.SetFormatter(formatter)
}

// SetLevel sets the log level for the default logger.
func (l *logger) SetLevel(level LogLevel) {
	l.level = level
}

// SetOutput sets the output for the default logger.
func (l *logger) SetOutput(out io.Writer) {
	l.out = out
}

// SetFormatter sets the formatter for the default logger.
func (l *logger) SetFormatter(formatter Formatter) {
	l.formatter = formatter
}

// WithFields returns a new Logger with the given fields.
func (l *logger) WithFields(fields ...any) Logger {
	return l.withFields(fields...)
}

// WithError returns a new Logger with the given error.
func (l *logger) WithError(err error) Logger {
	return l.withError(err)
}

// WithField returns a new Logger with the given key and value.
func (l *logger) WithField(key string, value any) Logger {
	return l.withField(key, value)
}

// Debug logs a message at debug level.
func Debug(msg string, fields ...any) {
	defaultLogger.Debug(msg, fields...)
}
