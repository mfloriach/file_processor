package utils

import (
	"context"
	"processor/types"
	"sync"
	"time"
)

func WorkerPool(
	ctx context.Context,
	out chan types.Movie,
	in chan<- types.Movie,
	j int,
	action func(context.Context, types.Movie) (types.Movie, error)) func() {
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

type workerPool struct {
	out    chan *types.Movie
	in     chan<- *types.Movie
	action func(context.Context, *types.Movie) error
	wg     sync.WaitGroup
	done   int
	mx     sync.Mutex
}

type Worker interface {
	Close()
	Add(ctx context.Context, numWorker int)
	Remove()
}

func NewWorkerPool(out chan *types.Movie, in chan<- *types.Movie, action func(context.Context, *types.Movie) error) Worker {
	return &workerPool{
		out:    out,
		in:     in,
		action: action,
		wg:     sync.WaitGroup{},
		done:   0,
		mx:     sync.Mutex{},
	}
}

func (wp *workerPool) Close() {
	close(wp.out)
	wp.wg.Wait()
}

func (wp *workerPool) Add(ctx context.Context, numWorker int) {
	for w := 0; w <= numWorker; w++ {
		wp.wg.Add(1)
		go func(workerID int) {
			defer wp.wg.Done()
			for a := range wp.out {
				wp.action(ctx, a)
				if wp.in == nil {
					continue
				}

				wp.in <- a

				wp.mx.Lock()
				defer wp.mx.Unlock()
				if wp.done != 0 {
					wp.done--
					return
				}
			}
		}(w)
	}
}

func (wp *workerPool) Remove() {
	wp.done++
	ticker := time.NewTicker(500 * time.Millisecond)
	for {
		if wp.done == 0 {
			return
		}
		<-ticker.C
	}
}
