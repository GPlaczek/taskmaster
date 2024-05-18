package api

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/GPlaczek/taskmaster/pkg/data"
	"github.com/gin-gonic/gin"
)

func (a *Api) getAttachments(c *gin.Context) {
	c.JSON(http.StatusOK, a.data.GetAttachments())
}

func (a *Api) getAttachment(c *gin.Context) {
	_id := c.Param("id")
	id, err := strconv.ParseInt(_id, 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	at := a.data.GetAttachment(id)
	fmt.Println(at)

	if at == nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.Header("ETag", hex.EncodeToString(at.ETag))
	c.JSON(http.StatusOK, &at)
}

func (a *Api) addAttachment(c *gin.Context) {
	var at data.AttachmentData

	if err := c.ShouldBindJSON(&at); err != nil {
		c.Status(http.StatusUnprocessableEntity)
		return
	}

	_a, err := a.data.AddAttachment(&at)
	if err != nil {
		if errors.Is(err, data.ErrMissingField) || errors.Is(err, data.ErrInvalidId) {
			c.Status(http.StatusBadRequest)
			return
		} else {
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	c.Header("ETag", hex.EncodeToString(_a.ETag))
	c.JSON(http.StatusOK, gin.H{"id": _a.ID})
}

func (a *Api) updateAttachment(c *gin.Context) {
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
	c.Status(http.StatusOK)
}

func (a *Api) deleteAttachment(c *gin.Context) {
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

	err = a.data.DeleteAttachment(id, rqet)
	if err != nil {
		c.Status(data.ErrToHttpStatus(err))
		return
	}

	c.Status(http.StatusOK)
}

func (a *Api) attachmentRoutes(router *gin.Engine) {
	router.GET("/attachments", a.getAttachments)
	router.POST("/attachments", a.addAttachment)
	router.GET("/attachments/:id", a.getAttachment)
	router.PUT("/attachments/:id", a.updateAttachment)
	router.DELETE("/attachments/:id", a.deleteAttachment)
}
