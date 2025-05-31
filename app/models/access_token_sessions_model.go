package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type AccessTokenSession struct {
	ObjectID    primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	AccessToken string             `json:"access_token" bson:"access_token"`
	Expired     time.Time          `json:"exp" bson:"exp"`
	UserData    *UserDataOnJWT     `json:"user_data" bson:"user_data"`
}
