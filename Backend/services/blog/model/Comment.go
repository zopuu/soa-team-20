package model

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID             uuid.UUID `json:"id" bson:"_id,omitempty"`
	UserId         string    `json:"userId" bson:"user_id"`
	BlogId         string    `json:"blogId" bson:"blog_id"`
	DateOfCreation time.Time `json:"dateOfCreation" bson:"date_of_creation"`
	Text           string    `json:"text" bson:"text"`
	LastEdit       time.Time `json:"lastEdit" bson:"last_edit"`
}

func CreateNewComment(userId string, blogId string, text string) *Comment {
	return &Comment{
		ID:             uuid.New(),
		UserId:         userId,
		BlogId:         blogId,
		DateOfCreation: time.Now(),
		Text:           text,
		LastEdit:       time.Now(),
	}
}
