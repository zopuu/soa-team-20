package repository

import (
	"github.com/Mihailo84/stakeholders-service/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	DatabaseConnection *gorm.DB
	Collection         *gorm.DB
}

func (repo *UserRepository) GetById(userId int) (*model.User, error) {
	var userById model.User
	result := repo.DatabaseConnection.First(&userById, userId)
	// result := repo.DatabaseConnection.Where(`"Username" = ?`, "Mihailo").First(&userById)
	if result.Error != nil {
		return nil, result.Error
	}
	return &userById, nil
}

func (repo *UserRepository) Update(userId int, updatedUser *model.User) (*model.User, error) {
	var existingUser model.User
	if err := repo.DatabaseConnection.First(&existingUser, userId).Error; err != nil {
		return nil, err
	}

	existingUser.FirstName = updatedUser.FirstName
	existingUser.LastName = updatedUser.LastName
	existingUser.ProfilePhoto = updatedUser.ProfilePhoto
	existingUser.Description = updatedUser.Description
	existingUser.Moto = updatedUser.Moto

	if err := repo.DatabaseConnection.Save(&existingUser).Error; err != nil {
		return nil, err
	}
	return &existingUser, nil
}

func (repo *UserRepository) GetAllUsers() ([]model.User, error) {
	var users []model.User
	if err := repo.DatabaseConnection.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
