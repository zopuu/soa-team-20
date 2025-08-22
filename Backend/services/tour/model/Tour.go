package model

import "github.com/google/uuid"

type Tour struct {
	Id          uuid.UUID      `json:"id" bson:"_id,omitempty"`
	AuthorId    string         `json:"authorId" bson:"authorId"`
	Title       string         `json:"title" bson:"title"`
	Description string         `json:"description" bson:"description"`
	Difficulty  TourDifficulty `json:"difficulty" bson:"difficulty"`
	Tags        []string       `json:"tags" bson:"tags"`
	Status      TourStatus     `json:"status" bson:"status"`
	Price       float64        `json:"price" bson:"price"`
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

func BeforeCreateTour(authorId string, title string, description string, tags []string, difficulty TourDifficulty) *Tour {
	return &Tour{
		Id:          uuid.New(),
		AuthorId:    authorId,
		Title:       title,
		Description: description,
		Difficulty:  difficulty,
		Tags:        tags,
		Status:      Draft,
		Price:       0,
	}
}
