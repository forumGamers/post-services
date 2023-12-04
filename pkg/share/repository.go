package share

import (
	"context"

	b "github.com/post-services/pkg/base"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ShareRepo interface {
	DeleteMany(ctx context.Context, postId primitive.ObjectID) error
}

type ShareRepoImpl struct {
	b.BaseRepoImpl
}

func NewShareRepo() ShareRepo {
	return &ShareRepoImpl{
		BaseRepoImpl: *b.NewBaseRepo(b.GetCollection(b.Share)),
	}
}

func (r *ShareRepoImpl) DeleteMany(ctx context.Context, postId primitive.ObjectID) error {
	return r.BaseRepoImpl.DeleteMany(ctx, bson.M{
		"postId": postId,
	})
}
