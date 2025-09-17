package model

import "time"

type CurrentLocation struct {
	UserId      string      `json:"userId" bson:"_id"`
	Coordinates Coordinates `json:"coordinates" bson:"coordinates"`
	UpdatedAt   time.Time   `json:"updatedAt" bson:"updatedAt"`
}
