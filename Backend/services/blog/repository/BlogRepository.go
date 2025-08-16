package repository

import (
	"context"
	"time"

	"blog.xws.com/model"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type BlogRepository struct {
	Collection *mongo.Collection
}

func (repo *BlogRepository) FindById(id uuid.UUID) (model.Blog, error) {
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
