package db

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type MongoDB struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	DBName   string `json:"dbName"`
	Password string `json:"password"`
}


// ConfigInit считывает конфигурационный файл, находящийся по пути filepath
func ConfigInit(filepath string) {
	log.Println("Считывание конфигурационного файла по пути ", filepath, "...")
	cfg, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Println("File error: ", err.Error())
		os.Exit(1)
	}
	err = json.Unmarshal(cfg, &settings)
	if err != nil {
		log.Println("Произошла ошибка при попытке преобразовать конфиг файл из JSON: ", err.Error())
	} else {
		log.Println("Файл конфигурации успешно считан")
	}
}

// InitMongoDBConnection создает клиента для монго, а также проверяет что подключение прошло без ошибок
func InitMongoDBConnection() *mongo.Client {
	cfg := settings.MongoDBConfig
	client, err := mongo.NewClient(options.Client().ApplyURI(fmt.Sprintf("%s://%s:%s", cfg.DBName, cfg.Host,
		cfg.Port)))
	if err != nil {
		log.Println("произошла ошибка при попытке инициализации клиента для mongo: ", err.Error())
		return nil
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Println("произошла ошибка при попытке инициализации соединения с mongo: ", err.Error())
		return nil
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Println("произошла ошибка при попытке проверки соединения с mongo: ", err.Error())
		return nil
	}
	return client
}
