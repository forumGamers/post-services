package like

import (
	"context"

	h "github.com/post-services/helper"
	m "github.com/post-services/models"
	b "github.com/post-services/pkg/base"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type LikeRepo interface {
	DeletePostLikes(ctx context.Context, postId primitive.ObjectID) error
	GetLikesByUserIdAndPostId(ctx context.Context, postId primitive.ObjectID, userId string, result *m.Like) error
	AddLikes(ctx context.Context, like *m.Like) (primitive.ObjectID, error)
	DeleteLike(ctx context.Context, postId primitive.ObjectID, userId string) error
}

type LikeRepoImpl struct {
	b.BaseRepoImpl
}

func NewLikeRepo() LikeRepo {
	return &LikeRepoImpl{
		BaseRepoImpl: *b.NewBaseRepo(b.GetCollection(b.Like)),
	}
}

func (r *LikeRepoImpl) DeletePostLikes(ctx context.Context, postId primitive.ObjectID) error {
	return r.DeleteMany(ctx, bson.M{"postId": postId})
}

func (r *LikeRepoImpl) GetLikesByUserIdAndPostId(ctx context.Context, postId primitive.ObjectID, userId string, result *m.Like) error {
	if err := r.DB.FindOne(ctx, bson.M{
		"userId": userId,
		"postId": postId,
	}).Decode(&result); err != nil {
		if err == mongo.ErrNoDocuments {
			return h.NotFount
		}
		return err
	}
	return nil
}

func (r *LikeRepoImpl) AddLikes(ctx context.Context, like *m.Like) (primitive.ObjectID, error) {
	result, err := r.DB.InsertOne(ctx, like)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

func (r *LikeRepoImpl) DeleteLike(ctx context.Context, postId primitive.ObjectID, userId string) error {
	if _, err := r.DB.DeleteOne(ctx, bson.M{
		"postId": postId,
		"userId": userId,
	}); err != nil {
		if err == mongo.ErrNoDocuments {
			return h.NotFount
		}
		return err
	}
	return nil
}
