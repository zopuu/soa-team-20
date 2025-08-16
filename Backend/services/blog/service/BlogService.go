package service

import (
	"fmt"

	"blog.xws.com/model"
	"blog.xws.com/repository"
	"github.com/google/uuid"
)

type BlogService struct {
	BlogRepository *repository.BlogRepository
}

func (service *BlogService) FindById(id uuid.UUID) (*model.Blog, error) {
	blog, err := service.BlogRepository.FindById(id)
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("Student with id %s not found", id))
	}

	return &blog, nil
}

func (service *BlogService) Create(blog *model.Blog) error {
	err := service.BlogRepository.Create(blog)
	if err != nil {
		return err
	}
	return nil
}
