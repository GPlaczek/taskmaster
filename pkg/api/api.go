package api

import (
	"strconv"
	"errors"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/GPlaczek/taskmaster/pkg/data"
)

type Api struct {
	data data.Data
}

func NewApi(d data.Data) *Api {
	return &Api{
		data: d,
	}
}

func (a *Api)getEvents(c *gin.Context) {
	c.JSON(http.StatusOK, a.data.GetEvents())
}

func (a *Api)getEvent(c *gin.Context) {
	_id := c.Param("id")
	id, err := strconv.ParseInt(_id, 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	ev := a.data.GetEvent(id)

	if ev == nil {
		c.Status(http.StatusNotFound)
		return
	}

	ev.ETagUpdate()
	c.Header("ETag", hex.EncodeToString(ev.ETagGet()))
	c.JSON(http.StatusOK, &ev)
}

func (a *Api)addEvent(c *gin.Context) {
	var ev data.EventData

	if err := c.BindJSON(&ev); err != nil {
		c.Status(http.StatusUnprocessableEntity)
		return
	}

	e, err := a.data.AddEvent(&ev)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Header("ETag", hex.EncodeToString(e.ETagGet()))
	c.Status(http.StatusOK)
}

func (a *Api)updateEvent(c *gin.Context) {
	_id := c.Param("id")
	id, err := strconv.ParseInt(_id, 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	event := a.data.GetEvent(id) 
	if event == nil {
		c.Status(http.StatusNotFound)
		return
	}

	et := c.GetHeader("If-Match")
	rqet, err := hex.DecodeString(et)

	if err != nil || !event.ETagCompare(rqet) {
		c.Status(http.StatusConflict)
		return
	}

	var ev data.EventData
	if err := c.BindJSON(&ev); err != nil {
		c.Status(http.StatusUnprocessableEntity)
		return
	}

	err = event.Update(&ev)
	if err != nil {
		println(err)
		c.Status(http.StatusBadRequest)
		return
	}

	c.Status(http.StatusOK)
}

func (a *Api)deleteEvent(c *gin.Context) {
	_id := c.Param("id")
	id, err := strconv.ParseInt(_id, 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	err = a.data.DeleteEvent(id)
	if err == nil {
		c.Status(http.StatusOK)
		return
	}

	if errors.Is(err, data.ErrNotFound) {
		c.Status(http.StatusNotFound)
	} else {
		c.Status(http.StatusInternalServerError)
	}
}

func (a *Api)RunServer() {
	router := gin.Default()
	router.GET("/events", a.getEvents)
	router.POST("/events", a.addEvent)
	router.GET("/events/:id", a.getEvent)
	router.PUT("/events/:id", a.updateEvent)
	router.DELETE("/events/:id", a.deleteEvent)
	router.Run("localhost:8080")
}
