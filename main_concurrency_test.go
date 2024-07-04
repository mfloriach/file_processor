package main

import (
	"context"
	"fmt"
	"processor/services"
	"processor/types"
	"processor/utils"
	"sync"
	"testing"
	"time"
)

func BenchmarkConcurrencySequencial(b *testing.B) {
	ctx := context.TODO()

	app := services.App{
		Config:   &services.Config{},
		Metadata: services.NewMock(time.Millisecond),
		Store:    services.NewMock(time.Millisecond),
		Compress: services.NewMock(time.Millisecond),
	}
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

func BenchmarkConcurrencyConcurrent(b *testing.B) {
	ctx := context.TODO()

	app := services.App{
		Config:   &services.Config{},
		Metadata: services.NewMock(time.Millisecond),
		Store:    services.NewMock(time.Millisecond),
		Compress: services.NewMock(time.Millisecond),
	}
	defer app.Close()

	in, c := app.Parallel(ctx, 10)
	defer c()

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

func BenchmarkConcurrencyParallel(b *testing.B) {
	ctx := context.TODO()

	for _, j := range []int{10} {
		b.Run(fmt.Sprintln(j), func(b *testing.B) {
			app := services.App{
				Config:   &services.Config{},
				Metadata: services.NewMock(time.Millisecond),
				Store:    services.NewMock(time.Millisecond),
				Compress: services.NewMock(time.Millisecond),
			}
			defer app.Close()

			in, c := app.Parallel(ctx, 10)
			defer c()

			b.ResetTimer()

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

func BenchmarkConcurrencyParPool(b *testing.B) {
	ctx := context.TODO()

	c1 := make(chan *types.Movie, 50)
	c2 := make(chan *types.Movie, 50)
	c3 := make(chan *types.Movie, 50)

	var pool = &sync.Pool{
		New: func() interface{} {
			return &types.Movie{}
		},
	}

	last := func(ctx context.Context, m *types.Movie) error {
		time.Sleep(time.Millisecond)
		pool.Put(m)
		return nil
	}

	bt := func(ctx context.Context, m *types.Movie) error {
		time.Sleep(time.Millisecond)
		return nil
	}

	wp1 := utils.WorkerPoolSync(ctx, c1, c2, 10, bt)
	wp2 := utils.WorkerPoolSync(ctx, c2, c3, 10, bt)
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
