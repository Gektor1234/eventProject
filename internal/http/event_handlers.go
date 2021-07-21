package http

import (
	"encoding/json"
	"eventProject/internal/app"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

type eventHandler struct {
	eventLogic app.EventLogic
}

func NewEventHandler(rg *echo.Group, eventLogic app.EventLogic) {
	h := eventHandler{eventLogic: eventLogic}
	rg.POST("/start", h.Start)
	rg.POST("/finish", h.Finish)
}

func (e *eventHandler) Start(c echo.Context) error {
	ctx := c.Request().Context()
	body := c.Request().Body
	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		return errors.Wrap(err, "не удалось считать тело запроса")
	}
	defer body.Close()
	event := app.Event{}
	err = json.Unmarshal(bodyBytes, &event)
	if err != nil {
		return errors.Wrap(err, "не удалось сериализовать тело запроса на структуру")
	}
	err = e.eventLogic.Start(ctx, event.Type)
	if err != nil {
		return errors.Wrap(err, "не удалось создать событие")
	}
	return c.NoContent(http.StatusCreated)
}

func (e *eventHandler) Finish(c echo.Context) error {
	ctx := c.Request().Context()
	body := c.Request().Body
	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		return errors.Wrap(err, "не удалось считать тело запроса")
	}
	defer body.Close()
	event := app.Event{}
	err = json.Unmarshal(bodyBytes, &event)
	if err != nil {
		return errors.Wrap(err, "не удалось сериализовать тело запроса на структуру")
	}
	err = e.eventLogic.Finish(ctx, event.Type)
	if err != nil {
		if err.Error() == "not found" {
			return c.JSON(http.StatusNotFound, "not found")
		}
		return errors.Wrap(err, "не удалось завершить событие")
	}
	return c.NoContent(http.StatusOK)
}
