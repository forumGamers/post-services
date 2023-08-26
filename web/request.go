package web

import "mime/multipart"

type PostForm struct {
	Text         string 				`form:"text" validate:"required_without=File"`
	AllowComment bool   				`form:"allowComment"`
	File         *multipart.FileHeader	`form:"file" validate:"required_without=Text"`
	Privacy		 string					`form:"privacy" validate:"oneof=Public FriendOnly Private"`
}