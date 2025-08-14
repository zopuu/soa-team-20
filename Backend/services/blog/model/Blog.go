package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Blog struct {
	ID             uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;column:Id;"`
	Title          string    `json:"title" gorm:"column:Title;type:varchar(255);not null"`
	Description    string    `json:"description" gorm:"column:Description;type:varchar(255)"`
	DateOfCreation time.Time `json:"date_of_creation" gorm:"column:DateOfCreation;type:date"`
	Images         []string  `json:"images" gorm:"column:Images;type:text[]"`
	Likes          []Like    `json:"likes" gorm:"column:Likes;type:text[]"`
}

func (blog *Blog) BeforeCreate(scope *gorm.DB) error {
	blog.ID = uuid.New()
	return nil
}
