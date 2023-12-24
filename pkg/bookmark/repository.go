package bookmark

import (
	"context"

	b "github.com/post-services/pkg/base"
)

func NewBookMarkRepo() BookmarkRepo {
	return &BookmarkRepoImpl{b.NewBaseRepo(b.GetCollection(b.Bookmark))}
}

func (r *BookmarkRepoImpl) CreateOne(ctx context.Context, data *Bookmark) error {
	if result, err := r.Create(ctx, data); err != nil {
		return err
	} else {
		data.Id = result
	}
	return nil
}

func (r *BookmarkRepoImpl) FindOne(ctx context.Context, query any, result *Bookmark) error {
	return r.FindOneByQuery(ctx, query, result)
}
