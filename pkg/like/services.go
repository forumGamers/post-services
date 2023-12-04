package like

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/post-services/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LikeService interface {
	InsertManyAndBindIds(ctx context.Context, likes []models.Like) error
}

type LikeServiceImpl struct {
	Repo     LikeRepo
	Validate *validator.Validate
}

func NewLikeService(repo LikeRepo, validate *validator.Validate) LikeService {
	return &LikeServiceImpl{
		Repo:     repo,
		Validate: validate,
	}
}

func (ls *LikeServiceImpl) InsertManyAndBindIds(ctx context.Context, likes []models.Like) error {
	var payload []any

	for _, data := range likes {
		data.Id = primitive.NilObjectID
		payload = append(payload, data)
	}

	ids, err := ls.Repo.CreateMany(ctx, payload)
	if err != nil {
		return err
	}

	for i := 0; i < len(ids.InsertedIDs); i++ {
		id := ids.InsertedIDs[i].(primitive.ObjectID)
		likes[i].Id = id
	}
	return nil
}
