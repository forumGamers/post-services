package comment

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	h "github.com/post-services/helper"
	m "github.com/post-services/models"
	"github.com/post-services/web"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommentService interface {
	ValidateComment(data *web.CommentForm) error
	CreatePayload(data web.CommentForm, postId primitive.ObjectID, userId string) m.Comment
	AuthorizeDeleteComment(data m.Comment, user m.User) error
	InsertManyAndBindIds(ctx context.Context, datas []web.CommentData) error
}

type CommentServiceImpl struct {
	Repo     CommentRepo
	Validate *validator.Validate
}

func NewCommentService(repo CommentRepo, validate *validator.Validate) CommentService {
	return &CommentServiceImpl{
		Repo:     repo,
		Validate: validate,
	}
}

func (ps *CommentServiceImpl) ValidateComment(data *web.CommentForm) error {
	return ps.Validate.Struct(data)
}

func (ps *CommentServiceImpl) CreatePayload(data web.CommentForm, postId primitive.ObjectID, userId string) m.Comment {
	return m.Comment{
		UserId:    userId,
		Text:      h.Encryption(data.Text),
		PostId:    postId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (ps *CommentServiceImpl) AuthorizeDeleteComment(data m.Comment, user m.User) error {
	//nanti yang punya post juga bisa hapus
	if user.UUID != data.UserId || user.LoggedAs != "Admin" {
		return h.AccessDenied
	}
	return nil
}

func (ps *CommentServiceImpl) InsertManyAndBindIds(ctx context.Context, datas []web.CommentData) error {
	var payload []any

	for _, data := range datas {
		payload = append(payload, data)
	}

	ids, err := ps.Repo.CreateMany(ctx, payload)
	if err != nil {
		return err
	}

	for i := 0; i < len(ids.InsertedIDs); i++ {
		id := ids.InsertedIDs[i].(primitive.ObjectID)
		datas[i].Id = id
	}
	return nil
}
