package logger

import (
	"io"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

/**
 * logger package provides a structured logging system using logrus.
 * It initializes multiple loggers for different log levels and outputs logs to both console and files (os.Stdout and lumberjack).
 * Each logger is configured with a specific log file, maximum size, number of backups, and age.
 * The loggers are initialized only once using sync.Once to ensure thread safety.
 * The package provides functions to log messages at different levels (Info, Warn, Error, Fatal, Panic, Trace, Debug).
 * The log files are stored in the "logs" directory, and each logger has its own file with specific naming conventions.
 */

var (
	once          sync.Once
	RequestLogger *logrus.Logger
	InfoLogger    *logrus.Logger
	WarnLogger    *logrus.Logger
	ErrorLogger   *logrus.Logger
	FatalLogger   *logrus.Logger
	PanicLogger   *logrus.Logger
	TraceLogger   *logrus.Logger
	DebugLogger   *logrus.Logger
)

const (
	// Log file paths for different log levels
	REQUEST_LOG_FILE = "logs/request.log"
	INFO_LOG_FILE    = "logs/info.log"
	WARN_LOG_FILE    = "logs/warn.log"
	ERROR_LOG_FILE   = "logs/error.log"
	FATAL_LOG_FILE   = "logs/fatal.log"
	PANIC_LOG_FILE   = "logs/panic.log"
	TRACE_LOG_FILE   = "logs/trace.log"
	DEBUG_LOG_FILE   = "logs/debug.log"

	// Log file size: maximum size of each log file before it is rotated
	// Sizes are defined in megabytes (MB)
	REQUEST_LOG_SIZE = 100
	INFO_LOG_SIZE    = 50
	WARN_LOG_SIZE    = 20
	ERROR_LOG_SIZE   = 20
	FATAL_LOG_SIZE   = 10
	PANIC_LOG_SIZE   = 10
	TRACE_LOG_SIZE   = 30
	DEBUG_LOG_SIZE   = 30

	// Log backup: number of backups to keep for each log file
	// This is used to control how many rotated log files are kept
	// before they are deleted
	REQUEST_LOG_BACKUPS = 7
	INFO_LOG_BACKUPS    = 5
	WARN_LOG_BACKUPS    = 10
	ERROR_LOG_BACKUPS   = 15
	FATAL_LOG_BACKUPS   = 10
	PANIC_LOG_BACKUPS   = 10
	TRACE_LOG_BACKUPS   = 3
	DEBUG_LOG_BACKUPS   = 5

	// Log age: maximum age in days for each log file before it is deleted
	// This is used to control how long logs are kept before being rotated out
	REQUEST_LOG_AGE = 7
	INFO_LOG_AGE    = 14
	WARN_LOG_AGE    = 30
	ERROR_LOG_AGE   = 90
	FATAL_LOG_AGE   = 180
	PANIC_LOG_AGE   = 180
	TRACE_LOG_AGE   = 3
	DEBUG_LOG_AGE   = 7

	// Log compression: whether to compress old log files
	// This is set to true to save disk space by compressing rotated log files
	CompressLogs = true
)

func Init() {
	once.Do(func() {
		// Using TextFormatter for log formatting
		// This allows for more human-readable logs
		formatter := &logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
		}

		// Initialize all loggers with the same formatter
		// This ensures that all loggers use the same format for consistency
		RequestLogger = GetRequestLogger(formatter)
		InfoLogger = GetInfoLogger(formatter)
		WarnLogger = GetWarnLogger(formatter)
		ErrorLogger = GetErrorLogger(formatter)
		FatalLogger = GetFatalLogger(formatter)
		PanicLogger = GetPanicLogger(formatter)
		TraceLogger = GetTraceLogger(formatter)
		DebugLogger = GetDebugLogger(formatter)
	})
}

func GetRequestLogger(formatter *logrus.TextFormatter) *logrus.Logger {
	// Create a new logger for request logging
	RequestLogger = logrus.New()
	RequestLogger.SetFormatter(formatter)
	RequestLogger.SetLevel(logrus.InfoLevel)
	RequestLogger.SetOutput(io.MultiWriter(os.Stdout, &lumberjack.Logger{
		Filename:   REQUEST_LOG_FILE,
		MaxSize:    REQUEST_LOG_SIZE,
		MaxBackups: REQUEST_LOG_BACKUPS,
		MaxAge:     REQUEST_LOG_AGE,
		Compress:   CompressLogs,
	}))

	return RequestLogger
}

