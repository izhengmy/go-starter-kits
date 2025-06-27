package zapx

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	levelDebug  = "debug"
	levelInfo   = "info"
	levelWarn   = "warn"
	levelError  = "error"
	levelDPanic = "dpanic"
	levelPanic  = "panic"
	levelFatal  = "fatal"
)

type Config struct {
	Path       string `mapstructure:"path"`
	Level      string `mapstructure:"level"`
	MaxSize    int    `mapstructure:"max-size"`
	MaxBackups int    `mapstructure:"max-backups"`
	MaxAge     int    `mapstructure:"max-age"`
}

func NewLogger(config Config) (*zap.Logger, func()) {
	rotator := &lumberjack.Logger{
		Filename:   dailyPath(config.Path),
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   true,
		LocalTime:  true,
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.MessageKey = "message"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.DateTime)
	jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)

	core := zapcore.NewCore(
		jsonEncoder,
		zapcore.AddSync(rotator),
		zap.NewAtomicLevelAt(transLevel(config.Level)),
	)

	stackTraceLevel := zap.NewAtomicLevelAt(zap.ErrorLevel)
	logger := zap.New(
		core,
		zap.AddCaller(),
		zap.AddStacktrace(stackTraceLevel),
		zap.AddCallerSkip(1),
	)

	return logger, func() {
		_ = logger.Sync()
	}
}

func transLevel(level string) zapcore.Level {
	switch level {
	case levelDebug:
		return zap.DebugLevel
	case levelInfo:
		return zap.InfoLevel
	case levelWarn:
		return zap.WarnLevel
	case levelError:
		return zap.ErrorLevel
	case levelDPanic:
		return zap.DPanicLevel
	case levelPanic:
		return zap.PanicLevel
	case levelFatal:
		return zap.FatalLevel
	default:
		return zap.InfoLevel
	}
}

func dailyPath(path string) string {
	date := time.Now().Format("2006-01-02")
	dir := filepath.Dir(path)
	filename := filepath.Base(path)
	ext := filepath.Ext(path)
	name := strings.TrimSuffix(filename, ext)
	newFilename := fmt.Sprintf("%s-%s%s", name, date, ext)
	return filepath.Join(dir, newFilename)
}
