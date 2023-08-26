package services

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	h "github.com/post-services/helper"
	m "github.com/post-services/models"
	r "github.com/post-services/repository"
	tp "github.com/post-services/third-party"
	"github.com/post-services/web"
)

type PostService interface {
	ValidatePostInput(data *web.PostForm) error
	CreatePostPayload(ctx context.Context,data *web.PostForm,user m.User,file tp.UploadFile)(m.Post,error)
}

type PostServiceImpl struct {
	Repo 		r.PostRepo
	Validate 	*validator.Validate
	// ImageKit 	tp.ImageKit
}

func NewPostService(repo r.PostRepo,validate *validator.Validate) PostService {
	return &PostServiceImpl{
		Repo: repo,
		Validate: validate,
		// ImageKit: ik,
	}
}

func (ps *PostServiceImpl) ValidatePostInput(data *web.PostForm) error {
	return ps.Validate.Struct(data)
}

func (ps *PostServiceImpl) CreatePostPayload(
	ctx context.Context,
	data *web.PostForm,
	user m.User,
	file tp.UploadFile,
	)(m.Post,error){
	// imageKitCh := make(chan tp.ImageKitResult)

	// if file.Data != nil {
	// 	go ps.ImageKit.Upload(ctx,file.Data,file.Name,file.Folder,imageKitCh)
	// }else {
	// 	go func () {
	// 		imageKitCh <- tp.ImageKitResult{ Url: "" ,FileId: "" ,Error: nil }
	// 	}()
	// }

	// result := ps.ImageKit.Upload(ctx,file.Data,file.Name,file.Folder)

	// if result.Error != nil {
	// 	return m.Post{},h.BadGateway
	// }

	return m.Post{
		UserId: user.Id,
		Text: h.Encryption(data.Text),
		Media: m.Media{
			// Url: result.Url,
			// Id: result.FileId,
			// Type: file.Folder,
			Url: "",
			Id: "",
			Type: "",
		},
		AllowComment: data.AllowComment,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Privacy: data.Privacy,
	},nil
}