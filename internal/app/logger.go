package app

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Logger provides structured logging
type Logger struct {
	infoLogger   *log.Logger
	errorLogger  *log.Logger
	actionLogger *log.Logger
	file         *os.File
}

// NewLogger creates a new logger
func NewLogger(logsDir string) (*Logger, error) {
	// Create logs directory
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create logs directory: %w", err)
	}

	// Create log file with timestamp
	timestamp := time.Now().Format("2006-01-02")
	logPath := filepath.Join(logsDir, fmt.Sprintf("linkedin-bot-%s.log", timestamp))

	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	// Create multi-writer (file + stdout)
	multiWriter := io.MultiWriter(os.Stdout, file)

	logger := &Logger{
		infoLogger:   log.New(multiWriter, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger:  log.New(multiWriter, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile),
		actionLogger: log.New(multiWriter, "[ACTION] ", log.Ldate|log.Ltime),
		file:         file,
	}

	return logger, nil
}

// Info logs an info message
func (l *Logger) Info(format string, v ...interface{}) {
	l.infoLogger.Printf(format, v...)
}

// Error logs an error message
func (l *Logger) Error(format string, v ...interface{}) {
	l.errorLogger.Printf(format, v...)
}

// Action logs an action
func (l *Logger) Action(actionType, target string, success bool) {
	status := "SUCCESS"
	if !success {
		status = "FAILED"
	}
	l.actionLogger.Printf("[%s] %s -> %s", status, actionType, target)
}

// Security logs a security event
func (l *Logger) Security(event string) {
	l.errorLogger.Printf("[SECURITY] %s", event)
}

// RateLimit logs a rate limit event
func (l *Logger) RateLimit(limitType string, reason string) {
	l.infoLogger.Printf("[RATE_LIMIT] %s: %s", limitType, reason)
}

// Close closes the logger
func (l *Logger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}
