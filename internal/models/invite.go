package models

import (
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateInvite struct {
	GroupID string `json:"group_id" validate:"required"`
	User    string `json:"user" validate:"required"`
}

type InvitationToken struct {
	UserLogin string
	GroupID   string
	jwt.StandardClaims
}

type Invitation struct {
	Invite_ID primitive.ObjectID `bson:"_id,omitempty"`
	Receiver  string             `bson:"receiver"`
	Sender    string             `bson:"sender"`
	GroupID   primitive.ObjectID `bson:"group_id"`
	GroupName string             `bson:"group_name"`
	Token     string             `bson:"token"`
	IsUsed    bool               `bson:"isUsed"`
}

type InvitationList struct {
	Invite_ID      primitive.ObjectID
	InvitationText string
	InvitationLink string
}
