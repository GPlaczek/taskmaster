package api

import (
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

func (a *Api) RunServer() {
	router := gin.Default()
	a.eventRoutes(router)
	a.attachmentRoutes(router)
	router.Run("localhost:8080")
}
