package services

import (
	"context"
	"processor/repo"
	"processor/types"
)

type store struct {
	config   *Config
	metadata types.MovieRepository
}

func NewStore(config *Config) types.Step {
	return &store{
		config:   config,
		metadata: repo.NewMovieMongoRepo(),
	}
}

func (a store) Sequencial(ctx context.Context, m types.Movie) (types.Movie, error) {
	if err := a.metadata.Save(ctx, &m); err != nil {
		return types.Movie{}, err
	}

	return m, nil
}

func (a store) SequencialSyncPool(ctx context.Context, m *types.Movie) error {
	if err := a.metadata.Save(ctx, m); err != nil {
		return err
	}

	return nil
}
