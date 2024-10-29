package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Poll struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	GroupID       primitive.ObjectID `bson:"group_id"`
	Creator       string             `bosn:"creator"`
	Title         string             `bson:"title"`
	FirstOption   string             `bson:"firstOption"`
	Votes1        []string           `bson:"votes1"`
	SecondOption  string             `bson:"secondOption"`
	Votes2        []string           `bson:"votes2"`
	EndTime       time.Time          `bson:"endtime"`
	IsEarlyClosed bool               `bson:"isEarlyClosed"`
}

type CreatePoll struct {
	GroupID      string `json:"groupID" validate:"required"`
	Title        string `json:"title" validate:"required"`
	FirstOption  string `json:"fstOption" validate:"required"`
	SecondOption string `json:"sndOption" validate:"required"`
	Duration     uint64 `json:"duration" validate:"required"`
}

type PollList struct {
	OpenPolls   []PrintPollList
	ClosedPolls []PrintPollList
}

type PrintPollList struct {
	ID               primitive.ObjectID
	Title            string
	Creator          string
	FirstOption      string
	FirstVotesCount  int
	SecondOption     string
	SecondVotesCount int
	EndTime          string
}
type AddVote struct {
	GroupID string `json:"groupID" validate:"required"`
	PollID  string `json:"pollID" validate:"required"`
	Option  string `json:"option" validate:"required"`
}
