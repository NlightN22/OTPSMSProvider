package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.SugaredLogger

// Init initializes the global logger with module name and timestamp.
func New(module string) (*zap.SugaredLogger, error) {
	cfg := zap.NewProductionConfig()
	cfg.InitialFields = map[string]interface{}{"module": module}
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	base, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	return base.Sugar(), nil
}
