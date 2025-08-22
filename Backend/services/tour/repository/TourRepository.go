package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"tour.xws.com/model"
)

type TourRepository struct {
	Collection *mongo.Collection
}

func (repo *TourRepository) GetAll() ([]model.Tour, error) {
	cursor, err := repo.Collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var tours []model.Tour
	for cursor.Next(context.TODO()) {
		var tour model.Tour
		if err := cursor.Decode(&tour); err != nil {
			return nil, err
		}
		tours = append(tours, tour)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return tours, nil
}

func (repo *TourRepository) GetAllByAuthor(authorId string) ([]model.Tour, error) {
	filter := bson.M{"authorId": authorId} // assuming you store it as "user_id" in MongoDB
	cursor, err := repo.Collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var tours []model.Tour
	for cursor.Next(context.TODO()) {
		var tour model.Tour
		if err := cursor.Decode(&tour); err != nil {
			return nil, err
		}
		tours = append(tours, tour)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return tours, nil
}

func (repo *TourRepository) Create(tour *model.Tour) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := repo.Collection.InsertOne(ctx, tour)
	return err
}

func (repo *TourRepository) Delete(id uuid.UUID) error {
	res, err := repo.Collection.DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("tour not found")
	}
	return nil
}

func (repo *TourRepository) Update(id uuid.UUID, updatedTour model.Tour) error {
	update := bson.M{
		"$set": bson.M{
			"title":       updatedTour.Title,
			"description": updatedTour.Description,
			"difficulty":  updatedTour.Difficulty,
			"tags":        updatedTour.Tags,
			"status":      updatedTour.Status,
			"price":       updatedTour.Price,
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
