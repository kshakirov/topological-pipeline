//go:build debug

package main

import (
	"log/slog"
	"os"
)

var logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
	AddSource: true,            // Выводит файл и строчку (main.go:42)
	Level:     slog.LevelDebug, // Включаем уровень DEBUG
}))

// Debug — выводит логи только при флаге -tags debug
func Debug(msg string, args ...any) {
	logger.Debug(msg, args...)
}
