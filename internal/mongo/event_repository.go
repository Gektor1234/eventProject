package mongo

import (
	"context"
	"eventProject/internal/app"
	"github.com/CossackPyra/pyraconv"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const (
	dbName         = "local"
	collectionName = "events"
)

type EventRepository struct {
	mongoClient *mongo.Client
}

func NewEventRepository(mongoClient *mongo.Client) app.EventRepository {
	return &EventRepository{mongoClient: mongoClient}
}

func (e EventRepository) GetList(ctx context.Context) (events []app.Event, err error) {
	events = []app.Event{}
	collection := e.mongoClient.Database(dbName).Collection(collectionName)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, errors.Wrap(err, "ошибка при попытке получить все события из монго")
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var event app.Event
		cursor.Decode(&event)
		events = append(events, event)
	}
	if err = cursor.Err(); err != nil {
		return nil, errors.Wrap(err, "ошибка курсора при мапинге на структуру всех событий")
	}
	return events, nil
}

func (e EventRepository) Start(ctx context.Context, event app.Event) (id primitive.ObjectID, err error) {
	collection := e.mongoClient.Database(dbName).Collection(collectionName)
	result, err := collection.InsertOne(ctx, event)
	if err != nil {
		return id, errors.Wrapf(err, "ошибка при создании события для типа: %s", event.Type)
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

func (e EventRepository) Finish(ctx context.Context, id primitive.ObjectID) (err error) {
	collection := e.mongoClient.Database(dbName).Collection(collectionName)
	_, err = collection.UpdateOne(ctx, bson.M{"_id": id},
		bson.D{
			{"$set", bson.D{{"state", 1}, {"finished_at", pyraconv.ToString(time.Now().Unix())}},
			},
		})
	if err != nil {
		return errors.Wrapf(err, "ошибка при попытке завершить событие с id: %v", id)
	}
	return nil
}
