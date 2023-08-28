package post

import (
	"context"

	h "github.com/post-services/helper"
	m "github.com/post-services/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostRepo interface {
	Create(ctx context.Context,data *m.Post)
	FindById(ctx context.Context,id primitive.ObjectID,data *m.Post) error
}

type PostRepoImpl struct {
	DB 		*mongo.Collection
}

func NewPostRepo(db *mongo.Collection) PostRepo {
	return &PostRepoImpl{
		DB: db,
	}
}

func (r *PostRepoImpl) Create(ctx context.Context,data *m.Post) {
	result,err := r.DB.InsertOne(ctx,data) 
	h.PanicIfError(err)

	data.Id = result.InsertedID.(primitive.ObjectID)
}

func (r *PostRepoImpl) FindById(ctx context.Context,id primitive.ObjectID,data *m.Post) error {
	if err := r.DB.FindOne(ctx,bson.M{
		"_id":id,
	}).Decode(data) ; err != nil {
		if err == mongo.ErrNoDocuments {
			return h.NotFount
		}
		return err
	}
	return nil
}
