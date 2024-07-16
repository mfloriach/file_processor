package services

import (
	"context"
	"errors"
	"processor/repo"
	"processor/types"
)

type metadata struct {
	config *Config
	api    types.MovieDetailsRepo
}

func NewMetadata(config *Config) types.Step {
	return &metadata{
		config: config,
		api:    repo.NewMetadataImdbApiRepo(),
	}
}

func (m *metadata) Sequencial(ctx context.Context, job types.Movie) (types.Movie, error) {
	if job.Metadata.Position == 3 {
		return job, errors.New("metadata step is not supported")
	}

	d, err := m.api.Get(ctx, job.Title)
	if err != nil {
		return types.Movie{}, err
	}

	job.Details = d

	return job, nil
}

func (m *metadata) SequencialSyncPool(ctx context.Context, job *types.Movie) error {
	d, err := m.api.Get(ctx, job.Title)
	if err != nil {
		return err
	}

	job.Details = d

	return nil
}
