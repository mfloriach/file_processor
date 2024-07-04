package repo

import (
	"context"
	"processor/types"
	"processor/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const collectionName = "test"

type mongoMovieRepository struct {
	collection *mongo.Collection
}

func NewMovieMongoRepo() types.MovieRepository {
	return &mongoMovieRepository{collection: utils.GetMongoClient().Collection(collectionName)}
}

func (r mongoMovieRepository) Save(ctx context.Context, movie *types.Movie) error {
	res, err := r.collection.InsertOne(ctx, movie)

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		movie.ID = oid.String()
	}

	return err
}
