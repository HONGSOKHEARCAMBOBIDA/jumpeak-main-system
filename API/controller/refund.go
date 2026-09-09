package controller

import (
	"context"
	"errors"
	"mysql/constant/share"
	"mysql/helper"
	"mysql/request"
	"mysql/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RefundController struct {
	service service.RefundService
}

func NewRefundController() RefundController {
	return RefundController{
		service: service.NewRefundService(),
	}
}

func (cr *RefundController) Create(c *gin.Context) {
	var input request.RefundRequestCreate
	userID, ok := helper.GetUserID(c)
	if !ok {
		return
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.Create(c.Request.Context(), userID, input); err != nil {
		share.RespondServiceError(c, err)
		return
	}
	share.ResponseSuccess(c, http.StatusOK, share.Created)
}

func (cr *RefundController) Get(c *gin.Context) {
	page, pageSize := helper.GetPagination(c)
	userID, ok := helper.GetUserID(c)
	if !ok {
		return
	}
	filter := map[string]string{
		"payment_id": c.Query("name"),
	}

	data, meta, err := cr.service.Get(c.Request.Context(), userID, request.Pagination{
		Page:     page,
		PageSize: pageSize,
	}, filter)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			share.ResponseError(c, http.StatusGatewayTimeout, err.Error())
			return
		}
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.ResponsePagination(c, 200, data, meta)
}
