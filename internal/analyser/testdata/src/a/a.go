package a

import "log/slog"

type logger struct{}

func (logger) Info(msg string, fields ...any) {}

func f() {
	slog.Info("starting server")
	slog.Info("Starting server")      // want "log message must start with a lowercase letter"
	slog.Info("запуск сервера") 	  // want "log message must be in English"
	slog.Warn("connection failed!!!") // want "log message must not contain special characters or emoji"
	slog.Debug("token validated")     // want "log message must not contain sensitive data"

	var l logger
	l.Info("password rotated") // want "log message must not contain sensitive data"
}
