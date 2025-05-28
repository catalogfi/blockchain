package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerOptions struct {
	webhookURL   string
	serviceName  string
	consoleLevel zapcore.Level
	webhookLevel zapcore.Level
}

func defaultLoggerOptions() *LoggerOptions {
	return &LoggerOptions{
		consoleLevel: zapcore.InfoLevel,
		webhookLevel: zapcore.ErrorLevel,
	}
}

func WithConsoleLevel(consoleLevel zapcore.Level) func(*LoggerOptions) {
	return func(logOpt *LoggerOptions) {
		logOpt.consoleLevel = consoleLevel
	}
}

func WithWebhookLevel(webhookLevel zapcore.Level) func(*LoggerOptions) {
	return func(logOpt *LoggerOptions) {
		logOpt.webhookLevel = webhookLevel
	}
}

func NewWebhookLogger(webhookURL string, ServiceName string, loggeropts ...func(*LoggerOptions)) *zap.Logger {
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
