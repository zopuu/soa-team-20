package model

import "time"

type OrderItem struct {
	TourID string  `bson:"tour_id" json:"tour_id"`
	Name   string  `bson:"name" json:"name"`
	Price  float64 `bson:"price" json:"price"`
}

type Cart struct {
	UserID string      `bson:"user_id" json:"user_id"`
	Items  []OrderItem `bson:"items" json:"items"`
	// Total is not stored, computed on read
}

type TourPurchaseToken struct {
	ID          string    `bson:"_id" json:"id"`
	UserID      string    `bson:"user_id" json:"user_id"`
	TourID      string    `bson:"tour_id" json:"tour_id"`
	TourName    string    `bson:"tour_name" json:"tour_name"`
	Price       float64   `bson:"price" json:"price"`
	PurchasedAt time.Time `bson:"purchased_at" json:"purchased_at"`
}
