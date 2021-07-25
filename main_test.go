package main

import (
	"bytes"
	"context"
	"eventProject/internal/app"
	"eventProject/internal/db"
	mongo2 "eventProject/internal/mongo"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"testing"
)

const (
	finishUrl              = "http://127.0.0.1:6655/v1/finish"
	startUrl               = "http://127.0.0.1:6655/v1/start"
	eventTypeForStartCase  = "testType1234"
	eventTypeForFinishCase = "testType123"
	noDocumentsError       = "mongo: no documents in result"
)

type testCases struct {
	name          string
	wantCode      int64
	wantEventType string
	wantState     int64
}

func testGetAndDeleteEvent(client *mongo.Client, eventType string) (event app.Event) {
	collection := client.Database(mongo2.DbName).Collection(mongo2.CollectionName)
	event = app.Event{}
	err := collection.FindOne(context.TODO(), bson.D{{"type", eventType}}).Decode(&event)
	if err != nil {
		if err.Error() != noDocumentsError {
			log.Fatal(err)
		}
	}
	_, err = collection.DeleteOne(context.TODO(), bson.D{{"type", eventType}})
	if err != nil {
		if err.Error() != noDocumentsError {
			log.Fatal(err)
		}
	}
	return event
}

func TestApi(t *testing.T) {
	db.ConfigInit("configs/local_settings.json")
	client := db.InitMongoDBConnection()
	testCasesStart := []testCases{{
		name:          "newEvent",
		wantCode:      http.StatusCreated,
		wantEventType: eventTypeForStartCase,
	}, {
		name:          "existEvent",
		wantCode:      http.StatusCreated,
		wantEventType: "",
	}}
	testCasesFinish := []testCases{{
		name:          "finishUnfinishedEvent",
		wantCode:      http.StatusOK,
		wantEventType: eventTypeForFinishCase,
		wantState:     1,
	}, {
		name:     "finishNonExistentEvent",
		wantCode: http.StatusNotFound,
	}}
	for _, tc := range testCasesStart {
		t.Run(tc.name, func(t *testing.T) {
			r := bytes.NewReader([]byte(`{"type" : "` + eventTypeForStartCase + `"}`))
			resp, err := http.Post(startUrl, "application/json", r)
			if err != nil {
				log.Fatal(err, "cant send request")
			}
			event := testGetAndDeleteEvent(client, eventTypeForStartCase)
			assert.Equal(t, tc.wantCode, int64(resp.StatusCode))
			assert.Equal(t, tc.wantEventType, event.Type)
		})
	}
	for i, tc := range testCasesFinish {
		t.Run(tc.name, func(t *testing.T) {
			r := bytes.NewReader([]byte(`{"type" : "` + eventTypeForFinishCase + `"}`))
			if i == 0 {
				_, err := http.Post(startUrl, "application/json", r)
				if err != nil {
					log.Fatal(err, "cant send request")
				}
			}
			r = bytes.NewReader([]byte(`{"type" : "` + eventTypeForFinishCase + `"}`))
			resp, err := http.Post(finishUrl, "application/json", r)
			if err != nil {
				log.Fatal(err, "cant send request")
			}
			event := testGetAndDeleteEvent(client, eventTypeForFinishCase)
			assert.Equal(t, tc.wantCode, int64(resp.StatusCode))
			assert.Equal(t, tc.wantEventType, event.Type)
			assert.Equal(t, tc.wantState, event.State)
		})
	}
}
