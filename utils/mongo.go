package utils

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var client *mongo.Client

func GetMongoClient() *mongo.Database {
	env := GetEnv()

	if client != nil {
		return client.Database(env.DBName)
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(env.DBConn))
	if err != nil {
		panic(err)
	}

	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	return client.Database(env.DBName)
}

func CloseClient() {
	if client != nil {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}
}
