package repository

import (
	"context"
	"errors"
	"time"

	"blog.xws.com/model"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type BlogRepository struct {
	Collection *mongo.Collection
}

func (repo *BlogRepository) GetAll() ([]model.Blog, error) {
	cursor, err := repo.Collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var blogs []model.Blog
	for cursor.Next(context.TODO()) {
		var blog model.Blog
		if err := cursor.Decode(&blog); err != nil {
			return nil, err
		}
		blogs = append(blogs, blog)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return blogs, nil
}

func (repo *BlogRepository) GetById(id uuid.UUID) (model.Blog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var blog model.Blog
	err := repo.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&blog)
	if err != nil {
		return blog, err
	}
	return blog, nil
}

func (repo *BlogRepository) Create(blog *model.Blog) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := repo.Collection.InsertOne(ctx, blog)
	return err
}

func (repo *BlogRepository) Delete(id uuid.UUID) error {
	res, err := repo.Collection.DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("blog not found")
	}
	return nil
}

func (repo *BlogRepository) Update(id uuid.UUID, updatedBlog model.Blog) error {
	update := bson.M{
		"$set": bson.M{
			"title":            updatedBlog.Title,
			"description":      updatedBlog.Description,
			"images":           updatedBlog.Images,
			"date_of_creation": updatedBlog.DateOfCreation, // optional
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

func (repo *BlogRepository) GetAllByUser(userID string) ([]model.Blog, error) {
	filter := bson.M{"userId": userID}
	cursor, err := repo.Collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	println("Cursor obtained for userID:", userID)
	println("NUmber of documents found:", cursor.RemainingBatchLength())
	var blogs []model.Blog
	for cursor.Next(context.TODO()) {
		var blog model.Blog
		if err := cursor.Decode(&blog); err != nil {
			return nil, err
		}
		blogs = append(blogs, blog)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return blogs, nil
}
