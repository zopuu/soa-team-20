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

type LikeRepository struct {
	Collection *mongo.Collection
}

func (repo *LikeRepository) CreateLike(like *model.Like) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := repo.Collection.InsertOne(ctx, like)
	return err
}

func (repo *LikeRepository) GetById(id uuid.UUID) (model.Like, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var like model.Like
	err := repo.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&like)
	if err != nil {
		return like, err
	}
	return like, nil
}

func (repo *LikeRepository) DeleteLike(userId string, blogId string) error {

	res, err := repo.Collection.DeleteOne(context.TODO(), bson.M{"user_id": userId, "blog_id": blogId})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("like not found")
	}
	return nil
}

func (repo *LikeRepository) GetLikesByBlogId(blogId string) (*[]model.Like, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"blog_id": blogId}

	cursor, err := repo.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var likes []model.Like
	if err = cursor.All(ctx, &likes); err != nil {
		return nil, err
	}

	return &likes, nil
}
