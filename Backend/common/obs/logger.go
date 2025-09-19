package obs

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(service string) *zap.Logger {
	enc := zap.NewProductionEncoderConfig()
	enc.TimeKey = "ts"
	enc.EncodeTime = zapcore.ISO8601TimeEncoder
	enc.LevelKey = "level"
	enc.MessageKey = "msg"

	lvl := zapcore.InfoLevel
	if v := os.Getenv("LOG_LEVEL"); v != "" {
		_ = lvl.UnmarshalText([]byte(v))
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(enc),
		zapcore.AddSync(os.Stdout),
		lvl,
	)

	if service == "" {
		service = os.Getenv("SERVICE_NAME")
		if service == "" {
			service = "service"
		}
	}
	return zap.New(core).With(zap.String("service", service))
}
