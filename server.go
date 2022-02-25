package main

import (
	"github.com/gin-gonic/gin"
	"golang-api-jwt/config"
	"golang-api-jwt/controllers"
	"golang-api-jwt/repository"
	"golang-api-jwt/services"
	"gorm.io/gorm"
)

var (
	db             *gorm.DB                   = config.SetupDatabaseConnection()
	userRepository repository.UserRepository  = repository.NewUserRepository(db)
	jwtService     services.JWTService        = services.NewJWTService()
	authService    services.AuthService       = services.NewAuthService(userRepository)
	authController controllers.AuthController = controllers.NewAuthController(authService, jwtService)
)

func main() {
	defer config.CloseDatabaseConnection(db)
	router := gin.Default()

	v1 := router.Group("/api/auth")
	{
		v1.POST("/login", authController.Login)
		v1.POST("/register", authController.Register)
	}
	router.Run(":5000")
}
