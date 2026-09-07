package utils

import (
	"mysql/constant/share"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetParamID(c *gin.Context) (int, bool) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return 0, false
	}
	if id <= 0 {
		share.ResponseError(c, http.StatusBadRequest, "id must be greater than 0")
		return 0, false
	}
	return id, true
}

func GetParamUUID(c *gin.Context) (string, bool) {
	id := c.Param("uuid")
	if id == "" {
		share.ResponseError(c, http.StatusBadRequest, "uuid is required")
		return "", false
	}
	return id, true
}
