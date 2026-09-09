package controller

import (
	"mysql/constant/share"
	"mysql/helper"
	"mysql/request"
	"mysql/service"
	"mysql/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BranchController struct {
	service service.BranchService
}

func NewBranchController() BranchController {
	return BranchController{
		service: service.NewBranchService(),
	}
}

func (cr *BranchController) Create(c *gin.Context) {
	var input request.BranchRequestCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.Create(c.Request.Context(), input); err != nil {
		share.RespondServiceError(c, err)
		return
	}
	share.ResponseSuccess(c, http.StatusOK, share.Created)
}

func (cr *BranchController) Update(c *gin.Context) {
	id, ok := utils.GetParamID(c)
	if !ok {
		return
	}
	var input request.BranchRequestUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cr.service.Update(c.Request.Context(), id, input); err != nil {
		share.RespondServiceError(c, err)
		return
	}
	share.ResponseSuccess(c, http.StatusOK, share.Updated)
}

func (cr *BranchController) GetBranchNoPagination(c *gin.Context) {
	companyID, ok := utils.GetParamID(c)
	if !ok {
		return
	}
	userID, ok := helper.GetUserID(c)
	if !ok {
		return
	}
	data, err := cr.service.GetBranchNoPagination(c.Request.Context(), userID, companyID)
	if err != nil {
		share.ResponseError(c, http.StatusGatewayTimeout, err.Error())
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}
