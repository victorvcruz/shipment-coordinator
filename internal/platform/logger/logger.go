package logger

import (
	"os"
	"path/filepath"

	"github.com/victorvcruz/shipment-coordinator/internal/platform/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(cfg *config.AppConfig) (*zap.Logger, error) {
	hostname, _ := os.Hostname()

	logDir := "./var/log/shipment-coordinator"
	err := os.MkdirAll(logDir, 0o755)
	if err != nil {
		return nil, err
	}
	logFilePath := filepath.Join(logDir, "app.log")

	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "ts"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.LevelKey = "level"
	encoderCfg.EncodeLevel = zapcore.LowercaseLevelEncoder

	level, err := zapcore.ParseLevel(cfg.Logging.Level)
	if err != nil {
		return nil, err
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(os.Stdout),
			zapcore.AddSync(file),
		),
		level,
	)

	logger := zap.New(core).With(
		zap.String("host.name", hostname),
		zap.String("env", cfg.Env),
		zap.String("service", cfg.Service),
		zap.String("version", cfg.Version))

	zap.ReplaceGlobals(logger)

	return logger, nil
}
