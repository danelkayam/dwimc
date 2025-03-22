package model

import (
	"time"
)

type Device struct {
	ID        string    `json:"id" bson:"_id"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
	Serial    string    `json:"serial" bson:"serial"`
	Name      string    `json:"name" bson:"name"`
}
