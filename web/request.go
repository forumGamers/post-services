package web

import (
	"mime/multipart"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostForm struct {
	File         *multipart.FileHeader `form:"file" validate:"required_without=Text"`
	Text         string                `form:"text" validate:"required_without=File"`
	AllowComment bool                  `form:"allowComment"`
	Privacy      string                `form:"privacy" validate:"oneof=Public FriendOnly Private"`
}

type CommentForm struct {
	Text string `json:"text" validate:"required" form:"text"`
}

type PostData struct {
	Id           primitive.ObjectID `json:"_id"`
	UserId       string             `json:"userId"`
	Text         string             `json:"text"`
	AllowComment bool               `json:"allowComment"`
	Privacy      string             `json:"privacy"`
	Media        struct {
		Type string `json:"type"`
		Url  string `json:"url"`
		Id   string `json:"id"`
	} `json:"media"`
}

type PostDatas struct {
	Datas []PostData `json:"datas" binding:"required"`
}

type LikeData struct {
	Id     primitive.ObjectID `json:"_id"`
	UserId string             `json:"userId"`
	PostId primitive.ObjectID `json:"postId"`
}

type LikeDatas struct {
	Datas []LikeData `json:"datas" binding:"required"`
}
