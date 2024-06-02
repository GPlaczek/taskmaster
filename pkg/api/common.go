package api;

import (
	"net/http"
	"encoding/hex"
	"strconv"

	"github.com/gin-gonic/gin"
)

func checkETag(c *gin.Context) []byte {
	et := c.GetHeader("If-Match")
	if et == "" {
		c.Status(http.StatusPreconditionRequired)
		return nil
	}

	rqet, err := hex.DecodeString(et)
	if err != nil {
		c.Status(http.StatusPreconditionFailed)
		return nil
	}

	return rqet
}

func getID(c *gin.Context) (int64, bool) {
	_id := c.Param("id")
	id, err := strconv.ParseInt(_id, 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return -1, false
	}

	return id, true
}
