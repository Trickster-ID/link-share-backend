package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type RefreshTokenSession struct {
	ObjectID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	RefreshToken string             `json:"refresh_token" bson:"refresh_token"`
	Expired      time.Time          `json:"exp" bson:"exp"`
	UserData     *UserDataOnJWT     `json:"user_data" bson:"user_data"`
}
