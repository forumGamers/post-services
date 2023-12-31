package like

import (
	"context"

	"github.com/post-services/errors"
	m "github.com/post-services/models"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type LikeRepoMockImpl struct {
	Mock mock.Mock
}

func (r *LikeRepoMockImpl) DeletePostLikes(ctx context.Context, postId primitive.ObjectID) error {
	return nil
}

func (r *LikeRepoMockImpl) GetLikesByUserIdAndPostId(ctx context.Context, postId primitive.ObjectID, userId string, result *m.Like) error {
	args := r.Mock.Called(ctx, postId, userId, result)
	switch args.Get(0) {
	case nil:
		return nil
	case errors.NewError("Data not found", 404):
		return errors.NewError("Data not found", 404)
	default:
		return nil
	}
}

func (r *LikeRepoMockImpl) AddLikes(ctx context.Context, like *m.Like) (primitive.ObjectID, error) {
	args := r.Mock.Called(ctx, like)
	if args.Get(1) != nil {
		return primitive.NilObjectID, args.Error(1)
	}
	return primitive.NewObjectID(), nil
}

func (r *LikeRepoMockImpl) DeleteLike(ctx context.Context, postId primitive.ObjectID, userId string) error {
	return nil
}

func (r *LikeRepoMockImpl) CreateMany(ctx context.Context, datas []any) (*mongo.InsertManyResult, error) {
	return nil, nil
}
