package main

import (
	"context"
	"flag"
	"processor/services"
	"processor/types"
	"processor/utils"

	"go.uber.org/zap"
)

func run(ctx context.Context, f types.Flags) {
	app := services.NewApp(f)
	defer app.Close()

	if err := app.Run(ctx, "test.txt"); err != nil {
		return
	}
}

func getArgs() types.Flags {
	mode := flag.String("mode", "parallel", "parallel, concurrent or sequential mode")
	workers := flag.Int("workers", 1, "number of workers")
	flag.Parse()

	return types.Flags{
		Mode:    *mode,
		Workers: *workers,
	}
}

func main() {
	logger := utils.GetLogger()

	t, err := utils.AddTracers()
	if err != nil {
		logger.Error("failed to add tracer", zap.Error(err))
	}
	defer t()

	f := getArgs()

	run(context.TODO(), f)
}
