package boundedwait

import (
	"context"

	"golang.org/x/sync/errgroup"
)

// ErrGroup provides a errgroup with a bounded amount of concurrency.
type ErrGroup struct {
	eg    *errgroup.Group
	limit chan struct{}
}

func WithContext(ctx context.Context, limit int) (*ErrGroup, context.Context) {
	eg, ctx := errgroup.WithContext(ctx)
	return &ErrGroup{limit: make(chan struct{}, limit), eg: eg}, ctx
}

func (g *ErrGroup) Go(f func() error) {
	g.limit <- struct{}{}
	g.eg.Go(func() error {
		defer func() {
			<-g.limit
		}()
		return f()
	})
}

func (g *ErrGroup) Wait() error {
	return g.eg.Wait()
}
