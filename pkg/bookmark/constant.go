package bookmark

import (
	"context"

	b "github.com/post-services/pkg/base"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookmarkRepo interface {
	CreateOne(ctx context.Context, data *Bookmark) error
	FindOne(ctx context.Context, query any, result *Bookmark) error
}

type BookmarkRepoImpl struct {
	b.BaseRepo
}

type BookmarkService interface {
	CreatePayload(postId primitive.ObjectID, userId string) Bookmark
}

type BookmarkServiceImpl struct {
	Repo BookmarkRepo
}
