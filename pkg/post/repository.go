package post

import (
	"context"

	m "github.com/post-services/models"
	b "github.com/post-services/pkg/base"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostRepo interface {
	Create(ctx context.Context, data *m.Post) error
	FindById(ctx context.Context, id primitive.ObjectID, data *m.Post) error
	GetSession() (mongo.Session, error)
	DeleteOne(ctx context.Context, id primitive.ObjectID) error
	CreateMany(ctx context.Context, datas []any) (*mongo.InsertManyResult, error)
}

type PostRepoImpl struct {
	b.BaseRepo
}

func NewPostRepo() PostRepo {
	return &PostRepoImpl{b.NewBaseRepo(b.GetCollection(b.Post))}
}

func (r *PostRepoImpl) Create(ctx context.Context, data *m.Post) error {
	result, err := r.BaseRepo.Create(ctx, data)
	if err != nil {
		return err
	}
	data.Id = result
	return nil
}

func (r *PostRepoImpl) FindById(ctx context.Context, id primitive.ObjectID, data *m.Post) error {
	return r.FindOneById(ctx, id, data)
}

func (r *PostRepoImpl) GetSession() (mongo.Session, error) {
	return r.BaseRepo.GetSession()
}

func (r *PostRepoImpl) DeleteOne(ctx context.Context, id primitive.ObjectID) error {
	return r.DeleteOneById(ctx, id)
}

func (r *PostRepoImpl) CreateMany(ctx context.Context, datas []any) (*mongo.InsertManyResult, error) {
	return r.InsertMany(ctx, datas)
}
