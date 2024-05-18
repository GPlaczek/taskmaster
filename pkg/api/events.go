package api

import (
	"encoding/hex"
	"errors"
	"net/http"
	"strconv"

	"github.com/GPlaczek/taskmaster/pkg/data"
	"github.com/gin-gonic/gin"
)

func (a *Api) getEvents(c *gin.Context) {
	c.JSON(http.StatusOK, a.data.GetEvents())
}

func (a *Api) getEvent(c *gin.Context) {
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

	c.Header("ETag", hex.EncodeToString(ev.ETag))
	c.JSON(http.StatusOK, &ev)
}

func (a *Api) addEvent(c *gin.Context) {
	var ev data.EventData

	if err := c.ShouldBindJSON(&ev); err != nil {
		c.Status(http.StatusUnprocessableEntity)
		return
	}

	e, err := a.data.AddEvent(&ev)
	if err != nil {
		if errors.Is(err, data.ErrMissingField) || errors.Is(err, data.ErrInvalidId) {
			c.Status(http.StatusBadRequest)
			return
		} else {
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	c.Header("ETag", hex.EncodeToString(e.ETag))
	c.JSON(http.StatusOK, gin.H{"id": e.ID})
}

func (a *Api) updateEvent(c *gin.Context) {
	_id := c.Param("id")
	id, err := strconv.ParseInt(_id, 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	et := c.GetHeader("If-Match")
	rqet, err := hex.DecodeString(et)
	if err != nil {
		c.Status(http.StatusConflict)
		return
	}

	c.Status(http.StatusUnprocessableEntity)
	var ev data.EventData
	if err := c.ShouldBindJSON(&ev); err != nil {
		c.Status(http.StatusUnprocessableEntity)
		return
	}

	e, err := a.data.UpdateEvent(id, &ev, rqet)
	if err != nil {
		c.Status(data.ErrToHttpStatus(err))
		return
	}

	c.Header("ETag", hex.EncodeToString(e.ETag))
	c.Status(http.StatusOK)
}

func (a *Api) deleteEvent(c *gin.Context) {
	_id := c.Param("id")
	id, err := strconv.ParseInt(_id, 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	et := c.GetHeader("If-Match")
	rqet, err := hex.DecodeString(et)
	if err != nil {
		c.Status(http.StatusConflict)
		return
	}

	err = a.data.DeleteEvent(id, rqet)
	if err != nil {
		c.Status(data.ErrToHttpStatus(err))
		return
	}

	c.Status(http.StatusOK)
}

func (a *Api) eventRoutes(router *gin.Engine) {
	router.GET("/events", a.getEvents)
	router.POST("/events", a.addEvent)
	router.GET("/events/:id", a.getEvent)
	router.PUT("/events/:id", a.updateEvent)
	router.DELETE("/events/:id", a.deleteEvent)
}
