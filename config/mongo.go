package config

import (
	"context"
	"fmt"
	"os"

	h "github.com/post-services/helper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database
const localDB = "mongodb://localhost:27017"

func Connection() {

	uri := os.Getenv("DATABASE_URL")

	if uri == "" {
		uri = localDB
	}

	client,err := mongo.Connect(context.Background(),options.Client().ApplyURI(uri))
	h.PanicIfError(err)

	fmt.Println("connection success")

	DB = client.Database("Post")
}

func TestingConnection() *mongo.Database {
	client,err := mongo.Connect(context.Background(),options.Client().ApplyURI(localDB))
	h.PanicIfError(err)

	return client.Database("Post_test")
}

func DisconnectConnection(client *mongo.Database) error {
	return client.Client().Disconnect(context.Background())
}