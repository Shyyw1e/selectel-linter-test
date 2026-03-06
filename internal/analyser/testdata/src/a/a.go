package a

import (
	"log/slog"

	"go.uber.org/zap"
)

func fSlog() {
	slog.Info("starting server")
	slog.Info("Starting server")      // want "log message must start with a lowercase letter"
	slog.Info("запуск сервера")       // want "log message must be in English"
	slog.Warn("connection failed!!!") // want "log message must not contain special characters or emoji"
	slog.Debug("token validated")     // want "log message must not contain sensitive data"
}

func fZap() {
	zl := zap.NewNop()

	zl.Info("starting zap logger")
	zl.Info("Starting zap logger") // want "log message must start with a lowercase letter"
	zl.Warn("ошибка подключения")  // want "log message must be in English"
	zl.Error("failed!!!")          // want "log message must not contain special characters or emoji"
	zl.Debug("token rotated")      // want "log message must not contain sensitive data"
}
