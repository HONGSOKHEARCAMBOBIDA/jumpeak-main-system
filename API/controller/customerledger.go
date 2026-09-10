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

type CustomerledgerController struct {
	service service.CustomerLedgerService
}

func NewCustomerledgerController() CustomerledgerController {
	return CustomerledgerController{
		service: service.NewCustomerLedgerService(),
	}
}

func (cr *CustomerledgerController) Get(c *gin.Context) {
	page, pageSize := helper.GetPagination(c)
	userID, ok := helper.GetUserID(c)
	if !ok {
		return
	}
	filter := map[string]string{
		"customer_id":    c.Query("customer_id"),
		"reference_type": c.Query("reference_type"),
		"date_from":      c.Query("date_from"),
		"date_to":        c.Query("date_to"),
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
