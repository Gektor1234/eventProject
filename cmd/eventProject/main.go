package main

import (
	"eventProject/internal/db"
	"eventProject/internal/http"
	"eventProject/internal/logic"
	"eventProject/internal/mongo"
	"github.com/labstack/echo"
	"log"
)

func main() {
	e := echo.New()
	db.ConfigInit("configs/local_settings.json")
	mongoClient := db.InitMongoDBConnection()
	eventRepository := mongo.NewEventRepository(mongoClient)
	eventLogic, err := logic.NewEventLogic(eventRepository)
	if err != nil {
		log.Fatal(err, "не удалось запустить конструктор событий")
	}
	http.NewEventHandler(e.Group("/v1"), eventLogic)
	go func() {
		log.Println("запуск http сервера на порту 6655")
		err := e.Start(":6655")
		if err != nil {
			log.Fatal(err, "не удалось запустить сервис echo")
		}
	}()
	select {}
}
