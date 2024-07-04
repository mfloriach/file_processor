package main

import (
	"context"
	"errors"
	"processor/types"
	"processor/utils"
	"sync"
	"testing"
	"time"

	cmap "github.com/orcaman/concurrent-map/v2"
)

const w = 20

func BenchmarkErrorsLock(b *testing.B) {
	ctx := context.TODO()

	c1 := make(chan *types.Movie, 50)
	c2 := make(chan *types.Movie, 50)
	c3 := make(chan *types.Movie, 50)
	errs := map[int64]error{}
	mx := sync.Mutex{}

	var pool = &sync.Pool{
		New: func() interface{} {
			return &types.Movie{}
		},
	}

	last := func(ctx context.Context, m *types.Movie) error {
		time.Sleep(time.Millisecond)
		pool.Put(m)
		if m.Metadata.Position%2 == 0 {
			mx.Lock()
			errs[m.Metadata.Position] = errors.New("fgfdgfdgfdgfd")
			mx.Unlock()
		}

		return nil
	}

	bt := func(ctx context.Context, m *types.Movie) error {
		time.Sleep(time.Millisecond)
		return nil
	}

	wp1 := utils.WorkerPoolSync(ctx, c1, c2, w, bt)
	wp2 := utils.WorkerPoolSync(ctx, c2, c3, w, bt)
	wp3 := utils.WorkerPoolSync(ctx, c3, nil, w, last)
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

func BenchmarkErrorsSyncMap(b *testing.B) {
	ctx := context.TODO()

	c1 := make(chan *types.Movie, 50)
	c2 := make(chan *types.Movie, 50)
	c3 := make(chan *types.Movie, 50)
	errs := sync.Map{}

	var pool = &sync.Pool{
		New: func() interface{} {
			return &types.Movie{}
		},
	}

	last := func(ctx context.Context, m *types.Movie) error {
		time.Sleep(time.Millisecond)
		pool.Put(m)
		if m.Metadata.Position%2 == 0 {
			errs.Store(m.ID, errors.New("fdgfdgfdgfd"))
		}

		return nil
	}

	bt := func(ctx context.Context, m *types.Movie) error {
		time.Sleep(time.Millisecond)
		return nil
	}

	wp1 := utils.WorkerPoolSync(ctx, c1, c2, w, bt)
	wp2 := utils.WorkerPoolSync(ctx, c2, c3, w, bt)
	wp3 := utils.WorkerPoolSync(ctx, c3, nil, w, last)
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

func BenchmarkErrorsConMap(b *testing.B) {
	ctx := context.TODO()

	c1 := make(chan *types.Movie, 50)
	c2 := make(chan *types.Movie, 50)
	c3 := make(chan *types.Movie, 50)
	errs := cmap.New[string]()

	var pool = &sync.Pool{
		New: func() interface{} {
			return &types.Movie{}
		},
	}

	last := func(ctx context.Context, m *types.Movie) error {
		time.Sleep(time.Millisecond)
		pool.Put(m)
		if m.Metadata.Position%2 == 0 {
			errs.Set(m.ID, "fdgfdgfdgfd")
		}

		return nil
	}

	bt := func(ctx context.Context, m *types.Movie) error {
		time.Sleep(time.Millisecond)
		return nil
	}

	wp1 := utils.WorkerPoolSync(ctx, c1, c2, w, bt)
	wp2 := utils.WorkerPoolSync(ctx, c2, c3, w, bt)
	wp3 := utils.WorkerPoolSync(ctx, c3, nil, w, last)
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

func BenchmarkErrorsChannel(b *testing.B) {
	ctx := context.TODO()

	c1 := make(chan *types.Movie, 50)
	c2 := make(chan *types.Movie, 50)
	c3 := make(chan *types.Movie, 50)
	errs := make(chan error, 50)
	a := map[int64]error{}

	var pool = &sync.Pool{
		New: func() interface{} {
			return &types.Movie{}
		},
	}

	last := func(ctx context.Context, m *types.Movie) error {
		time.Sleep(time.Millisecond)
		pool.Put(m)
		if m.Metadata.Position%2 == 0 {
			errs <- errors.New("dfdsfdsfds")
		}

		return nil
	}

	bt := func(ctx context.Context, m *types.Movie) error {
		time.Sleep(time.Millisecond)
		return nil
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		wg.Done()
		var i int64 = 0
		for e := range errs {
			a[i] = e
			i++
		}
	}()

	defer close(errs)
	defer wg.Wait()

	wp1 := utils.WorkerPoolSync(ctx, c1, c2, w, bt)
	wp2 := utils.WorkerPoolSync(ctx, c2, c3, w, bt)
	wp3 := utils.WorkerPoolSync(ctx, c3, nil, w, last)
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
