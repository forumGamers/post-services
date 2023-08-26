package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserPreference struct {
	Id 			primitive.ObjectID 	`json:"_id" bson:"_id,omitempty"`
	Tag			string				`json:"tag" bson:"tag,omitempty"`
	CreatedAt	time.Time
	UpdatedAt	time.Time
}