package server

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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

	_, err := MarshalAndTag(&ev)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Header("ETag", hex.EncodeToString(ev.ETag[:]))
	c.Status(http.StatusOK)
}

func removeEvent(c *gin.Context) {
	_id := c.Param("id")
	id, err := strconv.ParseUint(_id, 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	var ind int = -1
	for i := range events {
		if events[i].ID == id {
			ind = i
			break
		}
	}

	if ind == -1 {
		c.Status(http.StatusNotFound)
		return
	}

	events = append(events[:ind], events[ind+1:]...)
	c.Status(http.StatusOK)
}

func compareETag(e1, e2 []byte) bool {
	if len(e1) != 20 || len(e2) != 20 {
		return false
	}

	for i := range e1 {
		if e1[i] != e2[i] {
			return false
		}
	}

	return true
}

func updateEvent(c *gin.Context) {
	_id := c.Param("id")
	id, err := strconv.ParseUint(_id, 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	var event *Event = nil
	var ind int
	for ind = range events {
		if events[ind].ID == id {
			event = &events[ind]
		}
	}

	if event == nil {
		c.Status(http.StatusNotFound)
		return
	}

	et := c.GetHeader("If-Match")
	rqet, err := hex.DecodeString(et)

	if err != nil || !compareETag(rqet, event.ETag[:]) {
		c.Status(http.StatusConflict)
		return
	}

	var ev Event
	if err := c.BindJSON(&ev); err != nil {
		c.Status(http.StatusUnprocessableEntity)
		return
	}

	if ev.ID != event.ID {
		c.Status(http.StatusBadRequest)
		return
	}

	events[ind] = ev
	_, err = MarshalAndTag(&events[ind])
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func getEvent(c *gin.Context) {
	_id := c.Param("id")
	id, err := strconv.ParseUint(_id, 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	var event *Event = nil
	for i := range events {
		if events[i].ID == id {
			event = &events[i]
		}
	}

	if event == nil {
		c.Status(http.StatusNotFound)
		return
	}

	_, err = MarshalAndTag(event)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Header("ETag", hex.EncodeToString(event.ETag[:]))
	c.JSON(http.StatusOK, &event)
}

func RunServer() {
	router := gin.Default()
	router.GET("/events", getEvents)
	router.POST("/events", addEvent)
	router.GET("/events/:id", getEvent)
	router.PUT("/events/:id", updateEvent)
	router.DELETE("/events/:id", removeEvent)
	router.Run("localhost:8080")
}
