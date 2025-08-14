package model

import "github.com/google/uuid"

type Like struct {
	UserId uuid.UUID `json:"userId"`
	BlogId uuid.UUID `json:"blogId"`
}
