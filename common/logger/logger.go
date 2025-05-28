package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LoggerOptions defines a function type for configuring logger options
type LoggerOptions func(*loggerOptions)

// loggerOptions holds the configuration for the logger
type loggerOptions struct {
	webhookURL   string
	serviceName  string
	consoleLevel zapcore.Level
	webhookLevel zapcore.Level
}

// defaultLoggerOptions returns a loggerOptions with default values
func defaultLoggerOptions() *loggerOptions {
	return &loggerOptions{
		consoleLevel: zapcore.InfoLevel,
		webhookLevel: zapcore.ErrorLevel,
	}
}

// with consoleLevel sets the logging level for console output
func WithConsoleLevel(consoleLevel zapcore.Level) LoggerOptions {
	return func(logOpt *loggerOptions) {
		logOpt.consoleLevel = consoleLevel
	}
}

// with WebhookLevel sets the logging level for webhook output
func WithWebhookLevel(webhookLevel zapcore.Level) LoggerOptions {
	return func(logOpt *loggerOptions) {
		logOpt.webhookLevel = webhookLevel
	}
}

// NewWebhookLogger creates a new logger that logs to both console and a webhook.
//
//	params :
//
// - webhookURL: URL to send error logs to
//
// - ServiceName: Name of the service for logging context
//
// - loggeropts: Optional configuration functions to customize logger behavior
//
// If no webhook URL is provided, it returns a logger that only logs to stdout.
func NewWebhookLogger(webhookURL string, ServiceName string, loggeropts ...func(*loggerOptions)) *zap.Logger {
	defaultOpts := defaultLoggerOptions()
	defaultOpts.webhookURL = webhookURL
	defaultOpts.serviceName = ServiceName
	if len(loggeropts) > 0 {
		for _, opt := range loggeropts {
			opt(defaultOpts)
		}
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	consoleConfig := encoderConfig
	consoleConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(consoleConfig)

	stdoutSink := zapcore.AddSync(os.Stdout)
	stdoutCore := zapcore.NewCore(consoleEncoder, stdoutSink, defaultOpts.consoleLevel)

	// If no webhook URL is provided, return logger with stdout only
	if webhookURL == "" {
		return zap.New(stdoutCore)
	}

	// Create webhook core for errors
	jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)
	webhookSink := NewWebhookWriter(webhookURL, ServiceName)
	webhookCore := zapcore.NewCore(jsonEncoder, zapcore.AddSync(webhookSink), defaultOpts.webhookLevel)

	return zap.New(zapcore.NewTee(stdoutCore, webhookCore))
}
