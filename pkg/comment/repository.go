package comment

import (
	"context"

	m "github.com/post-services/models"
	b "github.com/post-services/pkg/base"
)

type CommentRepo interface {
	CreateComment(ctx context.Context, data *m.Comment) error
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
