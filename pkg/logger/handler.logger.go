package logger

import (
	"context"
	"sync"

	"github.com/GoReactors/backend-learning/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	globalLogger Logger
	once         sync.Once
)

type zapLogger struct {
	zap *zap.Logger
}

func (l *zapLogger) Debug(msg string, fields ...Field) {
	l.zap.Debug(msg, toZapFields(fields)...)
}

func (l *zapLogger) Info(msg string, fields ...Field) {
	l.zap.Info(msg, toZapFields(fields)...)
}

func (l *zapLogger) Warn(msg string, fields ...Field) {
	l.zap.Warn(msg, toZapFields(fields)...)
}

func (l *zapLogger) Error(msg string, fields ...Field) {
	l.zap.Error(msg, toZapFields(fields)...)
}

func (l *zapLogger) Fatal(msg string, fields ...Field) {
	l.zap.Fatal(msg, toZapFields(fields)...)
}

func (l *zapLogger) InfoCtx(ctx context.Context, msg string, fields ...Field) {
	fields = append(fields, TraceFieldsFromContext(ctx)...)
	l.Info(msg, fields...)
}

func (l *zapLogger) With(fields ...Field) Logger {
	return &zapLogger{zap: l.zap.With(toZapFields(fields)...)}
}

func (l *zapLogger) Sync() error {
	return l.zap.Sync()
}

func toZapFields(fields []Field) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, f := range fields {
		zapFields[i] = zap.Any(f.Key, f.Value)
	}
	return zapFields
}

func Initialize(cfg config.LoggingConfig) {
	once.Do(func() {
		zapCfg := zap.NewProductionConfig()

		if cfg.Development {
			zapCfg = zap.NewDevelopmentConfig()
			zapCfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		}

		zapCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		zapCfg.Sampling = &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		}

		level := zapcore.InfoLevel
		_ = level.UnmarshalText([]byte(cfg.Level))
		zapCfg.Level = zap.NewAtomicLevelAt(level)

		zLogger, err := zapCfg.Build(
			zap.AddCaller(),
			zap.AddCallerSkip(1),
			zap.AddStacktrace(zapcore.ErrorLevel),
		)
		if err != nil {
			panic(err)
		}

		globalLogger = &zapLogger{zap: zLogger}
	})
}

func FirstSyncInMain(l Logger) {
	if err := l.Sync(); err != nil {
		// Don't panic if sync fails (common when writing to stdout)
		l.Error("failed to sync logger", Field{Key: "error", Value: err})
	}
}

func Global() Logger {
	return globalLogger
}
