package model

import "time"

type User struct {
	Id           int       `gorm:"column:Id;primaryKey"`
	Username     string    `gorm:"column:Username"`
	Email        string    `gorm:"column:Email"`
	PasswordHash string    `gorm:"column:PasswordHash"`
	Role         string    `gorm:"column:Role"`
	CreatedAt    time.Time `gorm:"column:CreatedAt"`
	Description  string    `gorm:"column:Description"`
	FirstName    string    `gorm:"column:FirstName"`
	LastName     string    `gorm:"column:LastName"`
	Moto         string    `gorm:"column:Moto"`
	ProfilePhoto string    `gorm:"column:ProfilePhoto`
}

func (User) TableName() string { return "Users" }
