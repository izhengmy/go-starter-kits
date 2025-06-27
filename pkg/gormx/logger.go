package gormx

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
	gLogger "gorm.io/gorm/logger"
)

const (
	logLevelInfo  = "info"
	logLevelWarn  = "warn"
	logLevelError = "error"
)

type LogConfig struct {
	SlowThreshold time.Duration `mapstructure:"slow-threshold"`
	Level         string        `mapstructure:"level"`
}

type logger struct {
	zapLogger *zap.Logger
	level     gLogger.LogLevel
	config    LogConfig
}

func newLogger(zapLogger *zap.Logger, config LogConfig) *logger {
	level := gLogger.Warn

	switch config.Level {
	case logLevelInfo:
		level = gLogger.Info
	case logLevelWarn:
		level = gLogger.Warn
	case logLevelError:
		level = gLogger.Error
	}

	return &logger{
		zapLogger: zapLogger,
		level:     level,
		config:    config,
	}
}

func (l *logger) LogMode(level gLogger.LogLevel) gLogger.Interface {
	newLogger := *l
	newLogger.level = level
	return &newLogger
}

func (l *logger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= gLogger.Info {
		l.zapLogger.Info(fmt.Sprintf(msg, data...))
	}
}

func (l *logger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= gLogger.Warn {
		l.zapLogger.Warn(fmt.Sprintf(msg, data...))
	}
}

func (l *logger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= gLogger.Error {
		l.zapLogger.Error(fmt.Sprintf(msg, data...))
	}
}

func (l *logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.level <= gLogger.Silent {
		return
	}

	slowThreshold := l.config.SlowThreshold * time.Millisecond
	times := time.Since(begin)
	sql, rows := fc()
	fields := []zap.Field{
		zap.String("sql", sql),
		zap.Int64("rows", rows),
		zap.Duration("times", times),
	}

	switch {
	case l.level >= gLogger.Error && err != nil && !errors.Is(err, gLogger.ErrRecordNotFound):
		l.zapLogger.Error("sql error", append(fields, zap.Error(err))...)
	case l.level >= gLogger.Warn && times > slowThreshold && slowThreshold != 0:
		l.zapLogger.Warn("sql slow", fields...)
	case l.level == gLogger.Info:
		l.zapLogger.Info("sql executed", fields...)
	}
}
