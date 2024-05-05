package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Event struct {
	ID          uint64    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Date        time.Time `json:"date"`
}

var (
	events = []Event{
		{ID: 1, Name: "Wywalić śmiecie", Description: "", Date: time.Date(2024, time.May, 5, 10, 0, 0, 0, time.UTC)},
	}
	evId uint64 = 2
)

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

func RunServer() {
	router := gin.Default()
	router.GET("/events", getEvents)
	router.POST("/events", addEvent)
	router.Run("localhost:8080")
}
