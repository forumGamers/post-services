package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Share struct {
	Id        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	UserId    string             `json:"userId" bson:"userId,omitempty"`
	PostId    primitive.ObjectID `json:"postId" bson:"postId,omitempty"`
	Text      string             `json:"text" bson:"text"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
