package models

import (
	"time"
)

type Strategy struct {
	Name    string    `json:"name"`
	Mq      string    `json:"mq"`
	Ex      string    `json:"ex`
	Created time.Time `json:"created"`
}

type StrategyRequest struct {
	S_id     string    `json:"id"`
	Name    string    `json:"name"`
	Ex      string    `json:"ex"`
	Created time.Time `json:"created"`
}
