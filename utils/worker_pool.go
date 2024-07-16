package utils

import (
	"context"
	"processor/types"
	"sync"
)

func WorkerPool(
	ctx context.Context,
	out chan types.Movie,
	in chan<- types.Movie,
	j int,
	action func(context.Context, types.Movie) (types.Movie, error),
	errs *sync.Map,
) func() {
	var wg sync.WaitGroup

	for w := 0; w <= j; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for a := range out {
				_, err := action(ctx, a)
				if err != nil {
					errs.Store(a.Metadata.Position, err)
				}

				if in == nil {
					continue
				}

				in <- a
			}
		}(w)
	}

	return func() {
		close(out)
		wg.Wait()
	}
}

func WorkerPoolSync(
	ctx context.Context,
	out chan *types.Movie,
	in chan<- *types.Movie,
	j int,
	action func(context.Context, *types.Movie) error) func() {
	var wg sync.WaitGroup

	for w := 0; w <= j; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for a := range out {
				action(ctx, a)
				if in == nil {
					continue
				}

				in <- a
			}
		}(w)
	}

	return func() {
		close(out)
		wg.Wait()
	}
}
