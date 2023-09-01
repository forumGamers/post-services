package web

import "mime/multipart"

type PostForm struct {
	File         *multipart.FileHeader `form:"file" validate:"required_without=Text"`
	Text         string                `form:"text" validate:"required_without=File"`
	AllowComment bool                  `form:"allowComment"`
	Privacy      string                `form:"privacy" validate:"oneof=Public FriendOnly Private"`
}

type CommentForm struct {
	Text string `json:"text" validate:"required"`
}
