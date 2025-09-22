package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"tour.xws.com/model"
)

type TourExecutionRepository struct {
	Collection *mongo.Collection
}

func (r *TourExecutionRepository) Insert(te *model.TourExecution) error {
	_, err := r.Collection.InsertOne(context.TODO(), te)
	return err
}

func (r *TourExecutionRepository) GetActiveByUser(userId string) (*model.TourExecution, error) {
	var te model.TourExecution
	err := r.Collection.FindOne(context.TODO(), bson.M{"userId": userId, "status": model.TourExecActive}).Decode(&te)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &te, err
}

func (r *TourExecutionRepository) GetById(id uuid.UUID) (*model.TourExecution, error) {
	var te model.TourExecution
	err := r.Collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&te)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &te, err
}

func (r *TourExecutionRepository) Update(te *model.TourExecution) error {
	_, err := r.Collection.ReplaceOne(context.TODO(), bson.M{"_id": te.Id}, te)
	return err
}

func (r *TourExecutionRepository) SetStatus(id uuid.UUID, status model.TourExecutionStatus, endedAt *time.Time) error {
	update := bson.M{"$set": bson.M{"status": status}}
	if endedAt != nil {
		update["$set"].(bson.M)["endedAt"] = endedAt
	}
	_, err := r.Collection.UpdateOne(context.TODO(), bson.M{"_id": id}, update)
	return err
}

func (r *TourExecutionRepository) GetActiveByUserAndTour(userId string, tourId uuid.UUID) (*model.TourExecution, error) {
	var te model.TourExecution
	err := r.Collection.FindOne(
		context.TODO(),
		bson.M{"userId": userId, "tourId": tourId, "status": model.TourExecActive},
	).Decode(&te)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &te, err
}