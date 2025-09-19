package model

import (
	"time"

	"github.com/google/uuid"
)

type TourRating struct {
  Id           uuid.UUID `json:"id" bson:"_id,omitempty"`
  TourId       uuid.UUID `json:"tourId" bson:"tourId"`
  Rating       int       `json:"rating" bson:"rating"`
  Comment      string    `json:"comment,omitempty" bson:"comment,omitempty"`
  TouristName  string    `json:"touristName,omitempty" bson:"touristName,omitempty"`
  TouristEmail string    `json:"touristEmail,omitempty" bson:"touristEmail,omitempty"`
  VisitedAt    time.Time `json:"visitedAt,omitempty" bson:"visitedAt,omitempty"`
  CommentedAt  time.Time `json:"commentedAt,omitempty" bson:"commentedAt,omitempty"`
  CreatedAt    time.Time `json:"createdAt" bson:"createdAt"`

  Images       []string  `json:"images,omitempty" bson:"images,omitempty"`
}

