package controllers

import (
	"github.com/gin-gonic/gin"
	"golang-api-jwt/dto"
	"golang-api-jwt/entity"
	"golang-api-jwt/helper"
	"golang-api-jwt/services"
	"net/http"
	"strconv"
)

type AuthController interface {
	Login(ctx *gin.Context)
	Register(ctx *gin.Context)
}

type authController struct {
	authService services.AuthService
	jwtService  services.JWTService
}

func NewAuthController(authService services.AuthService, jwtService services.JWTService) AuthController {
	return &authController{authService: authService, jwtService: jwtService}
}

func (c *authController) Login(ctx *gin.Context) {
	var loginDTO dto.LoginDTO
	err := ctx.ShouldBind(&loginDTO)
	if err != nil {
		errorResponse := helper.SendErrorResponse(false, "Failed to process request", err.Error(), helper.EmptyObject{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse)
		return
	}
	authResult := c.authService.VerifyCredentials(loginDTO.Email, loginDTO.Password)
	if value, ok := authResult.(entity.User); ok {
		generatedToken := c.jwtService.GenerateToken(strconv.FormatUint(value.ID, 10))
		value.Token = generatedToken
		successResponse := helper.SendSuccessResponse(true, "OK!", value)
		ctx.JSON(http.StatusOK, successResponse)
		return
	}
	errorResponse := helper.SendErrorResponse(false, "Please check your email and password", "Invalid Credential", helper.EmptyObject{})
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse)
}

func (c *authController) Register(ctx *gin.Context) {
	var registerDTO dto.RegisterDTO
	errorDTO := ctx.ShouldBind(&registerDTO)
	if errorDTO != nil {
		errorResponse := helper.SendErrorResponse(false, "Failed to process request", errorDTO.Error(), helper.EmptyObject{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse)
		return
	}
	if !c.authService.IsDuplicateEmail(registerDTO.Email) {
		errorResponse := helper.SendErrorResponse(
			false, "Failed to process request", "Duplicate email", helper.EmptyObject{},
		)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse)
	} else {
		createdUser := c.authService.CreateUser(registerDTO)
		token := c.jwtService.GenerateToken(strconv.FormatUint(createdUser.ID, 10))
		createdUser.Token = token
		response := helper.SendSuccessResponse(true, "OK!", createdUser)
		ctx.JSON(http.StatusOK, response)
	}

}
