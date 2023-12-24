package reply

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/post-services/errors"
	h "github.com/post-services/helper"
	m "github.com/post-services/models"
	"github.com/post-services/pkg/comment"
	"github.com/post-services/web"
)

type ReplyService interface {
	ValidateReply(data *web.CommentForm) error
	CreatePayload(data web.CommentForm, userId string) m.ReplyComment
	AuthorizeDeleteReply(data m.ReplyComment, user m.User) error
}

type ReplyServiceImpl struct {
	Repo     comment.CommentRepo
	Validate *validator.Validate
}

func NewReplyService(repo comment.CommentRepo, validate *validator.Validate) ReplyService {
	return &ReplyServiceImpl{
		repo,
		validate,
	}
}

func (rs *ReplyServiceImpl) ValidateReply(data *web.CommentForm) error {
	return rs.Validate.Struct(data)
}

func (rs *ReplyServiceImpl) CreatePayload(data web.CommentForm, userId string) m.ReplyComment {
	return m.ReplyComment{
		UserId:    userId,
		Text:      h.Encryption(data.Text),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (rs *ReplyServiceImpl) AuthorizeDeleteReply(data m.ReplyComment, user m.User) error {
	//nanti yang punya post juga bisa hapus
	if user.UUID != data.UserId || user.LoggedAs != "Admin" {
		return errors.NewError("unauthorized", 401)
	}
	return nil
}
