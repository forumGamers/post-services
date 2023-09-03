package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Media struct {
	Url  string `json:"url" bson:"url,omitempty"`
	Type string `json:"type" bson:"type,omitempty"`
	Id   string `json:"id" bson:"id,omitempty"`
}

type Post struct {
	Id           primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	UserId       string             `json:"userId" bson:"userId,omitempty"`
	Text         string             `json:"text" bson:"text"`
	Media        Media
	AllowComment bool `json:"allowComment" bson:"allowComment" default:"true"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Tags         []string `json:"tags" bson:"tags,omitempty"`
	Privacy      string   `json:"privacy" bson:"privacy" default:"Public"`
}
