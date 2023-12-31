package like

import (
	"context"

	"github.com/post-services/errors"
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
	CreateMany(ctx context.Context, datas []any) (*mongo.InsertManyResult, error)
}

type LikeRepoImpl struct {
	b.BaseRepo
}

func NewLikeRepo() LikeRepo {
	return &LikeRepoImpl{b.NewBaseRepo(b.GetCollection(b.Like))}
}

func (r *LikeRepoImpl) DeletePostLikes(ctx context.Context, postId primitive.ObjectID) error {
	return r.DeleteManyByQuery(ctx, bson.M{"postId": postId})
}

func (r *LikeRepoImpl) GetLikesByUserIdAndPostId(ctx context.Context, postId primitive.ObjectID, userId string, result *m.Like) error {
	if err := r.FindOneByQuery(ctx, bson.M{
		"userId": userId,
		"postId": postId,
	}, &result); err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.NewError("Data not found", 404)
		}
		return err
	}
	return nil
}

func (r *LikeRepoImpl) AddLikes(ctx context.Context, like *m.Like) (primitive.ObjectID, error) {
	return r.Create(ctx, like)
}

func (r *LikeRepoImpl) DeleteLike(ctx context.Context, postId primitive.ObjectID, userId string) error {
	if err := r.DeleteOneByQuery(ctx, bson.M{
		"postId": postId,
		"userId": userId,
	}); err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.NewError("Data not found", 404)
		}
		return err
	}
	return nil
}

func (r *LikeRepoImpl) CreateMany(ctx context.Context, datas []any) (*mongo.InsertManyResult, error) {
	return r.InsertMany(ctx, datas)
}
