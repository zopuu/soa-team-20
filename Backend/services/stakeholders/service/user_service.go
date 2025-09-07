package service

import (
	"github.com/Mihailo84/stakeholders-service/model"
	"github.com/Mihailo84/stakeholders-service/repository"
)

type UserService struct {
	UserRepository *repository.UserRepository
}

func (service *UserService) Update(userId int, updatedUser *model.User) (*model.User, error) {
	return service.UserRepository.Update(userId, updatedUser)
}

func (service *UserService) GetById(userId int) (*model.User, error) {
	return service.UserRepository.GetById(userId)
}

func (service *UserService) GetAllUsers() ([]model.User, error) {
	return service.UserRepository.GetAllUsers()
}
