package api

import (
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/GPlaczek/taskmaster/pkg/data"
	"github.com/gin-gonic/gin"
)

func (a *Api) getEvents(c *gin.Context) {
	c.JSON(http.StatusOK, a.data.GetEvents())
}

func (a *Api) getEvent(c *gin.Context) {
	id, ok := getID(c)
	if !ok {
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
	e, err := a.data.AddEvent()
	if err != nil {
		c.Status(data.ErrToHttpStatus(err))
		return
	}

	c.Header("ETag", hex.EncodeToString(e.ETag))
	c.Header("Location", fmt.Sprintf("/events/%d", *e.ID))
	c.Status(http.StatusCreated)
}

func (a *Api) updateEvent(c *gin.Context) {
	id, ok := getID(c)
	if !ok {
		return
	}

	rqet := checkETag(c)
	if rqet == nil {
		return
	}

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
	c.JSON(http.StatusOK, e)
}

func (a *Api) deleteEvent(c *gin.Context) {
	id, ok := getID(c)
	if !ok {
		return
	}

	rqet := checkETag(c)
	if rqet == nil {
		return
	}

	if err := a.data.DeleteEvent(id, rqet); err != nil {
		c.Status(data.ErrToHttpStatus(err))
		return
	}

	c.Status(http.StatusNoContent)
}

func (a *Api) bindAttachment(c *gin.Context) {
	eid, ok := getID(c)
	if !ok {
		return
	}

	var at data.AttachmentData
	if err := c.ShouldBindJSON(&at); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if at.ID == nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if err := a.data.BindAttachment(eid, *at.ID); err != nil {
		c.Status(data.ErrToHttpStatus(err))
		return
	}

	c.Status(http.StatusCreated)
}

func (a *Api) getBoundAttachments(c *gin.Context) {
	eid, ok := getID(c)
	if !ok {
		return
	}

	ad, err := a.data.GetBoundAttachments(eid)
	if err != nil {
		c.Status(data.ErrToHttpStatus(err))
		return
	}

	c.JSON(http.StatusOK, ad)
}

func (a *Api) unbindAttachment(c *gin.Context) {
	eid, ok := getID(c)
	if !ok {
		return
	}

	aid, ok := getID2(c, "aid")
	if !ok {
		return
	}

	err := a.data.UnbindAttachment(eid, aid)
	if err != nil {
		c.Status(data.ErrToHttpStatus(err))
		return
	}

	c.Status(http.StatusNoContent)
}

func (a *Api) eventRoutes(router *gin.Engine) {
	router.GET("/events", a.getEvents)
	router.POST("/events", a.addEvent)

	eventRouter := router.Group("/events/:id")
	eventRouter.GET("", a.getEvent)
	eventRouter.PUT("", a.updateEvent)
	eventRouter.DELETE("", a.deleteEvent)

	eventRouter.POST("/attachments", a.bindAttachment)
	eventRouter.GET("/attachments", a.getBoundAttachments)
	eventRouter.DELETE("/attachments/:aid", a.unbindAttachment) 
}
