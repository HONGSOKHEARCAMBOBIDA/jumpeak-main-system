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

type CustomerController struct {
	service service.CustomerService
}

func NewCustomerController() CustomerController {
	return CustomerController{
		service: service.NewCustomerService(),
	}
}

func (cr *CustomerController) Get(c *gin.Context) {
	page, pageSize := helper.GetPagination(c)
	userID, ok := helper.GetUserID(c)
	if !ok {
		return
	}
	filter := map[string]string{
		"name":       c.Query("name"),
		"company_id": c.Query("company_id"),
		"branch_id":  c.Query("branch_id"),
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

func (cr *CustomerController) Create(c *gin.Context) {
	var input request.CustomerRequestCreate
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

func (cr *CustomerController) Update(c *gin.Context) {
	id, ok := utils.GetParamID(c)
	if !ok {
		return
	}
	userID, ok := helper.GetUserID(c)
	if !ok {
		return
	}
	var input request.CustomerRequestUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.Update(c.Request.Context(), id, userID, input); err != nil {
		share.RespondServiceError(c, err)
		return
	}
	share.ResponseSuccess(c, http.StatusOK, share.Updated)
}
