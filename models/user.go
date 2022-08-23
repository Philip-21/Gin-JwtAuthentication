package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `bson:"_id"` //declaring the ID as a numeric type in bson format
	First_name    string             `json:"first_name" validate:"required, min=2, max=100"`
	Last_name     string             `json:"last_name" validate:"required, min=2, max=100"`
	Password      string             `json:"Password" validate:"required, min=6"`
	Email         string             `json:"email" validate:"email, required"`
	Phone         string             `json:"phone" validate:"required,  max=11"`
	Token         string             `json:"token"`
	User_type     string             `json:"user_type" validate:"required, eq=ADMIN|eq=USER"` //using the same concept in Enums
	Refresh_token string             `json:"refresh_token"`
	Created_at    time.Time          `json:"created_at"`
	Updated_at    time.Time          `json:"updated_at"`
	User_id       string             `json:"user_id"`
}
