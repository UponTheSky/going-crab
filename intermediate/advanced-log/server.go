package main

import (
	"log/slog"
	"os"
)

func main() {
	writer := os.Stdout
	handlerOptions := slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug}
	handler := slog.NewTextHandler(writer, &handlerOptions)
	logger := slog.New(handler)

	logger.Info("hey", "Mr.Chow", "Ken Jeong")
}
