package main

import (
	"context"
	"log/slog"
	"os"
)

func main() {
	// initialize a text handler, and a logger
	writer := os.Stdout
	handlerOptions := slog.HandlerOptions{AddSource: false, Level: slog.LevelDebug}
	handler := slog.NewTextHandler(writer, &handlerOptions)
	logger := slog.New(handler)

	// logger.Info("hey", "Mr.Chow", "Ken Jeong")
	loggerCtx := context.Background()
	logger.LogAttrs(loggerCtx, slog.LevelInfo, "hey", slog.String("MrChow", "Ken Jeong"))

	hangoverLogger := logger.WithGroup("hangover")
	hangoverLogger.LogAttrs(loggerCtx, slog.LevelInfo, "hey", slog.String("MrChow", "Ken Jeong"))

	withLogger := hangoverLogger.With("Alan", "Zach Galifianakis")
	withLogger.LogAttrs(loggerCtx, slog.LevelInfo, "hey", slog.String("MrChow", "Ken Jeong"))
}
