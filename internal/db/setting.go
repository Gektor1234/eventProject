package db

var settings Settings

//главная структура для конфигов в которую при необходимости
//можно дабавить еще подключений
type Settings struct {
	MongoDBConfig MongoDB `json:"mongoDB"`
}
