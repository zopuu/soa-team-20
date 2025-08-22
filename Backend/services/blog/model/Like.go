package model

import (
	"time"
)

type Like struct {
	UserId         string    `json:"userId" bson:"user_id"`
	BlogId         string    `json:"blogId" bson:"blog_id"`
	DateOfCreation time.Time `json:"dateOfCreation" bson:"date_of_creation"`
}
