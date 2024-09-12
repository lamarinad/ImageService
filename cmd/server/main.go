package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"ImageService/internal/app/server"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := server.Run(ctx); err != nil {
		slog.ErrorContext(ctx, "app stopped with error: %w", err)

		return
	}

	slog.InfoContext(ctx, "app stopped")
}
