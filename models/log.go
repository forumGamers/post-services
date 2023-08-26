package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Log struct {
	Id 				primitive.ObjectID	`json:"id" bson:"id,omitempty"`
	Path         	string 				`json:"path" bson:"path,omitempty"`
	UserId       	int    				`json:"UserId" bson:"userId,omitempty"`
	Method       	string 				`json:"method" bson:"method,omitempty"`
	StatusCode   	int    				`json:"statusCode" bson:"statusCode,omitempty"`
	ResponseTime 	int    				`json:"responseTime" bson:"responseTime,omitempty"`
	Origin       	string 				`json:"origin"`
	CreatedAt		time.Time			`json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt		time.Time			`json:"updatedAt" bson:"updatedAt,omitempty"`
}
