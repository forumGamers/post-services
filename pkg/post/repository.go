package post

import (
	"context"

	h "github.com/post-services/helper"
	m "github.com/post-services/models"
	b "github.com/post-services/pkg/base"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostRepo interface {
	Create(ctx context.Context,data *m.Post)
	FindById(ctx context.Context,id primitive.ObjectID,data *m.Post) error
	GetSession() (mongo.Session,error)
	DeleteOne(ctx context.Context,id primitive.ObjectID) error
}

type PostRepoImpl struct {
	b.BaseRepoImpl
}

func NewPostRepo() PostRepo {
	return &PostRepoImpl{
		BaseRepoImpl: *b.NewBaseRepo(b.GetCollection(b.Post)),
	}
}

func (r *PostRepoImpl) Create(ctx context.Context,data *m.Post) {
	result,err := r.DB.InsertOne(ctx,data) 
	h.PanicIfError(err)

	data.Id = result.InsertedID.(primitive.ObjectID)
}

func (r *PostRepoImpl) FindById(ctx context.Context,id primitive.ObjectID,data *m.Post) error {
	return r.FindOneById(ctx,id,data) 
}

func (r *PostRepoImpl) GetSession() (mongo.Session,error) {
	return r.DB.Database().Client().StartSession()
}

func (r *PostRepoImpl) DeleteOne(ctx context.Context,id primitive.ObjectID) error {
	return r.DeleteOneById(ctx,id)
}