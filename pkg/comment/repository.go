package comment

import (
	"context"

	m "github.com/post-services/models"
	b "github.com/post-services/pkg/base"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CommentRepo interface {
	CreateComment(ctx context.Context, data *m.Comment) error
	CreateReply(ctx context.Context, id primitive.ObjectID, data *m.ReplyComment) error
	FindById(ctx context.Context, id primitive.ObjectID, data *m.Comment) error
	DeleteOne(ctx context.Context, id primitive.ObjectID) error
	CreateMany(ctx context.Context, datas []any) (*mongo.InsertManyResult, error)
	DeleteReplyByPostId(ctx context.Context, postId primitive.ObjectID) error
	FindReplyById(ctx context.Context, id, replyId primitive.ObjectID, data *m.ReplyComment) error
	DeleteOneReply(ctx context.Context, id, replyId primitive.ObjectID) error
	DeleteMany(ctx context.Context, postId primitive.ObjectID) error
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

func (r *CommentRepoImpl) CreateMany(ctx context.Context, datas []any) (*mongo.InsertManyResult, error) {
	return r.DB.InsertMany(ctx, datas)
}

func (r *CommentRepoImpl) CreateReply(ctx context.Context, id primitive.ObjectID, data *m.ReplyComment) error {
	if _, err := r.BaseRepoImpl.DB.UpdateByID(ctx, bson.M{
		"_id": id,
	}, bson.M{
		"$push": bson.M{
			"reply": data,
		},
	}); err != nil {
		return err
	}
	return nil
}

func (r *CommentRepoImpl) DeleteReplyByPostId(ctx context.Context, postId primitive.ObjectID) error {
	cursor, err := r.BaseRepoImpl.DB.Find(ctx, bson.M{"postId": postId})
	if err != nil {
		return err
	}

	defer cursor.Close(ctx)
	var commentIds []primitive.ObjectID
	for cursor.Next(ctx) {
		var comment struct {
			CommentId primitive.ObjectID `bson:"_id"`
		}
		if err := cursor.Decode(&comment); err != nil {
			return err
		}
		commentIds = append(commentIds, comment.CommentId)
	}

	if len(commentIds) > 0 {
		if _, err := r.DB.DeleteMany(ctx, bson.M{
			"commentId": bson.M{
				"$in": commentIds,
			},
		}); err != nil {
			return err
		}
	}
	return nil
}

func (r *CommentRepoImpl) FindReplyById(ctx context.Context, id, replyId primitive.ObjectID, data *m.ReplyComment) error {
	if err := r.BaseRepoImpl.DB.FindOne(ctx, bson.M{
		"_id": id,
		"reply": bson.M{
			"$elemMatch": bson.M{
				"_id": replyId,
			},
		},
	}).Decode(&data); err != nil {
		return err
	}
	return nil
}

func (r *CommentRepoImpl) DeleteOneReply(ctx context.Context, id, replyId primitive.ObjectID) error {
	if _, err := r.BaseRepoImpl.DB.DeleteOne(ctx, bson.M{
		"_id": id,
		"reply": bson.M{
			"$pull": bson.M{
				"$elemMatch": bson.M{
					"_id": replyId,
				},
			},
		},
	}); err != nil {
		return err
	}
	return nil
}

func (r *CommentRepoImpl) DeleteMany(ctx context.Context, postId primitive.ObjectID) error {
	return r.BaseRepoImpl.DeleteMany(ctx, bson.M{
		"postId": postId,
	})
}
