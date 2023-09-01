package comment

import (
	"context"

	m "github.com/post-services/models"
	b "github.com/post-services/pkg/base"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommentRepo interface {
	CreateComment(ctx context.Context, data *m.Comment) error
	FindById(ctx context.Context, id primitive.ObjectID, data *m.Comment) error
	DeleteOne(ctx context.Context, id primitive.ObjectID) error
}

type CommentRepoImpl struct {
	b.BaseRepoImpl
}

func NewCommentRepo() CommentRepo {
	return &CommentRepoImpl{
		BaseRepoImpl: *b.NewBaseRepo(b.GetCollection(b.Comment)),
	}
}

func (r *CommentRepoImpl) CreateComment(ctx context.Context, data *m.Comment) error {
	result, err := r.BaseRepoImpl.Create(ctx, &data)
	if err != nil {
		return err
	}
	data.Id = result
	return nil
}

func (r *CommentRepoImpl) FindById(ctx context.Context, id primitive.ObjectID, data *m.Comment) error {
	return r.FindOneById(ctx, id, data)
}

func (r *CommentRepoImpl) DeleteOne(ctx context.Context, id primitive.ObjectID) error {
	return r.DeleteOneById(ctx, id)
}
