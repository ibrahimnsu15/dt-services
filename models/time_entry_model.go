package models

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

type TimeEntry struct {
	Id        bson.ObjectId `json:"id" bson:"_id"`
	CreatedAt time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time     `json:"updated_at" bson:"updated_at"`
	StartTime time.Time     `json:"start_time" bson:"start_time"`
	EndTime   *time.Time    `json:"end_time" bson:"end_time"`
}

func NewTimeEntry() *TimeEntry {
	return &TimeEntry{
		Id:        bson.NewObjectId(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
