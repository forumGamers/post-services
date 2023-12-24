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
	b.BaseRepo
}

func NewShareRepo() ShareRepo {
	return &ShareRepoImpl{b.NewBaseRepo(b.GetCollection(b.Share))}
}

func (r *ShareRepoImpl) DeleteMany(ctx context.Context, postId primitive.ObjectID) error {
	return r.DeleteManyByQuery(ctx, bson.M{"postId": postId})
}
