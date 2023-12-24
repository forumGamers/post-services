package like_test

import (
	"context"
	"testing"

	"github.com/post-services/errors"
	m "github.com/post-services/models"
	l "github.com/post-services/pkg/like"
	v "github.com/post-services/validations"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var mockRepo = &l.LikeRepoMockImpl{Mock: mock.Mock{}}

var mockService = l.LikeServiceImpl{
	Repo:     mockRepo,
	Validate: v.GetValidator(),
}

func TestMain(m *testing.M) {
	m.Run()
}

func TestGetLikesByUserIdAndPostId_NotFound(t *testing.T) {
	defer func() {
		mockRepo.Mock.ExpectedCalls = nil
	}()
	ctx := context.Background()
	postId, _ := primitive.ObjectIDFromHex("64e2ff258c78c4a3ff840e9d")
	userId := "1"
	var data m.Like
	mockRepo.Mock.On("GetLikesByUserIdAndPostId", ctx, postId, userId, &data).Return(errors.NewError("Data not found", 404))

	err := mockService.Repo.GetLikesByUserIdAndPostId(ctx, postId, userId, &data)

	assert.NotNil(t, err)
	assert.Equal(t, errors.NewError("Data not found", 404), err, "Error must be not found error")
}

func TestGetLikesByUserIdAndPostId_Success(t *testing.T) {
	defer func() {
		mockRepo.Mock.ExpectedCalls = nil
	}()
	ctx := context.Background()
	postId, _ := primitive.ObjectIDFromHex("64e2ff258c78c4a3ff840e9d")
	userId := "1"
	var data m.Like
	mockRepo.Mock.On("GetLikesByUserIdAndPostId", ctx, postId, userId, &data).Return(nil)

	assert.Nil(t, mockService.Repo.GetLikesByUserIdAndPostId(ctx, postId, userId, &data))
}
