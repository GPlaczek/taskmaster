package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/GPlaczek/taskmaster/pkg/data"
	"github.com/gin-gonic/gin"
)

func (a *Api) addMerge(c *gin.Context) {
	var md data.MergeData

	if err := c.ShouldBindJSON(&md); err != nil {
		c.Status(http.StatusUnprocessableEntity)
		return
	}

	ed, md1, err := a.data.MergeEvents(&md)
	if err != nil {
		c.Status(data.ErrToHttpStatus(err))
		return
	}

	c.Header("Content-Location", fmt.Sprintf("/events/%d", ed.ID))
	c.JSON(http.StatusOK, md1)
}

func (a *Api) getMerges(c *gin.Context) {
	c.JSON(http.StatusOK, a.data.GetMerges())
}

func (a *Api) getMerge(c *gin.Context) {
	_id := c.Param("id")
	id, err := strconv.ParseInt(_id, 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	ev := a.data.GetMerge(id)

	if ev == nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, &ev)
}

func (a *Api) mergeRoutes(router *gin.Engine) {
	router.POST("/merges", a.addMerge)
	router.GET("/merges", a.getMerges)
	router.GET("/merges/:id", a.getMerge)
}
