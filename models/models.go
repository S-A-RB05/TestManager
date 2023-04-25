package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Strategy struct {
	Id      string    `json:"id"`
	Name    string    `json:"name"`
	Ex      string    `json:"ex`
	Created time.Time `json:"created"`
}

type StrategyRequest struct {
	Id      primitive.ObjectID `bson:"_id,omitempty"`
	Name    string             `bson:"name"`
	Ex      string             `bson:"ex"`
	Created time.Time          `bson:"created"`
}
