package model

import (
	"time"

	"github.com/google/uuid"
)

type Blog struct {
	ID             uuid.UUID `json:"id" bson:"_id,omitempty"`
	Title          string    `json:"title" bson:"title"`
	Description    string    `json:"description" bson:"description"`
	DateOfCreation time.Time `json:"date_of_creation" bson:"date_of_creation"`
	Images         []string  `json:"images" bson:"images"`
	Likes          []Like    `json:"likes" bson:"likes"`
}

func CreateNewBlog(title, description string, images []string) *Blog {
	return &Blog{
		ID:             uuid.New(),
		Title:          title,
		Description:    description,
		DateOfCreation: time.Now(),
		Images:         images,
		Likes:          []Like{},
	}
}
