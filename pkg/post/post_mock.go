package post

import (
	"context"

	m "github.com/post-services/models"
	tp "github.com/post-services/third-party"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostRepoMockImpl struct {
	Mock mock.Mock
}

type ImagekitMockImpl struct {
	Mock mock.Mock
}

func (repo *PostRepoMockImpl) Create(ctx context.Context, data *m.Post) error { return nil }
func (repo *PostRepoMockImpl) FindById(ctx context.Context, id primitive.ObjectID, data *m.Post) error {
	return nil
}
func (repo *PostRepoMockImpl) DeletePost(ctx context.Context, id primitive.ObjectID) error {
	return nil
}
func (repo *PostRepoMockImpl) GetSession() (mongo.Session, error)                         { return nil, nil }
func (repo *PostRepoMockImpl) DeleteOne(ctx context.Context, id primitive.ObjectID) error { return nil }

func (ik *ImagekitMockImpl) UploadFile(ctx context.Context, file []byte, fileName string, folder string) tp.ImageKitResult {
	return tp.ImageKitResult{}
}
func (ik *ImagekitMockImpl) UpdateImage(ctx context.Context, file []byte, fileName string, folder string, updatedFileID string, resultCh chan<- tp.ImageKitResult) {
}
func (ik *ImagekitMockImpl) Delete(ctx context.Context, imageId string, ch chan<- error) {}
