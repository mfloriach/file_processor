package repo

import (
	"context"
	"processor/types"
)

type mockMovieRepository struct {
}

func NewMovieMockRepo() types.MovieRepository {
	return &mockMovieRepository{}
}

func (r mockMovieRepository) Save(ctx context.Context, movie *types.Movie) error {
	return nil
}
