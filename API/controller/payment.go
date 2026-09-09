package controller

import (
	"context"
	"errors"
	"mysql/constant/share"
	"mysql/helper"
	"mysql/request"
	"mysql/service"
	"mysql/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PaymentController struct {
	service service.PaymentService
}

func NewPaymentController() PaymentController {
	return PaymentController{
		service: service.NewPaymentService(),
	}
}

func (cr *PaymentController) Create(c *gin.Context) {
	var input request.PaymentRequestCreate
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

func (cr *PaymentController) Void(c *gin.Context) {
	var input request.PaymentRequestVoid
	userID, ok := helper.GetUserID(c)
	if !ok {
		return
	}
	id, ok := utils.GetParamID(c)
	if !ok {
		return
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.Void(c.Request.Context(), id, userID, input); err != nil {
		share.RespondServiceError(c, err)
		return
	}
	share.ResponseSuccess(c, http.StatusOK, share.Created)
}

func (cr *PaymentController) Get(c *gin.Context) {
	page, pageSize := helper.GetPagination(c)
	userID, ok := helper.GetUserID(c)
	if !ok {
		return
	}
	filter := map[string]string{
		"customer_id": c.Query("customer_id"),
		"method":      c.Query("method"),
		"status":      c.Query("status"),
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
