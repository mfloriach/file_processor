package services

import (
	"context"
	"processor/types"
	"time"
)

type mock struct {
	t time.Duration
}

func NewMock(d time.Duration) types.Step {
	return &mock{
		t: d,
	}
}

func (mk mock) Sequencial(ctx context.Context, m types.Movie) (types.Movie, error) {
	time.Sleep(mk.t)
	return m, nil
}

func (mk mock) SequencialSyncPool(ctx context.Context, m *types.Movie) error {
	time.Sleep(mk.t)
	return nil
}
