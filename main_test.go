package main

import (
	"context"
	"fmt"
	"processor/services"
	"processor/types"
	"processor/utils"
	"sync"
	"testing"
)

func BenchmarkMainSequencial(b *testing.B) {
	ctx := context.TODO()

	app := services.NewApp()
	defer app.Close()

	for i := 0; i < b.N; i++ {
		m := types.Movie{
			Title: "test",
			Metadata: types.Metadata{
				Position: int64(i),
				Context:  ctx,
			},
		}

		app.Sequencial(ctx, m)
	}
}

func BenchmarkMainConcurrent(b *testing.B) {
	ctx := context.TODO()

	app := services.NewApp()
	defer app.Close()

	in, close := app.Parallel(ctx)
	defer close()

	for i := 0; i < b.N; i++ {
		in <- types.Movie{
			Title: "test",
			Metadata: types.Metadata{
				Position: int64(i),
				Context:  ctx,
			},
		}
	}
}

func BenchmarkMainParallel(b *testing.B) {
	ctx := context.TODO()

	app := services.NewApp()
	defer app.Close()

	for _, j := range []int{10} {
		b.ResetTimer()

		b.Run(fmt.Sprintln(j), func(b *testing.B) {
			in, close := app.Parallel(ctx, j)
			defer close()

			for i := 0; i < b.N; i++ {
				in <- types.Movie{
					Title: "test",
					Metadata: types.Metadata{
						Position: int64(i),
						Context:  context.Background(),
					},
				}
			}
		})
	}
}

func BenchmarkMainParPool(b *testing.B) {
	ctx := context.TODO()

	app := services.NewApp()
	defer app.Close()

	c1 := make(chan *types.Movie, 50)
	c2 := make(chan *types.Movie, 50)
	c3 := make(chan *types.Movie, 50)

	var pool = &sync.Pool{
		New: func() interface{} {
			return &types.Movie{}
		},
	}

	last := func(ctx context.Context, m *types.Movie) error {
		app.Store.SequencialSyncPool(ctx, m)
		pool.Put(m)
		return nil
	}

	wp1 := utils.WorkerPoolSync(ctx, c1, c2, 10, app.Metadata.SequencialSyncPool)
	wp2 := utils.WorkerPoolSync(ctx, c2, c3, 10, app.Compress.SequencialSyncPool)
	wp3 := utils.WorkerPoolSync(ctx, c3, nil, 10, last)
	defer wp3()
	defer wp2()
	defer wp1()

	for i := 0; i < b.N; i++ {
		obj := pool.Get().(*types.Movie)
		obj.Title = "test"
		obj.Metadata.Position = int64(i)
		obj.Metadata.Context = ctx

		c1 <- obj
	}
}
