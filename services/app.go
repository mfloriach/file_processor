package services

import (
	"context"
	"fmt"
	"processor/types"
	"processor/utils"
	"sync"

	"go.uber.org/zap"
)

type Mode string

const (
	ModeSequencial  Mode = "sequencial"
	ModeConcurrency Mode = "concurrency"
	ModeParallel    Mode = "parallel"
)

type App struct {
	Config   *Config
	Metadata types.Step
	Store    types.Step
	Compress types.Step
}

type Config struct {
	Mode   Mode
	Errs   *sync.Map
	Env    *utils.Env
	Cfg    *utils.Config
	Logger *zap.Logger
	Flags  types.Flags
}

func NewApp(f ...types.Flags) App {
	flg := types.Flags{}
	if len(f) > 0 {
		fmt.Println(f, len(f))
		flg = f[0]
	}

	config := &Config{
		Errs:   &sync.Map{},
		Env:    utils.GetEnv(),
		Cfg:    utils.GetConfig(),
		Logger: utils.GetLogger(),
		Flags:  flg,
	}

	return App{
		Config:   config,
		Metadata: NewMetadata(config),
		Store:    NewStore(config),
		Compress: NewCompress(config),
	}
}

func (app *App) Run(ctx context.Context, paths ...string) error {
	iter, err := NewIngestorFileReader(paths...)
	if err != nil {
		return err
	}

	app.Config.Logger.Info("start processing files", zap.String("mode", app.Config.Flags.Mode), zap.Int("num of worker", app.Config.Flags.Workers))

	if Mode(app.Config.Flags.Mode) == ModeSequencial {
		for iter.HasNext() {
			m := iter.Next()
			if err := app.Sequencial(ctx, m); err != nil {
				return err
			}
		}
	}

	if Mode(app.Config.Flags.Mode) == ModeConcurrency {
		in, close := app.Parallel(ctx)
		defer close()
		for iter.HasNext() {
			in <- iter.Next()
		}
	}

	if Mode(app.Config.Flags.Mode) == ModeParallel {
		in, close := app.Parallel(ctx, app.Config.Flags.Workers)
		defer close()
		for iter.HasNext() {
			in <- iter.Next()
		}
	}

	return nil
}

func (app *App) Close() {
	logger := utils.GetLogger()
	logger.Sync()
	utils.CloseClient()
}

func (app *App) Sequencial(ctx context.Context, m types.Movie) error {

	m1, err := app.Metadata.Sequencial(ctx, m)
	if err != nil {
		return err
	}

	m2, err := app.Compress.Sequencial(ctx, m1)
	if err != nil {
		return err
	}

	if _, err := app.Store.Sequencial(ctx, m2); err != nil {
		return err
	}

	return nil
}

func (app *App) Parallel(ctx context.Context, numWorker ...int) (chan types.Movie, func()) {
	w := 1
	if len(numWorker) > 0 {
		w = numWorker[0]
	}

	in := make(chan types.Movie, 50)
	c2 := make(chan types.Movie, 50)
	c3 := make(chan types.Movie, 50)

	w1 := utils.WorkerPool(ctx, in, c2, w, app.Metadata.Sequencial)
	w2 := utils.WorkerPool(ctx, c2, c3, w, app.Compress.Sequencial)
	w3 := utils.WorkerPool(ctx, c3, nil, w, app.Store.Sequencial)

	return in, func() {
		defer w3()
		defer w2()
		defer w1()
	}
}
