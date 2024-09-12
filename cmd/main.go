package main

import (
	"ImageService/internal/app/api"
	"ImageService/internal/pkg/service"
	"ImageService/internal/pkg/worker"
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	port     = 8080
	imageDir = "./images" // TODO: вынести в конфиг
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx); err != nil {
		slog.Log(ctx, slog.LevelError, "app stopped: %w", err)
	}
}

func run(ctx context.Context) error {
	workers := provideWorkers(ctx)
	defer stopWorkers(ctx, workers)

	var (
		imageSVC   = service.NewImage(imageDir)
		handler    = api.NewHandler(imageSVC)
		serverHTTP = api.NewHTTP(port, handler)
	)

	go func() {
		if err := serverHTTP.ServeHTTP(); err != nil {

		}
	}()
	select {
	case <-ctx.Done():
		log.Println("shutdown in 3 seconds.")
		time.Sleep(3 * time.Second)
	}
	return nil
}

func provideWorkers(ctx context.Context) []worker.Worker {
	imageConvertor := worker.NewImageConvertor(imageDir)

	go func() {
		if err := imageConvertor.Start(ctx); err != nil {
			slog.Log(ctx, slog.LevelError, fmt.Sprintf("image convertor: %v", err))
		}
	}()

	return []worker.Worker{imageConvertor}
}

func stopWorkers(ctx context.Context, workers []worker.Worker) {
	for _, w := range workers {
		if err := w.Stop(); err != nil {
			slog.Log(ctx, slog.LevelError, fmt.Sprintf("stop workers: %v", err))
		}
	}
}
