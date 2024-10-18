package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Login        string             `json:"login" validate:"required,min=6" bson:"login"`
	Email        string             `json:"email" validate:"required,min=6" bson:"email"`
	Password     string             `json:"password,omitempty" validate:"required,min=6" bson:"-"`
	PasswordHash string             `json:"-" bson:"hashed_password"`
}

type SignUp struct {
	Login    string `json:"login" validate:"required,min=6,max=15"`
	Email    string `json:"email" validate:"required,min=6,max=20"`
	Password string `json:"password,omitempty" validate:"required,min=6,max=30"`
}
type LoginRequest struct {
	Option   string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}
