package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"tour.xws.com/model"
)

type CurrentLocationRepository struct {
	Collection *mongo.Collection
}

func (r *CurrentLocationRepository) GetByUserId(userId string) (*model.CurrentLocation, error) {
	var loc model.CurrentLocation
	err := r.Collection.FindOne(context.TODO(), bson.M{"_id": userId}).Decode(&loc)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &loc, err
}

func (r *CurrentLocationRepository) Upsert(userId string, coords model.Coordinates) error {
	now := time.Now()
	_, err := r.Collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": userId},
		bson.M{"$set": bson.M{
			"coordinates": coords,
			"updatedAt":   now,
		}},
		&options.UpdateOptions{Upsert: boolPtr(true)},
	)
	return err
}

func boolPtr(b bool) *bool { return &b }
