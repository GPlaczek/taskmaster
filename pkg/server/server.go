package server

import (
	"time"
	"net/http"
	"github.com/gin-gonic/gin"
)

type Event struct {
	ID          uint64    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Date        time.Time `json:"date"`
}

var events = []Event{
	{ID: 1, Name: "Wywalić śmiecie", Description: "", Date: time.Date(2024, time.May, 5, 10, 0, 0, 0, time.UTC)},
}

func getEvents(c *gin.Context) {
	c.JSON(http.StatusOK, events)
}

func RunServer() {
	router := gin.Default()
	router.GET("/events", getEvents)
	router.Run("localhost:8080")
}
