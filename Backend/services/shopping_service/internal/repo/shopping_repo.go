package repo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/zopuu/soa-team-20/Backend/services/shopping_service/internal/model"
)

type ShoppingRepo struct {
	carts  *mongo.Collection
	tokens *mongo.Collection
}

func NewShoppingRepo(db *mongo.Database) *ShoppingRepo {
	return &ShoppingRepo{
		carts:  db.Collection("carts"),
		tokens: db.Collection("tokens"),
	}
}

// GetCart returns the cart for a user, or empty cart if none
func (r *ShoppingRepo) GetCart(ctx context.Context, userID string) (*model.Cart, error) {
	var cart model.Cart
	err := r.carts.FindOne(ctx, bson.M{"user_id": userID}).Decode(&cart)
	if err == mongo.ErrNoDocuments {
		return &model.Cart{UserID: userID, Items: []model.OrderItem{}}, nil
	}
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *ShoppingRepo) SaveCart(ctx context.Context, cart *model.Cart) error {
	ctx2, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_, err := r.carts.UpdateOne(ctx2, bson.M{"user_id": cart.UserID},
		bson.M{"$set": bson.M{"items": cart.Items}},
		options.Update().SetUpsert(true))
	return err
}

func (r *ShoppingRepo) ClearCart(ctx context.Context, userID string) error {
	_, err := r.carts.DeleteOne(ctx, bson.M{"user_id": userID})
	return err
}

func (r *ShoppingRepo) SaveTokens(ctx context.Context, tokens []model.TourPurchaseToken) error {
	if len(tokens) == 0 {
		return nil
	}
	docs := make([]interface{}, 0, len(tokens))
	for i := range tokens {
		docs = append(docs, tokens[i])
	}
	_, err := r.tokens.InsertMany(ctx, docs)
	return err
}

func (r *ShoppingRepo) GetTokens(ctx context.Context, userID string) ([]model.TourPurchaseToken, error) {
	cur, err := r.tokens.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var res []model.TourPurchaseToken
	for cur.Next(ctx) {
		var t model.TourPurchaseToken
		if err := cur.Decode(&t); err != nil {
			return nil, err
		}
		res = append(res, t)
	}
	return res, nil
}

// Optional index creation helper
func (r *ShoppingRepo) EnsureIndexes(ctx context.Context) error {
	// index tokens.user_id
	_, err := r.tokens.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "user_id", Value: 1}},
	})
	return err
}
