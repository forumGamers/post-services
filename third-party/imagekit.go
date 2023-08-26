package thirdparty

import (
	"context"
	"os"

	"github.com/codedius/imagekit-go"
	h "github.com/post-services/helper"
)

type ImageKitService interface {
	Upload(ctx context.Context,file []byte,fileName string,folder string)ImageKitResult
	UpdateImage(ctx context.Context,file []byte,fileName string,folder string,updatedFileID string,resultCh chan<- ImageKitResult)
	Delete(ctx context.Context,imageId string,ch chan<- error)
}

type ImageKit struct {
	Client *imagekit.Client
}

type ImageKitResult struct {
	Url    string
	FileId string
	Error  error
}

type UploadFile struct {
	Data []byte
	Name string
	Folder string 
}

func ImageKitConnection() ImageKitService {
	ik,err := imagekit.NewClient(&imagekit.Options{
		PublicKey: os.Getenv("IMAGEKIT_PRIVATE_KEY"),
		PrivateKey: os.Getenv("IMAGEKIT_PUBLIC_KEY"),
	})
	h.PanicIfError(err)

	return &ImageKit{
		Client: ik,
	}
}

func (ik ImageKit) Upload(ctx context.Context,file []byte,fileName string,folder string) ImageKitResult {
	uploadResponse, err := ik.Client.Upload.ServerUpload(ctx,&imagekit.UploadRequest{
		File:              file,
		FileName:          fileName,
		UseUniqueFileName: true,
		Folder:            folder,
	})

	return ImageKitResult{
		Url: uploadResponse.URL,
		FileId: uploadResponse.FileID,
		Error: err,
	}
}

func (ik *ImageKit) UpdateImage(ctx context.Context,file []byte,fileName string,folder string,updatedFileID string,ch chan<- ImageKitResult){
	go func(){
		ch <- ik.Upload(ctx,file,fileName,folder)
	}()

	go func(){
		ik.Client.Media.DeleteFile(ctx,updatedFileID)
	}()
}

func (ik *ImageKit) Delete(ctx context.Context,imageId string,ch chan<- error) {
	ch <- ik.Client.Media.DeleteFile(ctx,imageId) 
}
