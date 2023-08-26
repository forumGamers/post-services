package config

import (
	"context"
	"fmt"
	"os"

	h "github.com/post-services/helper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connection() *mongo.Database {

	uri := os.Getenv("DATABASE_URL")

	if uri == "" {
		uri = "mongodb://localhost:27017"
	}

	client,err := mongo.Connect(context.Background(),options.Client().ApplyURI(uri))
	h.PanicIfError(err)

	fmt.Println("connection success")

	return client.Database("Post")
}