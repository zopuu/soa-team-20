package service

import (
	"github.com/google/uuid"
	"tour.xws.com/model"
	"tour.xws.com/repository"
)

type TourService struct {
	TourRepository *repository.TourRepository
}

func (service *TourService) GetAllTours() ([]model.Tour, error) {
	return service.TourRepository.GetAll()
}
func (service *TourService) GetAllByAuthor(userId string) ([]model.Tour, error) {
	return service.TourRepository.GetAllByAuthor(userId)
}

func (service *TourService) Create(tour *model.Tour) error {
	err := service.TourRepository.Create(model.BeforeCreateTour(tour.AuthorId, tour.Title, tour.Description, tour.Tags, tour.Difficulty))
	if err != nil {
		return err
	}
	return nil
}

func (service *TourService) Delete(id uuid.UUID) error {
	return service.TourRepository.Delete(id)
}

func (service *TourService) Update(id uuid.UUID, updatedTour model.Tour) error {
	return service.TourRepository.Update(id, updatedTour)
}

func (service *TourService) GetById(id uuid.UUID) (model.Tour, error) {
	return service.TourRepository.GetById(id)
}
