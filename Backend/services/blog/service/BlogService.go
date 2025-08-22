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

func (service *BlogService) GetAllBlogs() ([]model.Blog, error) {
	return service.BlogRepository.GetAll()
}

func (service *BlogService) GetById(id uuid.UUID) (*model.Blog, error) {
	blog, err := service.BlogRepository.GetById(id)
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("Blog with id %s not found", id))
	}

	return &blog, nil
}

func (service *BlogService) Create(blog *model.Blog) error {
	err := service.BlogRepository.Create(model.BeforeCreate(blog.UserId, blog.Title, blog.Description, blog.Images))
	if err != nil {
		return err
	}
	return nil
}

func (service *BlogService) Delete(id uuid.UUID) error {
	return service.BlogRepository.Delete(id)
}

func (service *BlogService) Update(id uuid.UUID, updatedBlog model.Blog) error {
	return service.BlogRepository.Update(id, updatedBlog)
}

func (service *BlogService) GetAllByUser(userId string) ([]model.Blog, error) {
	return service.BlogRepository.GetAllByUser(userId)
}
