package repository

import (
	"context"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"tour.xws.com/model"
)

type TourRatingRepository struct {
	Collection *mongo.Collection
}

func (r *TourRatingRepository) Create(m *model.TourRating) error {
	_, err := r.Collection.InsertOne(context.TODO(), m)
	return err
}

func (r *TourRatingRepository) GetByTour(tourId uuid.UUID) ([]model.TourRating, error) {
	cur, err := r.Collection.Find(context.TODO(), bson.M{"tourId": tourId})
	if err != nil { return nil, err }
	defer cur.Close(context.TODO())

	var out []model.TourRating
	for cur.Next(context.TODO()) {
		var item model.TourRating
		if err := cur.Decode(&item); err != nil { return nil, err }
		out = append(out, item)
	}
	return out, cur.Err()
}
