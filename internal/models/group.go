package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type CreateGroup struct {
	Name        string   `json:"name" validate:"required"`
	Invitations []string `json:"invites"`
}

type GroupList struct {
	ID           primitive.ObjectID
	Name         string
	MembersCount int
}

type Group struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	LeaderLogin string             `json:"leader_login" bson:"leader_login"`
	Members     []string           `json:"members" bson:"members"`
	IsActive    bool               `json:"-" bson:"isActive"`
}

type BlackList struct {
	GroupID   primitive.ObjectID `bson:"group_id"`
	Blacklist []string           `bson:"blacklist"`
}
