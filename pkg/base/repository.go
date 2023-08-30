package base

import (
	"context"

	cfg "github.com/post-services/config"
	h "github.com/post-services/helper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CollectionName string

const (
	Post 	CollectionName = "post"
	Like    CollectionName = "like"
	Comment CollectionName = "comment"
	Reply   CollectionName = "replyComment"
	Share   CollectionName = "share"
)

type BaseRepo interface {
	DeleteMany(ctx context.Context,filter any) error
	DeleteOneById(ctx context.Context,id primitive.ObjectID) error
	FindOneById(ctx context.Context,id primitive.ObjectID,data any) error
}

type BaseRepoImpl struct {
	DB *mongo.Collection
}

func NewBaseRepo(db *mongo.Collection) *BaseRepoImpl {
	return &BaseRepoImpl{
		DB: db,
	}
}

func (r *BaseRepoImpl) DeleteMany(ctx context.Context,filter any) error {
	if _,err := r.DB.DeleteMany(ctx,filter) ; err != nil {
		if err == mongo.ErrNoDocuments {
			return h.NotFount
		}
		return err
	}
	return nil
}

func (r *BaseRepoImpl) DeleteOneById(ctx context.Context,id primitive.ObjectID) error {
	if _,err := r.DB.DeleteOne(ctx,bson.M{
		"_id":id,
	}) ; err != nil {
		if err == mongo.ErrNoDocuments {
			return h.NotFount
		}
		return err
	}
	return nil
}

func(r *BaseRepoImpl) FindOneById(ctx context.Context,id primitive.ObjectID,data any) error {
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

func GetCollection(name CollectionName) *mongo.Collection {
	return cfg.DB.Collection(string(name))
}