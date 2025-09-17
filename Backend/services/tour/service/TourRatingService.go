package service

import (
	"github.com/google/uuid"
	"tour.xws.com/model"
	"tour.xws.com/repository"
)

type TourRatingService struct {
	Repo *repository.TourRatingRepository
}

func (s *TourRatingService) Create(m *model.TourRating) error { return s.Repo.Create(m) }
func (s *TourRatingService) GetByTour(tourId uuid.UUID) ([]model.TourRating, error) {
	return s.Repo.GetByTour(tourId)
}
