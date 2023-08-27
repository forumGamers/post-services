package services_test

import (
	"context"
	"testing"

	"github.com/joho/godotenv"
	m "github.com/post-services/models"
	s "github.com/post-services/services"
	tp "github.com/post-services/third-party"
	v "github.com/post-services/validations"
	"github.com/post-services/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostRepoMock struct {
	Mock 	mock.Mock
}

type ImagekitMock struct {
	Mock 	mock.Mock
}

func(repo *PostRepoMock) Create(ctx context.Context,data *m.Post){}
func(repo *PostRepoMock) FindById(ctx context.Context, id primitive.ObjectID,data *m.Post) error{return nil}

func(ik *ImagekitMock) UploadFile(ctx context.Context,file []byte,fileName string,folder string) tp.ImageKitResult{return tp.ImageKitResult{}}
func(ik *ImagekitMock)UpdateImage(ctx context.Context,file []byte,fileName string,folder string,updatedFileID string,resultCh chan<- tp.ImageKitResult){}
func(ik *ImagekitMock) Delete(ctx context.Context,imageId string,ch chan<- error){}

var PostRepo = &PostRepoMock{Mock: mock.Mock{}}
var ImageKit = &ImagekitMock{Mock: mock.Mock{}}

var PostServiceTest = &s.PostServiceImpl{
	Repo: PostRepo,
	Validate: v.GetValidator(),
	ImageKit: ImageKit,
} 

func TestMain(m *testing.M) {
	godotenv.Load()

	m.Run()
}

func TestValidatePostInput(t *testing.T) {
	datas := []struct{
		Name		string
		Data		web.PostForm
		Error	bool
	}{
		{
			Name: "File = nil & Text = ''",
			Data: web.PostForm{
				File: nil,
				Text: "",
				Privacy: "Public",
				AllowComment: true,
			},
			Error: true,
		},
	}

	for _,data := range datas {
		t.Run(data.Name,func(t *testing.T) {
			err := PostServiceTest.ValidatePostInput(&data.Data)

			if data.Error {
				assert.NotNil(t,err)
			}else {
				assert.Nil(t,err)
			}
		})
	}
}