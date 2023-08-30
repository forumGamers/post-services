package like

import (
	"context"

	b "github.com/post-services/pkg/base"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LikeRepo interface {
	DeletePostLikes(ctx context.Context,postId primitive.ObjectID) error
}

type LikeRepoImpl struct {
	b.BaseRepoImpl
}

func NewLikeRepo() LikeRepo {
	return &LikeRepoImpl{
		BaseRepoImpl: *b.NewBaseRepo(b.GetCollection(b.Like)),
	}
}

func (r *LikeRepoImpl) DeletePostLikes(ctx context.Context,postId primitive.ObjectID) error {
	return r.DeleteMany(ctx,bson.M{ "postId": postId })
}
