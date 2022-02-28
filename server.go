package main

import (
	"github.com/gin-gonic/gin"
	"golang-api-jwt/config"
	"golang-api-jwt/controllers"
	"golang-api-jwt/middlewares"
	"golang-api-jwt/repository"
	"golang-api-jwt/services"
	"gorm.io/gorm"
)

var (
	db             *gorm.DB                   = config.SetupDatabaseConnection()
	userRepository repository.UserRepository  = repository.NewUserRepository(db)
	bookRepository repository.BookRepository  = repository.NewBookRepository(db)
	jwtService     services.JWTService        = services.NewJWTService()
	authService    services.AuthService       = services.NewAuthService(userRepository)
	userService    services.UserService       = services.NewSUserService(userRepository)
	bookService    services.BookService       = services.NewBookService(bookRepository)
	authController controllers.AuthController = controllers.NewAuthController(authService, jwtService)
	userController controllers.UserController = controllers.NewUserController(userService, jwtService)
	bookController controllers.BookController = controllers.NewBookController(bookService, jwtService)
)

func main() {
	defer config.CloseDatabaseConnection(db)
	router := gin.Default()

	authRoutes := router.Group("/api/auth")
	{
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/register", authController.Register)
	}
	userRoutes := router.Group("/api/v1/users", middlewares.AuthorizeJWT(jwtService))
	{
		userRoutes.GET("/profile", userController.MyProfile)
		userRoutes.PUT("/profile", userController.Update)
	}
	bookRoutes := router.Group("/api/v1/books", middlewares.AuthorizeJWT(jwtService))
	{
		bookRoutes.GET("/", bookController.All)
		bookRoutes.POST("/", bookController.Insert)
		bookRoutes.GET("/:id", bookController.FindByID)
		bookRoutes.PUT("/:id", bookController.Update)
		bookRoutes.DELETE("/:id", bookController.Delete)
	}
	router.Run(":5000")

}
