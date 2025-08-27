package service

import (
	"fmt"
	"time"

	"blog.xws.com/model"
	"blog.xws.com/repository"
	"github.com/google/uuid"
)

type CommentService struct {
	CommentRepository *repository.CommentRepository
}

func (service *CommentService) CreateComment(comment *model.Comment) error {
	comment.DateOfCreation = time.Now()
	err := service.CommentRepository.CreateComment(comment)
	if err != nil {
		return err
	}
	return nil
}

func (service *CommentService) UpdateComment(comment *model.Comment) error {
	err := service.CommentRepository.UpdateComment(comment.ID, comment)
	if err != nil {
		return err
	}
	return nil
}

func (service *CommentService) Delete(id uuid.UUID) error {
	err := service.CommentRepository.DeleteComment(id)

	if err != nil {
		return err
	}
	return nil
}

func (service *CommentService) GetByBlogId(id uuid.UUID) (*[]model.Comment, error) {
	return service.CommentRepository.GetCommentsByBlogId(id)
}

func (service *CommentService) GetById(id uuid.UUID) (*model.Comment, error) {
	comment, err := service.CommentRepository.GetById(id)
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("Comment with id %s not found", id))
	}

	return &comment, nil
}
