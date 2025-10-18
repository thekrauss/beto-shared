// pkg/logger/logger.go
package logger

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.SugaredLogger

type Config struct {
	Level  string `mapstructure:"level"`  // "debug", "info", "warn", "error"
	Format string `mapstructure:"format"` // "json" ou "console"
}

func InitLogger(cfg Config) *zap.SugaredLogger {
	var level zapcore.Level
	switch cfg.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	var encoder zapcore.Encoder
	if cfg.Format == "console" {
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	}

	core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	log = logger.Sugar()
	return log
}

func L() *zap.SugaredLogger {
	if log == nil {
		InitLogger(Config{Level: "info", Format: "console"})
	}
	return log
}

func FromContext(ctx context.Context) *zap.SugaredLogger {
	if ctx == nil {
		return L()
	}
	if v := ctx.Value("request_id"); v != nil {
		return L().With("request_id", v)
	}
	return L()
}

func Sync() {
	if log != nil {
		_ = log.Sync()
	}
}
