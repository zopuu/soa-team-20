package service

import (
	"tour.xws.com/model"
	"tour.xws.com/repository"
)

type CurrentLocationService struct {
	Repo *repository.CurrentLocationRepository
}

func (s *CurrentLocationService) Get(userId string) (*model.CurrentLocation, error) {
	return s.Repo.GetByUserId(userId)
}

func (s *CurrentLocationService) Set(userId string, coords model.Coordinates) error {
	return s.Repo.Upsert(userId, coords)
}
