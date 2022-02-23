package main

import (
	"github.com/gin-gonic/gin"
	"golang-api-jwt/controllers"
)

var (
	authController controllers.AuthController = controllers.NewAuthController()
)

func main() {
	router := gin.Default()

	v1 := router.Group("/api/auth")
	{
		v1.POST("/login", authController.Login)
		v1.POST("/register", authController.Register)
	}
	router.Run(":5000")
}
