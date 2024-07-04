package services

import (
	"context"
	"processor/repo"
	"processor/types"
	"processor/utils"
)

type compress struct {
	config          *Config
	movieStorageSrc types.ObjectStorage
	movieStorageOut types.ObjectStorage
}

func NewCompress(config *Config) types.Step {
	client := utils.NewMinioObjectStorage()

	env := utils.GetEnv()

	return &compress{
		config:          config,
		movieStorageSrc: repo.NewMovieStorageMinioRepo(client, env.MinioBucketNameSrc),
		movieStorageOut: repo.NewMovieStorageMinioRepo(client, env.MinioBucketNameOut),
	}
}

func (c *compress) Sequencial(ctx context.Context, job types.Movie) (types.Movie, error) {
	reader, size, err := c.movieStorageSrc.Get(ctx, "test.mp4")
	if err != nil {
		return types.Movie{}, err
	}
	// defer reader.Close()

	// out := &bytes.Buffer{}
	// if err := ffmpeg.Input("pipe:", ffmpeg.KwArgs{}).WithInput(reader).Output("pipe:", ffmpeg.KwArgs{
	// 	"acodec": "pcm_s16le",
	// 	"f":      "wav",
	// 	"ac":     "1",
	// 	"ar":     "16000",
	// }).WithOutput(out).OverWriteOutput().Silent(true).Run(); err != nil {
	// 	panic(err)
	// 	return types.Movie{}, err
	// }

	if err := c.movieStorageOut.Put(ctx, job.Title+"test.mp4", reader, size); err != nil {
		return types.Movie{}, err
	}

	return job, nil
}

func (c *compress) SequencialSyncPool(ctx context.Context, job *types.Movie) error {
	reader, size, err := c.movieStorageSrc.Get(ctx, "architecture_high_level.svg")
	if err != nil {
		return err
	}
	// defer reader.Close()

	// in := new(bytes.Buffer)
	// in.ReadFrom(reader)

	// out := &bytes.Buffer{}
	// if err := ffmpeg.Input("pipe:0", ffmpeg.KwArgs{}).WithInput(in).Output("pipe:1", ffmpeg.KwArgs{
	// 	"acodec": "pcm_s16le",
	// 	"f":      "wav",
	// 	"ac":     "1",
	// 	"ar":     "16000",
	// }).WithOutput(out).OverWriteOutput().Run(); err != nil {
	// 	fmt.Println("error", err)
	// 	return
	// }

	if err := c.movieStorageOut.Put(ctx, job.Title+"architecture_high_level.svg", reader, size); err != nil {
		return err
	}

	return nil
}
