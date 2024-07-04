package types

import (
	"context"
	"io"
	"time"
)

type Movie struct {
	ID          string       `json:"id"`
	Title       string       `json:"title"`
	Year        string       `json:"year"`
	ReleaseDate string       `json:"release_date"`
	Details     MovieDetails `json:"details"`
	Metadata    Metadata     `json:"-"`
}

type MovieDetails struct {
	Description string   `json:"description"`
	Stars       float32  `json:"start"`
	Actors      []string `json:"actors"`
}

type Metadata struct {
	Position        int64           `json:"position"`
	StartProcessing time.Time       `json:"start_processing"`
	FileName        string          `json:"file_name"`
	Context         context.Context `json:"-"`
}

type MovieDetailsRepo interface {
	Get(ctx context.Context, title string) (MovieDetails, error)
}

type MovieRepository interface {
	Save(context.Context, *Movie) error
}

type Iterator interface {
	HasNext() bool
	Next() Movie
}

type Injestor interface {
	Run(Iterator)
}

type Step interface {
	Sequencial(context.Context, Movie) (Movie, error)
	SequencialSyncPool(context.Context, *Movie) error
}

type ObjectStorageProvider interface {
	List(ctx context.Context, bucketName string)
	CreateBucket(ctx context.Context, name string)
	Get(ctx context.Context, bucketName string, objectName string) (io.Reader, int64, error)
	Put(ctx context.Context, bucketName string, objectName string, reader io.Reader, size int64) error
}

type ObjectStorage interface {
	Get(ctx context.Context, objectName string) (io.Reader, int64, error)
	Put(ctx context.Context, objectName string, reader io.Reader, size int64) error
}

type Flags struct {
	Mode    string
	Workers int
}
