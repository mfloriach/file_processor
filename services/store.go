package services

import (
	"context"
	"encoding/json"
	"processor/repo"
	"processor/types"
	"processor/utils"
	"time"

	"github.com/go-redis/redis/v8"
)

type store struct {
	config   *Config
	metadata types.MovieRepository
	cache    *redis.Client
}

func NewStore(config *Config) types.Step {
	return &store{
		config:   config,
		metadata: repo.NewMovieMongoRepo(),
		cache:    utils.GetRedisClient(),
	}
}

func (a store) Sequencial(ctx context.Context, m types.Movie) (types.Movie, error) {
	if err := a.metadata.Save(ctx, &m); err != nil {
		return types.Movie{}, err
	}

	b, err := json.Marshal(m)
	if err != nil {
		return types.Movie{}, err
	}

	if err := a.cache.Set(ctx, m.ID, b, time.Minute); err.Err() != nil {
		return types.Movie{}, err.Err()
	}

	return m, nil
}

func (a store) SequencialSyncPool(ctx context.Context, m *types.Movie) error {
	if err := a.metadata.Save(ctx, m); err != nil {
		return err
	}

	return nil
}
