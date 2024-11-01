package models

import "time"

type Message struct {
	User    string    `json:"user_login" bson:"user_login"`
	Content string    `json:"content" bson:"content"`
	GroupID string    `json:"group_id" bson:"group_id"`
	Time    time.Time `bson:"time"`
}

type UserConn struct {
	UserLogin string
	GroupID   string
}
