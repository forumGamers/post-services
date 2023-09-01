package comment

import (
	"time"

	"github.com/go-playground/validator/v10"
	h "github.com/post-services/helper"
	m "github.com/post-services/models"
	"github.com/post-services/web"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommentService interface {
	ValidateComment(data *web.CommentForm) error
	CreatePayload(data web.CommentForm, postId primitive.ObjectID, userId int) m.Comment
	AuthorizeDeleteComment(data m.Comment, user m.User) error
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

func (ps *CommentServiceImpl) CreatePayload(data web.CommentForm, postId primitive.ObjectID, userId int) m.Comment {
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
	if user.Id != data.UserId || user.Role != "Admin" {
		return h.AccessDenied
	}
	return nil
}