func GetInfoLogger(formatter *logrus.TextFormatter) *logrus.Logger {
	// Create a new logger for info logging
	InfoLogger = logrus.New()
	InfoLogger.SetFormatter(formatter)
	InfoLogger.SetLevel(logrus.InfoLevel)
	InfoLogger.SetOutput(io.MultiWriter(os.Stdout, &lumberjack.Logger{
		Filename:   INFO_LOG_FILE,
		MaxSize:    INFO_LOG_SIZE,
		MaxBackups: INFO_LOG_BACKUPS,
		MaxAge:     INFO_LOG_AGE,
		Compress:   CompressLogs,
	}))

	return InfoLogger
}

func GetWarnLogger(formatter *logrus.TextFormatter) *logrus.Logger {
	// Create a new logger for warn logging
	WarnLogger = logrus.New()
	WarnLogger.SetFormatter(formatter)
	WarnLogger.SetLevel(logrus.WarnLevel)
	WarnLogger.SetOutput(io.MultiWriter(os.Stdout, &lumberjack.Logger{
		Filename:   WARN_LOG_FILE,
		MaxSize:    WARN_LOG_SIZE,
		MaxBackups: WARN_LOG_BACKUPS,
		MaxAge:     WARN_LOG_AGE,
		Compress:   CompressLogs,
	}))

	return WarnLogger
}

func GetErrorLogger(formatter *logrus.TextFormatter) *logrus.Logger {
	// Create a new logger for error logging
	ErrorLogger = logrus.New()
	ErrorLogger.SetFormatter(formatter)
	ErrorLogger.SetLevel(logrus.ErrorLevel)
	ErrorLogger.SetOutput(io.MultiWriter(os.Stdout, &lumberjack.Logger{
		Filename:   ERROR_LOG_FILE,
		MaxSize:    ERROR_LOG_SIZE,
		MaxBackups: ERROR_LOG_BACKUPS,
		MaxAge:     ERROR_LOG_AGE,
		Compress:   CompressLogs,
	}))

	return ErrorLogger
}

func GetFatalLogger(formatter *logrus.TextFormatter) *logrus.Logger {
	// Create a new logger for fatal logging
	FatalLogger = logrus.New()
	FatalLogger.SetFormatter(formatter)
	FatalLogger.SetLevel(logrus.FatalLevel)
	FatalLogger.SetOutput(io.MultiWriter(os.Stdout, &lumberjack.Logger{
		Filename:   FATAL_LOG_FILE,
		MaxSize:    FATAL_LOG_SIZE,
		MaxBackups: FATAL_LOG_BACKUPS,
		MaxAge:     FATAL_LOG_AGE,
		Compress:   CompressLogs,
	}))

	return FatalLogger
}

func GetPanicLogger(formatter *logrus.TextFormatter) *logrus.Logger {
	// Create a new logger for panic logging
	PanicLogger = logrus.New()
	PanicLogger.SetFormatter(formatter)
	PanicLogger.SetLevel(logrus.PanicLevel)
	PanicLogger.SetOutput(io.MultiWriter(os.Stdout, &lumberjack.Logger{
		Filename:   PANIC_LOG_FILE,
		MaxSize:    PANIC_LOG_SIZE,
		MaxBackups: PANIC_LOG_BACKUPS,
		MaxAge:     PANIC_LOG_AGE,
		Compress:   CompressLogs,
	}))

	return PanicLogger
}

func GetTraceLogger(formatter *logrus.TextFormatter) *logrus.Logger {
	// Create a new logger for trace logging
	TraceLogger = logrus.New()
	TraceLogger.SetFormatter(formatter)
	TraceLogger.SetLevel(logrus.TraceLevel)
	TraceLogger.SetOutput(io.MultiWriter(os.Stdout, &lumberjack.Logger{
		Filename:   TRACE_LOG_FILE,
		MaxSize:    TRACE_LOG_SIZE,
		MaxBackups: TRACE_LOG_BACKUPS,
		MaxAge:     TRACE_LOG_AGE,
		Compress:   CompressLogs,
	}))

	return TraceLogger
}

