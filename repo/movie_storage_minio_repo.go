package repo

import (
	"context"
	"io"
	"processor/types"
)

type movieStorageMinioRepo struct {
	client     types.ObjectStorageProvider
	bucketName string
}

func NewMovieStorageMinioRepo(client types.ObjectStorageProvider, bucketName string) types.ObjectStorage {
	client.CreateBucket(context.Background(), bucketName)
	return &movieStorageMinioRepo{client: client, bucketName: bucketName}
}

func (o *movieStorageMinioRepo) Get(ctx context.Context, name string) (io.Reader, int64, error) {
	return o.client.Get(ctx, o.bucketName, name)
}

func (o *movieStorageMinioRepo) Put(ctx context.Context, objectName string, reader io.Reader, size int64) error {
	return o.client.Put(ctx, o.bucketName, objectName, reader, size)
}
