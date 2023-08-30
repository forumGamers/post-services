package comment

import (
	b "github.com/post-services/pkg/base"
	"go.mongodb.org/mongo-driver/mongo"
)

type CommentRepo interface {}

type CommentRepoImpl struct {
	b.BaseRepoImpl
}

func NewCommentRepo(db *mongo.Collection) CommentRepo {
	return &CommentRepoImpl{
		BaseRepoImpl: *b.NewBaseRepo(b.GetCollection(b.Comment)),
	}
}