package model

import (
	"time"
)

type Like struct {
	UserId         string    `json:"userId" bson:"userId"`
	BlogId         string    `json:"blogId" bson:"blogId"`
	DateOfCreation time.Time `json:"dateOfCreation" bson:"date_of_creation"`
}

func CreateNewLike(userId string, blogId string) *Like {
	return &Like{
		UserId:         userId,
		BlogId:         blogId,
		DateOfCreation: time.Now(),
	}
}
