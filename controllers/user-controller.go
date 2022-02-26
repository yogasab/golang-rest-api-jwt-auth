package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang-api-jwt/dto"
	"golang-api-jwt/helper"
	"golang-api-jwt/services"
	"net/http"
	"strconv"
)

type UserController interface {
	Update(context *gin.Context)
	MyProfile(context *gin.Context)
}

type userController struct {
	userService services.UserService
	jwtService  services.JWTService
}

func NewUserController(userService services.UserService, jwtService services.JWTService) UserController {
	return &userController{
		userService: userService,
		jwtService:  jwtService,
	}
}

func (c *userController) Update(context *gin.Context) {
	var userUpdateDTO dto.UserUpdateDTO
	errDTO := context.ShouldBind(&userUpdateDTO)
	if errDTO != nil {
		errorResponse := helper.SendErrorResponse(false, "Failed to process request", errDTO.Error(), helper.EmptyObject{})
		context.AbortWithStatusJSON(http.StatusBadRequest, errorResponse)
		return
	}
	authHeader := context.GetHeader("Authorization")
	token, errToken := c.jwtService.VerifyToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}
	claims := token.Claims.(jwt.MapClaims)
	id, err := strconv.ParseUint(fmt.Sprintf("%v", claims["user_id"]), 10, 64)
	if err != nil {
		panic(err.Error())
	}
	userUpdateDTO.ID = id
	updatedUser := c.userService.UpdateProfile(userUpdateDTO)
	successResponse := helper.SendSuccessResponse(true, "OK!", updatedUser)
	context.JSON(http.StatusOK, successResponse)
}

func (c *userController) MyProfile(context *gin.Context) {
	authHeader := context.GetHeader("Authorization")
	token, err := c.jwtService.VerifyToken(authHeader)
	if err != nil {
		panic(err.Error())
	}
	claims := token.Claims.(jwt.MapClaims)
	userID := fmt.Sprintf("%v", claims["user_id"])
	userProfile := c.userService.MyProfile(userID)
	successResponse := helper.SendSuccessResponse(true, "OK!", userProfile)
	context.JSON(http.StatusOK, successResponse)
}
