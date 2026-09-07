package controller

import (
	"log"
	"mysql/constant/share"
	"mysql/helper"
	"mysql/request"
	"mysql/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	service service.AuthService
}

func NewAuthController() AuthController {
	return AuthController{
		service: service.NewAuthService(),
	}
}

func (cr *AuthController) Login(c *gin.Context) {
	var input request.AuthRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		share.ResponseError(c, 400, err.Error())
		return
	}

	result, err := cr.service.Login(input, c)
	if err != nil {
		log.Printf("Login error: %v", err)
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.RespondDate(c, http.StatusOK, result)
}

func (cr *AuthController) Refresh(c *gin.Context) {

	// var input request.RefreshTokenRequest
	///	log.Printf(input.RefreshToken)
	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		share.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := cr.service.RefreshToken(cookie, c)
	if err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.RespondDate(c, http.StatusOK, result)
}

func (cr *AuthController) GetUserData(c *gin.Context) {
	userlog, ok := helper.GetUserID(c)
	if !ok {
		share.ResponseError(c, http.StatusUnauthorized, "invalid user context")
		return
	}
	data, err := cr.service.GetUserData(c, userlog)
	if err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}

func (cr *AuthController) GetRole(c *gin.Context) {
	userID, ok := helper.GetUserID(c)
	if !ok {
		share.ResponseError(c, http.StatusUnauthorized, "please login")
		return
	}
	data, err := cr.service.GetRole(c, userID)
	if err != nil {
		share.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	share.RespondDate(c, http.StatusOK, data)
}
