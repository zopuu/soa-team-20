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

type CommentRepository struct {
	Collection *mongo.Collection
}

func (repo *CommentRepository) CreateComment(comment *model.Comment) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := repo.Collection.InsertOne(ctx, comment)
	return err
}

func (repo *CommentRepository) GetById(id uuid.UUID) (model.Comment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var comment model.Comment
	err := repo.Collection.FindOne(ctx, bson.M{"blog_id": id}).Decode(&comment)
	if err != nil {
		return comment, err
	}
	return comment, nil
}

func (repo *CommentRepository) UpdateComment(id uuid.UUID, comment *model.Comment) error {
	update := bson.M{
		"$set": bson.M{
			"text":      comment.Text,
			"last_edit": time.Now(),
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

func (repo *CommentRepository) DeleteComment(id uuid.UUID) error {

	res, err := repo.Collection.DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("comment not found")
	}
	return nil
}

func (repo *CommentRepository) GetCommentsByBlogId(blogId uuid.UUID) (*[]model.Comment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"blog_id": blogId.String()}
	//filter := bson.M{}
	cursor, err := repo.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var comments []model.Comment
	println("Repository comments count", comments)
	if err = cursor.All(ctx, &comments); err != nil {
		return nil, err
	}

	return &comments, nil
}
