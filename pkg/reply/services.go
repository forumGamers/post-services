package reply

import (
	"time"

	"github.com/go-playground/validator/v10"
	h "github.com/post-services/helper"
	m "github.com/post-services/models"
	"github.com/post-services/web"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReplyService interface {
	ValidateReply(data *web.CommentForm) error
	CreatePayload(data web.CommentForm, commentId primitive.ObjectID, userId string) m.ReplyComment
	AuthorizeDeleteReply(data m.ReplyComment, user m.User) error
}

type ReplyServiceImpl struct {
	Repo     ReplyRepo
	Validate *validator.Validate
}

func NewReplyService(repo ReplyRepo, validate *validator.Validate) ReplyService {
	return &ReplyServiceImpl{
		Repo:     repo,
		Validate: validate,
	}
}

func (rs *ReplyServiceImpl) ValidateReply(data *web.CommentForm) error {
	return rs.Validate.Struct(data)
}

func (rs *ReplyServiceImpl) CreatePayload(data web.CommentForm, commentId primitive.ObjectID, userId string) m.ReplyComment {
	return m.ReplyComment{
		UserId:    userId,
		Text:      h.Encryption(data.Text),
		CommentId: commentId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (rs *ReplyServiceImpl) AuthorizeDeleteReply(data m.ReplyComment, user m.User) error {
	//nanti yang punya post juga bisa hapus
	if user.UUID != data.UserId || user.Role != "Admin" {
		return h.AccessDenied
	}
	return nil
}
