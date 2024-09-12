package worker

import "context"

type Worker interface {
	Start(ctx context.Context) error
	Stop() error
}
