package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Strategy struct {
	Id      string    `json:"id"`
	UserId  string    `json:"userid"`
	Name    string    `json:"name"`
	Ex      string    `json:"ex`
	Created time.Time `json:"created"`
}

type StrategyRequest struct {
	Id      primitive.ObjectID `bson:"_id,omitempty"`
	UserId  string             `json:"userid"`
	Name    string             `bson:"name"`
	Ex      string             `bson:"ex"`
	Created time.Time          `bson:"created"`
}

type Test struct {
	Id      string    `json:"id"`
	UserId  string    `json:"userid"`
	StratId string    `json:"strat_id"`
	Created time.Time `json:"created"`
}
