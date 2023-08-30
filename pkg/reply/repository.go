package reply

import (
	"context"

	b "github.com/post-services/pkg/base"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReplyRepo interface {
	DeleteReplyByPostId(ctx context.Context,postId primitive.ObjectID) error
}

type ReplyRepoImpl struct {
	b.BaseRepoImpl
}

func NewReplyRepo() ReplyRepo {
	return &ReplyRepoImpl{
		BaseRepoImpl: *b.NewBaseRepo(b.GetCollection(b.Reply)),
	}
}

func (r *ReplyRepoImpl) DeleteReplyByPostId(ctx context.Context,postId primitive.ObjectID) error {
	commentCollection := b.GetCollection(b.Comment)
	cursor,err := commentCollection.Find(ctx,bson.M{"postId":postId}) 
	if err != nil {
		return err
	}

	defer cursor.Close(ctx)
	var commentIds []primitive.ObjectID
	for cursor.Next(ctx) {
		var comment struct {
			CommentId 	primitive.ObjectID	`bson:"_id"`
		}
		if err := cursor.Decode(&comment) ; err != nil {
			return err
		}
		commentIds = append(commentIds, comment.CommentId)
	}

	if _,err := r.DB.DeleteMany(ctx,bson.M{
		"commentId": bson.M{
			"$in":commentIds,
		},
	}) ; err != nil {
		return err
	}

	return nil
}