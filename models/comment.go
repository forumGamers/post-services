package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	Id        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	UserId    string             `json:"userId" bson:"userId,omitempty"`
	Text      string             `json:"text" bson:"text,omitempty"`
	PostId    primitive.ObjectID `json:"postId" bson:"postId,omitempty"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
