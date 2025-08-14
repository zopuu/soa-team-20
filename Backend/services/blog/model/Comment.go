package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Comment struct {
	ID             uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;column:Id;"`
	UserId         uuid.UUID `json:"userId" gorm:"column:UserId;type:uuid;not null"`
	DateOfCreation time.Time `json:"dateOfCreation" gorm:"column:DateOfCreation;type:date"`
	Text           string    `json:"text" gorm:"column:Text;type:varchar(255);not null"`
	LastEdit       time.Time `json:"lastEdit" gorm:"column:LastEdit;type:date"`
}

func (comment *Comment) BeforeCreate(scope *gorm.DB) error {
	comment.ID = uuid.New()
	return nil
}
