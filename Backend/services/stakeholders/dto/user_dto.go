package dto

import (
	"time"

	"github.com/google/uuid"
)

type UserUpdateRequest struct {
	Description  string    `json:"description" gorm:"primaryKey"`
	FirstName    string    `json:"FirstName" gorm:"primaryKey"`
	LastName     string    `json:"LastName" gorm:"primaryKey"`
	Moto         string    `json:"Moto" gorm:"primaryKey"`
	ProfilePhoto string    `json:"ProfilePhoto" gorm:"primaryKey"`

}

type UserResponse struct {
	Id 			 uuid.UUID `json:"id" gorm:"primaryKey"`
	Username     string    `json:"username" gorm:"primaryKey"`
	Email        string    `json:"email" gorm:"primaryKey"`
	PasswordHash string    `json:"passwordHash" gorm:"primaryKey"`
	Role         string    `json:"role" gorm:"primaryKey"`
	CreatedAt    time.Time `json:"createdAt" gorm:"primaryKey"`
	UpdatedAt    time.Time `json:"updatedAt" gorm:"primaryKey"`
	Description  string    `json:"description" gorm:"primaryKey"`
	FirstName    string    `json:"FirstName" gorm:"primaryKey"`
	LastName     string    `json:"LastName" gorm:"primaryKey"`
	Moto         string    `json:"Moto" gorm:"primaryKey"`
	ProfilePhoto string    `json:"ProfilePhoto" gorm:"primaryKey"`
}