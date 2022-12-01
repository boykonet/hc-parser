package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zlog struct {
	log *zap.Logger
}

func NewLogger() ILogger {
	zapLogger, err := zap.NewProduction()
	if err != nil {
		return nil
	}
	return &zlog{
		log: zapLogger,
	}
}

func (l *zlog) Info(msg string, fields ...zapcore.Field) {
	l.log.Info(msg, fields...)
}

func (l *zlog) Error(msg string, fields ...zapcore.Field) {
	l.log.Error(msg, fields...)
}

func (l *zlog) Fatal(msg string, fields ...zapcore.Field) {
	l.log.Fatal(msg, fields...)
}

func (l *zlog) With(fields ...zapcore.Field) ILogger {
	return &zlog{
		log: l.log.With(fields...),
	}
}