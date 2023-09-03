package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReplyComment struct {
	Id        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	UserId    string             `json:"userId" bson:"userId,omitempty"`
	Text      string             `json:"text" bson:"text,omitempty"`
	CommentId primitive.ObjectID `json:"commentId" bson:"commentId,omitempty"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
