package share

import (
	"mysql/constant/apperror"

	"github.com/gin-gonic/gin"
)

func RespondServiceError(c *gin.Context, err error) {
	status := apperror.HTTPStatus(err)
	msg := apperror.ClientMessage(err)
	ResponseError(c, status, msg)
}
