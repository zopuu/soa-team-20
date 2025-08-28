package service

import (
	"github.com/google/uuid"
	"tour.xws.com/model"
	"tour.xws.com/repository"
)

type KeyPointService struct {
	KeyPointRepository *repository.KeyPointRepository
}

func (service *KeyPointService) GetAllKeyPoints() ([]model.KeyPoint, error) {
	return service.KeyPointRepository.GetAll()
}

func (service *KeyPointService) GetAllByTour(tourId uuid.UUID) ([]model.KeyPoint, error) {
	return service.KeyPointRepository.GetAllByTour(tourId)
}

func (service *KeyPointService) GetAllByTourSortedByCreatedAt(tourId uuid.UUID) ([]model.KeyPoint, error) {
	return service.KeyPointRepository.GetAllByTourSortedByCreatedAt(tourId)
}

func (service *KeyPointService) Create(keyPoint *model.KeyPoint) error {
	err := service.KeyPointRepository.Create(model.BeforeCreateKeyPoint(keyPoint.TourId, keyPoint.Coordinates, keyPoint.Title, keyPoint.Description, keyPoint.Image))
	if err != nil {
		return err
	}
	return nil
}

func (service *KeyPointService) Delete(id uuid.UUID) error {
	return service.KeyPointRepository.Delete(id)
}

func (service *KeyPointService) Update(id uuid.UUID, updatedKeyPoint model.KeyPoint) error {
	return service.KeyPointRepository.Update(id, updatedKeyPoint)
}
