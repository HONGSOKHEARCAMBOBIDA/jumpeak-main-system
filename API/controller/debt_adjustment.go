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

type DebtAdjustmentController struct {
	service service.DebtAdjustmentService
}

func NewDebtAdjustmentController() DebtAdjustmentController {
	return DebtAdjustmentController{
		service: service.NewDebtAdjustmentService(),
	}
}

func (cr *DebtAdjustmentController) Get(c *gin.Context) {
	page, pageSize := helper.GetPagination(c)
	userID, ok := helper.GetUserID(c)
	if !ok {
		return
	}
	filter := map[string]string{
		"customer_id": c.Query("customer_id"),
		"type":        c.Query("type"),
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

func (cr *DebtAdjustmentController) Create(c *gin.Context) {
	var input request.DebtAdjustmentRequestCreate
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
