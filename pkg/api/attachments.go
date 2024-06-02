package api

import (
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/GPlaczek/taskmaster/pkg/data"
	"github.com/gin-gonic/gin"
)

func (a *Api) getAttachments(c *gin.Context) {
	c.JSON(http.StatusOK, a.data.GetAttachments())
}

func (a *Api) getAttachment(c *gin.Context) {
	id, ok := getID(c)
	if !ok {
		return
	}

	at := a.data.GetAttachment(id)
	if at == nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.Header("ETag", hex.EncodeToString(at.ETag))
	c.JSON(http.StatusOK, &at)
}

func (a *Api) addAttachment(c *gin.Context) {
	_a, err := a.data.AddAttachment()
	if err != nil {
		c.Status(data.ErrToHttpStatus(err))
		return
	}

	c.Header("ETag", hex.EncodeToString(_a.ETag))
	c.Header("Location", fmt.Sprintf("/attachments/%d", *_a.ID))
	c.Status(http.StatusCreated)
}

func (a *Api) updateAttachment(c *gin.Context) {
	id, ok := getID(c)
	if !ok {
		return
	}

	rqet := checkETag(c)
	if rqet == nil {
		return
	}

	var at data.AttachmentData
	if err := c.ShouldBindJSON(&at); err != nil {
		c.Status(http.StatusUnprocessableEntity)
		return
	}

	e, err := a.data.UpdateAttachment(id, &at, rqet)
	if err != nil {
		c.Status(data.ErrToHttpStatus(err))
		return
	}

	c.Header("ETag", hex.EncodeToString(e.ETag))
	c.JSON(http.StatusOK, e)
}

func (a *Api) deleteAttachment(c *gin.Context) {
	id, ok := getID(c)
	if !ok {
		return
	}

	rqet := checkETag(c)
	if rqet == nil {
		return
	}

	
	if err := a.data.DeleteAttachment(id, rqet); err != nil {
		c.Status(data.ErrToHttpStatus(err))
		return
	}

	c.Status(http.StatusNoContent)
}

func (a *Api) attachmentRoutes(router *gin.Engine) {
	router.GET("/attachments", a.getAttachments)
	router.POST("/attachments", a.addAttachment)
	router.GET("/attachments/:id", a.getAttachment)
	router.PUT("/attachments/:id", a.updateAttachment)
	router.DELETE("/attachments/:id", a.deleteAttachment)
}
