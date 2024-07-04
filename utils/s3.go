package utils

import (
	"context"
	"fmt"
	"io"
	"log"
	"processor/types"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type minioObjStorage struct {
	client *minio.Client
	env    *Env
}

func NewMinioObjectStorage() types.ObjectStorageProvider {
	env := GetEnv()

	client, err := minio.New(env.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(env.MinioAccessKey, env.MinioSecretKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
	}

	return &minioObjStorage{client: client, env: env}
}

func (m *minioObjStorage) CreateBucket(ctx context.Context, name string) {
	if err := m.client.MakeBucket(ctx, name, minio.MakeBucketOptions{Region: m.env.MinioLocation}); err != nil {
		exists, errBucketExists := m.client.BucketExists(ctx, name)
		if errBucketExists == nil && exists {
			// log.Printf("We already own %s\n", name)
		} else {
			log.Fatalln(err)
		}
		return
	}

	log.Printf("Successfully created %s\n", name)
}

func (m *minioObjStorage) List(ctx context.Context, bucketName string) {
	for object := range m.client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		UseV1:     true,
		Recursive: true,
	}) {
		if object.Err != nil {
			fmt.Println(object.Err)
		}
		fmt.Println(object)
	}
}

func (m *minioObjStorage) Get(ctx context.Context, bucketName string, objectName string) (io.Reader, int64, error) {
	reader, err := m.client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, 0, err
	}

	stat, err := reader.Stat()
	if err != nil {
		return nil, 0, err
	}

	return reader, stat.Size, nil
}

func (m *minioObjStorage) Put(ctx context.Context, bucketName string, objectName string, reader io.Reader, size int64) error {
	if _, err := m.client.PutObject(ctx, m.env.MinioBucketNameOut, objectName, reader, size, minio.PutObjectOptions{ContentType: "application/octet-stream"}); err != nil {
		return err
	}

	return nil
}
