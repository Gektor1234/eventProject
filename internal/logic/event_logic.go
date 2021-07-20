package logic

import (
	"context"
	"eventApi/internal/app"
	"github.com/CossackPyra/pyraconv"
	"github.com/pkg/errors"
	"log"
	"time"
)

type EventLogic struct {
	eventRepository app.EventRepository
	cache           app.Cache
}

func NewEventLogic(eventRepository app.EventRepository) app.EventLogic {
	e := EventLogic{eventRepository: eventRepository, cache: *app.NewCache()}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	//в случае перезапуска сервиса нужно заполнить кэш всеми незавершеными событиями
	events, err := e.eventRepository.GetList(ctx)
	if err != nil {
		log.Println(err)
	}
	for _, event := range events {
		if event.State == 1 {
			continue
		}
		e.cache.Set(event.Type, event.Id)
	}
	return &e
}

func (e EventLogic) Start(ctx context.Context, eventType string) (err error) {
	//если с таким типом незавершенное событие уже есть, то новое не создаем
	if _, ok := e.cache.Get(eventType); ok {
		return nil
	}
	event := app.Event{
		Type:      eventType,
		StartedAt: pyraconv.ToString(time.Now().Unix()),
	}
	id, err := e.eventRepository.Start(ctx, event)
	if err != nil {
		return errors.Wrap(err, "ошибка в репозитории")
	}
	e.cache.Set(eventType, id)
	return nil
}

func (e EventLogic) Finish(ctx context.Context, eventType string) (err error) {
	//Если незавершенного события такого типа в кэше не найдено, отбиваем 404ую
	eventId, ok := e.cache.Get(eventType)
	if !ok {
		return errors.New("not found")
	}
	err = e.eventRepository.Finish(ctx, eventId)
	if err != nil {
		return err
	}
	//после успешного завершения события,удаляем его из кэша
	e.cache.Delete(eventType)
	return nil
}
