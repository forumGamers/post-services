package like

import (
	"github.com/go-playground/validator/v10"
)

type LikeService interface{}

type LikeServiceImpl struct {
	Repo     LikeRepo
	Validate *validator.Validate
}

func NewLikeService(repo LikeRepo,validate *validator.Validate) LikeService {
	return &LikeServiceImpl{
		Repo: repo,
		Validate: validate,
	}
}

