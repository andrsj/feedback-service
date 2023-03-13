package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/andrsj/feedback-service/pkg/logger"
)

func New() *zapLogger {
	//nolint
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		MessageKey:     "msg",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.Lock(os.Stdout),
		zapcore.InfoLevel,
	)

	logger := zap.New(core)

	return &zapLogger{
		logger: logger,
	}
}

type zapLogger struct {
	logger *zap.Logger
}

var _ logger.Logger = (*zapLogger)(nil)

func (l *zapLogger) Named(name string) logger.Logger {
	return &zapLogger{
		logger: l.logger.Named(name),
	}
}

func (l *zapLogger) Debug(message string, args logger.M) {
	l.logger.Debug(message, toFields(args)...)
}

func (l *zapLogger) Info(message string, args logger.M) {
	l.logger.Info(message, toFields(args)...)
}

func (l *zapLogger) Warn(message string, args logger.M) {
	l.logger.Warn(message, toFields(args)...)
}

func (l *zapLogger) Error(message string, args logger.M) {
	l.logger.Error(message, toFields(args)...)
}

func (l *zapLogger) Fatal(message string, args logger.M) {
	l.logger.Fatal(message, toFields(args)...)
}

func toFields(args logger.M) []zap.Field {
	fields := make([]zap.Field, len(args))
	i := 0

	for k, v := range args {
		fields[i] = zap.Any(k, v)
		i++
	}

	return fields
}