func GetDebugLogger(formatter *logrus.TextFormatter) *logrus.Logger {
	// Create a new logger for debug logging
	DebugLogger = logrus.New()
	DebugLogger.SetFormatter(formatter)
	DebugLogger.SetLevel(logrus.DebugLevel)
	DebugLogger.SetOutput(io.MultiWriter(os.Stdout, &lumberjack.Logger{
		Filename:   DEBUG_LOG_FILE,
		MaxSize:    DEBUG_LOG_SIZE,
		MaxBackups: DEBUG_LOG_BACKUPS,
		MaxAge:     DEBUG_LOG_AGE,
		Compress:   CompressLogs,
	}))

	return DebugLogger
}

// GetLogger returns a singleton instance of logrus.Logger
func GetLogger(level logrus.Level) *logrus.Logger {
	if RequestLogger == nil || InfoLogger == nil ||
		WarnLogger == nil || ErrorLogger == nil ||
		FatalLogger == nil || PanicLogger == nil ||
		TraceLogger == nil || DebugLogger == nil {
		// Initialize the loggers if they are not already initialized
		// This ensures that the loggers are only initialized once
		Init()
	}

	// Set the log level for the logger
	switch level {
	case logrus.InfoLevel:
		return InfoLogger
	case logrus.WarnLevel:
		return WarnLogger
	case logrus.ErrorLevel:
		return ErrorLogger
	case logrus.FatalLevel:
		return FatalLogger
	case logrus.PanicLevel:
		return PanicLogger
	case logrus.TraceLevel:
		return TraceLogger
	default:
		return DebugLogger
	}
}

// Log functions for different log levels
func Info(msg string, fields logrus.Fields) {
	logger := GetLogger(logrus.InfoLevel)
	if fields != nil {
		logger.WithFields(fields).Info(msg)
	} else {
		logger.Info(msg)
	}
}

func Warn(msg string, fields logrus.Fields) {
	logger := GetLogger(logrus.WarnLevel)
	if fields != nil {
		logger.WithFields(fields).Warn(msg)
	} else {
		logger.Warn(msg)
	}
}

func Error(msg string, fields logrus.Fields) {
	logger := GetLogger(logrus.ErrorLevel)
	if fields != nil {
		logger.WithFields(fields).Error(msg)
	} else {
		logger.Error(msg)
	}
}

func Fatal(msg string, fields logrus.Fields) {
	logger := GetLogger(logrus.FatalLevel)
	if fields != nil {
		logger.WithFields(fields).Fatal(msg)
	} else {
		logger.Fatal(msg)
	}
}

func Panic(msg string, fields logrus.Fields) {
	logger := GetLogger(logrus.PanicLevel)
	if fields != nil {
		logger.WithFields(fields).Panic(msg)
	} else {
		logger.Panic(msg)
	}
}

func Trace(msg string, fields logrus.Fields) {
	logger := GetLogger(logrus.TraceLevel)
	if fields != nil {
		logger.WithFields(fields).Trace(msg)
	} else {
		logger.Trace(msg)
	}
}

func Debug(msg string, fields logrus.Fields) {
	logger := GetLogger(logrus.DebugLevel)
	if fields != nil {
		logger.WithFields(fields).Debug(msg)
	} else {
		logger.Debug(msg)
	}
}

func Exit() {
	// Exit all loggers gracefully
	RequestLogger.Exit(0)
	InfoLogger.Exit(0)
	WarnLogger.Exit(0)
	ErrorLogger.Exit(0)
	FatalLogger.Exit(0)
	PanicLogger.Exit(0)
	TraceLogger.Exit(0)
	DebugLogger.Exit(0)

	// Clean up loggers
	if RequestLogger != nil {
		RequestLogger = nil
	}
	if InfoLogger != nil {
		InfoLogger = nil
	}
	if WarnLogger != nil {
		WarnLogger = nil
	}
	if ErrorLogger != nil {
		ErrorLogger = nil
	}
	if FatalLogger != nil {
		FatalLogger = nil
	}
	if PanicLogger != nil {
		PanicLogger = nil
	}
	if TraceLogger != nil {
		TraceLogger = nil
	}
	if DebugLogger != nil {
		DebugLogger = nil
	}
}
