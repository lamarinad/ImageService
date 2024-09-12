package worker

import (
	"ImageService/internal/config"
	"context"
	"fmt"
	"log/slog"
	"sync"

	"ImageService/internal/pkg/worker" // TODO: dix
)

func Run(ctx context.Context) error {
	var (
		wg      sync.WaitGroup
		cfg     = config.Default
		workers = provideWorkers(ctx, cfg.Worker, &wg)
	)

	select {
	case <-ctx.Done():
		stopWorkers(ctx, workers)
	}

	wg.Wait()

	return nil
}

func provideWorkers(ctx context.Context, cfg *config.Worker, wg *sync.WaitGroup) []worker.Worker {
	imageConvertor := worker.NewImageConvertor(cfg.ImageDir, cfg.FrequencyS)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := imageConvertor.Start(ctx); err != nil {
			slog.Warn(fmt.Sprintf("image convertor: %v", err))
		}
	}()

	return []worker.Worker{imageConvertor}
}

func stopWorkers(ctx context.Context, workers []worker.Worker) {
	for _, w := range workers {
		if err := w.Stop(); err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("stop worker: %v", err))
		}
	}
}
