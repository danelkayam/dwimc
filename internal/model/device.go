package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Device struct {
	ID        bson.ObjectID `json:"id" bson:"_id"`
	CreatedAt time.Time     `json:"created_at" bson:"createdAt"`
	UpdatedAt time.Time     `json:"updated_at" bson:"updatedAt"`
	Serial    string        `json:"serial" bson:"serial"`
	Name      string        `json:"name" bson:"name"`
}
