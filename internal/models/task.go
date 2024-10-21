package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	GroupID   primitive.ObjectID `bson:"group_id"`
	Title     string             `bson:"title"`
	StartTime time.Time          `bson:"start_time"`
	Duration  int                `bson:"duration"`
	EndTime   time.Time          `bson:"end_time"`
}

type CreateTask struct {
	GroupID   string    `json:"group_id"`
	Title     string    `json:"title"`
	StartTime StartTime `json:"start_time"`
	Duration  Duration  `json:"duration"`
}

type StartTime struct {
	StartDate string `json:"start_date" example:"2024-10-21"`
	StartTime string `json:"start_time" example:"14:00"`
}

type Duration struct {
	DurDays    int `json:"days" example:"0"`
	DurHours   int `json:"hours" example:"2"`
	DurMinutes int `json:"minutes" example:"30"`
}

func (c CreateTask) IsEmptyUpdate() bool {
	return c.Title == "" && c.StartTime.IsFullEmpty() && c.Duration.IsEmpty()
}

func (c CreateTask) IsEmpty() bool {
	return c.Title == "" || c.StartTime.IsEmpty() || c.Duration.IsEmpty()
}
func (s StartTime) IsMissingPart() bool {
	return (s.StartDate == "" && s.StartTime != "") || (s.StartDate != "" && s.StartTime == "")
}

func (s StartTime) IsEmpty() bool {
	return s.StartDate == "" || s.StartTime == ""
}

func (s StartTime) IsFullEmpty() bool {
	return s.StartDate == "" && s.StartTime == ""
}

func (d Duration) IsEmpty() bool {
	return d.DurDays <= 0 && d.DurHours <= 0 && d.DurMinutes <= 0
}
