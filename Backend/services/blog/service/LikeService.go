package service

import (
	//"fmt"
	"log"
	"time"

	"blog.xws.com/model"
	"blog.xws.com/repository"
)

type LikeService struct {
	LikeRepository *repository.LikeRepository
}

func (service *LikeService) CreateLike(like *model.Like) error {
	like.DateOfCreation = time.Now()

	err := service.LikeRepository.CreateLike(model.CreateNewLike(like.UserId, like.BlogId))
	log.Println("Creating like:", err)
	log.Println("Creating like object:", like)
	if err != nil {
		return err
	}
	return nil
}

func (service *LikeService) Delete(userId string, blogId string) error {
	err := service.LikeRepository.DeleteLike(userId, blogId)

	if err != nil {
		return err
	}
	return nil
}

func (service *LikeService) GetByBlogId(id string) (*[]model.Like, error) {
	return service.LikeRepository.GetLikesByBlogId(id)
}

/*func (service *LikeService) GetById(id uuid.UUID) (*model.Like, error) {
	comment, err := service.LikeRepository.GetById(id)
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("Like with id %s not found", id))
	}

	return &comment, nil
}*/
