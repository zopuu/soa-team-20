package model

import (
	"time"

	"github.com/google/uuid"
)

type Tour struct {
	Id            uuid.UUID      `json:"id" bson:"_id,omitempty"`
	AuthorId      string         `json:"authorId" bson:"authorId"`
	Title         string         `json:"title" bson:"title"`
	Description   string         `json:"description" bson:"description"`
	Difficulty    TourDifficulty `json:"difficulty" bson:"difficulty"`
	Tags          []string       `json:"tags" bson:"tags"`
	Status        TourStatus     `json:"status" bson:"status"`
	Price         float64        `json:"price" bson:"price"`
	Distance      float64        `json:"distance" bson:"distance"`
	PublishedAt   time.Time      `json:"publishedAt" bson:"publishedAt"`
	ArchivedAt    time.Time      `json:"archivedAt" bson:"archivedAt"`
	Duration      float64        `json:"duration" bson:"duration"`
	TransportType TransportType  `json:"transportType" bson:"transportType"`
}

type TourStatus int

const (
	Draft TourStatus = iota
	Published
	Archived
)

type TourDifficulty int

const (
	Beginner TourDifficulty = iota
	Intermediate
	Advanced
	Pro
)

type TransportType int

const (
	Walking TransportType = iota
	Bicycle
	Bus
)

func BeforeCreateTour(authorId string, title string, description string, tags []string, difficulty TourDifficulty) *Tour {
	return &Tour{
		Id:            uuid.New(),
		AuthorId:      authorId,
		Title:         title,
		Description:   description,
		Difficulty:    difficulty,
		Tags:          tags,
		Status:        Draft,
		Price:         0,
		Distance:      0,
		PublishedAt:   time.Time{},
		ArchivedAt:    time.Time{},
		Duration:      0,
		TransportType: 0,
	}
}
