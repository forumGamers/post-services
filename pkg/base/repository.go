package base

import (
	"context"

	cfg "github.com/post-services/config"
	"github.com/post-services/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CollectionName string

const (
	Post    CollectionName = "post"
	Like    CollectionName = "like"
	Comment CollectionName = "comment"
	Reply   CollectionName = "replyComment"
	Share   CollectionName = "share"
	Log     CollectionName = "log"
)

type BaseRepo interface {
	DeleteMany(ctx context.Context, filter any) error
	DeleteOneById(ctx context.Context, id primitive.ObjectID) error
	FindOneById(ctx context.Context, id primitive.ObjectID, data any) error
	InsertMany(ctx context.Context, data []any) (*mongo.InsertManyResult, error)
}

type BaseRepoImpl struct {
	DB *mongo.Collection
}

func NewBaseRepo(db *mongo.Collection) *BaseRepoImpl {
	return &BaseRepoImpl{
		DB: db,
	}
}

func (r *BaseRepoImpl) DeleteMany(ctx context.Context, filter any) error {
	if _, err := r.DB.DeleteMany(ctx, filter); err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.NewError("Data not found", 404)
		}
		return err
	}
	return nil
}

func (r *BaseRepoImpl) DeleteOneById(ctx context.Context, id primitive.ObjectID) error {
	if _, err := r.DB.DeleteOne(ctx, bson.M{
		"_id": id,
	}); err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.NewError("Data not found", 404)
		}
		return err
	}
	return nil
}

func (r *BaseRepoImpl) FindOneById(ctx context.Context, id primitive.ObjectID, data any) error {
	if err := r.DB.FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(data); err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.NewError("Data not found", 404)
		}
		return err
	}
	return nil
}

func (r *BaseRepoImpl) InsertMany(ctx context.Context, data []any) (*mongo.InsertManyResult, error) {
	return r.DB.InsertMany(ctx, data)
}

func (r *BaseRepoImpl) Create(ctx context.Context, data any) (primitive.ObjectID, error) {
	result, err := r.DB.InsertOne(ctx, data)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

func GetCollection(name CollectionName) *mongo.Collection {
	return cfg.DB.Collection(string(name))
}
