package models

import "github.com/dgrijalva/jwt-go"

type CreateInvite struct {
	GroupID string `json:"group_id" validate:"required"`
	User    string `json:"user" validate:"required"`
}

type InvitationToken struct {
	UserLogin string
	GroupID   string
	jwt.StandardClaims
}

type Invitation struct{
	Receiver string `bson:"receiver"`
	Sender	string	`bson:"sender"`
	GroupName string `bson:"group_name"`
	Token string `bson:"token"`
	IsUsed bool `bson:"isUsed"`
}

type InvitationList struct{
	InvitationText string
	InvitationLink string
}