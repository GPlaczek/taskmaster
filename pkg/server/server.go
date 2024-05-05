package server

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

var (
	events = []Event{
		{ID: 1, Name: "Wywalić śmiecie", Description: "", Date: time.Date(2024, time.May, 5, 10, 0, 0, 0, time.UTC)},
	}
	evId uint64 = 2
)

type ETag interface {
	update(tag [20]byte)
}

type Event struct {
	ID          uint64    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Date        time.Time `json:"date"`
	ETag        [20]byte  `json:"-"`
}

func (e *Event) update(tag [20]byte) {
	e.ETag = tag
}

func MarshalAndTag[ET ETag](inst ET) ([]byte, error) {
	d, err := json.Marshal(inst)
	if err != nil {
		return nil, err
	}

	t := sha1.Sum(d)
	inst.update(t)

	return d, nil
}

func getEvents(c *gin.Context) {
	c.JSON(http.StatusOK, events)
}

func addEvent(c *gin.Context) {
	var ev Event

	if err := c.BindJSON(&ev); err != nil {
		c.Status(http.StatusUnprocessableEntity)
		return
	}

	ev.ID = evId
	evId++

	events = append(events, ev)

	c.JSON(http.StatusOK, gin.H{"id": ev.ID})
}

func getEvent(c *gin.Context) {
	_id := c.Param("id")
	id, err := strconv.ParseUint(_id, 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	var event *Event = nil
	for i, _ := range events {
		if events[i].ID == id {
			event = &events[i]
		}
	}

	if event == nil {
		c.Status(http.StatusNotFound)
		return
	}

	d, err := MarshalAndTag[*Event](event)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Header("ETag", fmt.Sprintf("%x", event.ETag))
	c.Data(http.StatusOK, "application/json", d)
}

func RunServer() {
	router := gin.Default()
	router.GET("/events", getEvents)
	router.POST("/events", addEvent)
	router.GET("/events/:id", getEvent)
	router.Run("localhost:8080")
}
