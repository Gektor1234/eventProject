package app

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Event struct {
	Id         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Type       string             `json:"type,omitempty" bson:"type,omitempty"`
	State      int64              `json:"state,omitempty" bson:"state,omitempty"`
	StartedAt  string             `json:"started_at,omitempty" bson:"started_at,omitempty"`
	FinishedAt string             `json:"finished_at,omitempty" bson:"finished_at,omitempty"`
}

type EventRepository interface {
	Start(ctx context.Context, event Event) (id primitive.ObjectID, err error)
	Finish(ctx context.Context, id primitive.ObjectID) (err error)
	GetList(ctx context.Context) (events []Event, err error)
}

type EventLogic interface {
	Start(ctx context.Context, eventType string) (err error)
	Finish(ctx context.Context, eventType string) (err error)
}
