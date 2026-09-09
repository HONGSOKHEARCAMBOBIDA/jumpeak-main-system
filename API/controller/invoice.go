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

type InvoiceController struct {
	service service.InvoiceService
}

func NewInvoiceController() InvoiceController {
	return InvoiceController{
		service: service.NewInvoiceService(),
	}
}

func (cr *InvoiceController) Get(c *gin.Context) {
	page, pageSize := helper.GetPagination(c)
	userID, ok := helper.GetUserID(c)
	if !ok {
		return
	}
	filter := map[string]string{
		"customer_id":    c.Query("customer_id"),
		"status":         c.Query("status"),
		"invoice_number": c.Query("invoice_number"),
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

func (cr *InvoiceController) Create(c *gin.Context) {
	var input request.InvoiceRequestCreate
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

func (cr *InvoiceController) Cancel(c *gin.Context) {
	var input request.InvoiceRequestCancel
	id, ok := utils.GetParamID(c)
	if !ok {
		return
	}
	userID, ok := helper.GetUserID(c)
	if !ok {
		return
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.Cancel(c.Request.Context(), id, userID, input); err != nil {
		share.RespondServiceError(c, err)
		return
	}
	share.ResponseSuccess(c, http.StatusOK, share.Created)
}
