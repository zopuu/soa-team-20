package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"tour.xws.com/model"
)

type KeyPointRepository struct {
	Collection *mongo.Collection
}

func (repo *KeyPointRepository) GetAll() ([]model.KeyPoint, error) {
	cursor, err := repo.Collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var keyPoints []model.KeyPoint
	for cursor.Next(context.TODO()) {
		var keyPoint model.KeyPoint
		if err := cursor.Decode(&keyPoint); err != nil {
			return nil, err
		}
		keyPoints = append(keyPoints, keyPoint)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return keyPoints, nil
}

func (repo *KeyPointRepository) GetAllByTour(tourId uuid.UUID) ([]model.KeyPoint, error) {
	filter := bson.M{"tourId": tourId}
	cursor, err := repo.Collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var keyPoints []model.KeyPoint
	for cursor.Next(context.TODO()) {
		var keyPoint model.KeyPoint
		if err := cursor.Decode(&keyPoint); err != nil {
			return nil, err
		}
		keyPoints = append(keyPoints, keyPoint)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return keyPoints, nil
}

func (repo *KeyPointRepository) GetAllByTourSortedByCreatedAt(tourId uuid.UUID) ([]model.KeyPoint, error) {
	filter := bson.M{"tourId": tourId}
	opts := options.Find().SetSort(bson.D{{"createdAt", 1}})
	cursor, err := repo.Collection.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var keyPoints []model.KeyPoint
	for cursor.Next(context.TODO()) {
		var keyPoint model.KeyPoint
		if err := cursor.Decode(&keyPoint); err != nil {
			return nil, err
		}
		keyPoints = append(keyPoints, keyPoint)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return keyPoints, nil
}

func (repo *KeyPointRepository) Create(keyPoint *model.KeyPoint) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := repo.Collection.InsertOne(ctx, keyPoint)
	return err
}

func (repo *KeyPointRepository) Delete(id uuid.UUID) error {
	res, err := repo.Collection.DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("keyPoint not found")
	}
	return nil
}

func (repo *KeyPointRepository) Update(id uuid.UUID, updatedKeyPoint model.KeyPoint) error {
	update := bson.M{
		"$set": bson.M{
			"coordinates": updatedKeyPoint.Coordinates,
			"title":       updatedKeyPoint.Title,
			"description": updatedKeyPoint.Description,
			"image":       updatedKeyPoint.Image,
		},
	}

	res, err := repo.Collection.UpdateOne(context.TODO(), bson.M{"_id": id}, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
