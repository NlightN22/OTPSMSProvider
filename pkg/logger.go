package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.SugaredLogger

func Init(level string) error {
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(level)); err != nil {
		return fmt.Errorf("invalid log level %q: %w", level, err)
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(zapLevel)
	logger, err := cfg.Build(zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	if err != nil {
		return fmt.Errorf("build logger: %w", err)
	}
	Log = logger.Sugar()
	zap.ReplaceGlobals(logger)

	return nil
}

func New(module string) *zap.SugaredLogger {
	return zap.L().
		With(zap.String("module", module)).
		Sugar()
}
