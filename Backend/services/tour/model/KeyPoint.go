package model

import "github.com/google/uuid"

type KeyPoint struct {
	Id          uuid.UUID   `json:"id" bson:"_id"`
	TourId      uuid.UUID   `json:"tourId" bson:"tourId"`
	Coordinates Coordinates `json:"coordinates" bson:"coordinates"`
	Title       string      `json:"title" bson:"title"`
	Description string      `json:"description" bson:"description"`
	Image       string      `json:"image" bson:"image"`
}

type Coordinates struct {
	Latitude  float64 `json:"latitude" bson:"latitude"`
	Longitude float64 `json:"longitude" bson:"longitude"`
}

func BeforeCreateKeyPoint(tourId uuid.UUID, coordinates Coordinates, title string, description string, image string) *KeyPoint {
	return &KeyPoint{
		Id:          uuid.New(),
		TourId:      tourId,
		Coordinates: coordinates,
		Title:       title,
		Description: description,
		Image:       image,
	}
}
