package logger

import (
	"go.uber.org/zap"
)

var Log *zap.Logger

func InitLogger(env string) {
	var err error

	if env == "prod" {
		Log, err = zap.NewProduction()
	} else if env == "dev" {
		Log, err = zap.NewDevelopment()
	} else {
		Log = zap.NewNop()
	}

	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}

	Log = Log.With(zap.String("app", "orders-app"), zap.String("env", env))
}

func SyncLogger() {
	if Log != nil {
		_ = Log.Sync()
	}
}
