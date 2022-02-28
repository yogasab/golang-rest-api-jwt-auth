package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang-api-jwt/dto"
	"golang-api-jwt/entity"
	"golang-api-jwt/helper"
	"golang-api-jwt/services"
	"log"
	"net/http"
	"strconv"
)

type BookController interface {
	Insert(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	FindByID(c *gin.Context)
	All(c *gin.Context)
}

type bookController struct {
	bookService services.BookService
	JWTService  services.JWTService
}

func NewBookController(bookService services.BookService, JWTService services.JWTService) BookController {
	return &bookController{bookService: bookService, JWTService: JWTService}
}

func (c *bookController) All(ctx *gin.Context) {
	var books []entity.Book = c.bookService.All()
	successResponse := helper.SendSuccessResponse(true, "Books fetched successfully", books)
	ctx.JSON(http.StatusOK, successResponse)
}

func (c *bookController) FindByID(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 0, 0)
	if err != nil {
		errorResponse := helper.SendErrorResponse(false, "No param id was not found", err.Error(), helper.EmptyObject{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse)
		return
	}
	var book entity.Book = c.bookService.FindByID(uint64(id))
	if (book == entity.Book{}) {
		errorResponse := helper.SendErrorResponse(false, "Book not found", err.Error(), helper.EmptyObject{})
		ctx.JSON(http.StatusNotFound, errorResponse)
		return
	}
	successResponse := helper.SendSuccessResponse(true, "Boot fetched successfully", book)
	ctx.JSON(http.StatusOK, successResponse)
}

func (c *bookController) Insert(ctx *gin.Context) {
	var bookCreateDTO dto.BookCreateDTO
	errDTO := ctx.ShouldBind(&bookCreateDTO)
	if errDTO != nil {
		errorResponse := helper.SendErrorResponse(false, "Failed to process request", errDTO.Error(), helper.EmptyObject{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse)
		return
	} else {
		authHeader := ctx.GetHeader("Authorization")
		userID := c.getUserIDFromToken(authHeader)
		convertedUserID, err := strconv.ParseUint(userID, 10, 64)
		if err == nil {
			bookCreateDTO.UserID = convertedUserID
		}
		result := c.bookService.Insert(bookCreateDTO)
		successResponse := helper.SendSuccessResponse(true, "Book created successfully", result)
		ctx.JSON(http.StatusCreated, successResponse)
	}
}

func (c *bookController) Update(ctx *gin.Context) {
	var updatedBookDTO dto.BookUpdateDTO
	errDTO := ctx.ShouldBind(&updatedBookDTO)
	if errDTO != nil {
		errorResponse := helper.SendErrorResponse(false, "Failed to process request", errDTO.Error(), helper.EmptyObject{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse)
		return
	}
	authHeaderToken := ctx.GetHeader("Authorization")
	token, errToken := c.JWTService.VerifyToken(authHeaderToken)
	if errToken != nil {
		errorResponse := helper.SendErrorResponse(false, "Token is not verified", errToken.Error(), helper.EmptyObject{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse)
		return
	}
	claims := token.Claims.(jwt.MapClaims)
	userID := fmt.Sprintf("%v", claims["user_id"])
	if !c.bookService.IsAllowedToEdit(userID, updatedBookDTO.ID) {
		id, errID := strconv.ParseUint(userID, 10, 64)
		if errID == nil {
			updatedBookDTO.UserID = id
		}
		result := c.bookService.Update(updatedBookDTO)
		successResponse := helper.SendSuccessResponse(true, "Book updated successfully", result)
		ctx.JSON(http.StatusOK, successResponse)
	} else {
		response := helper.SendErrorResponse(false, "You dont have permission", "You are not the owner", helper.EmptyObject{})
		ctx.JSON(http.StatusForbidden, response)
	}
}

func (c *bookController) Delete(ctx *gin.Context) {
	var book entity.Book
	bookID, err := strconv.ParseUint(ctx.Param("id"), 0, 0)
	if err != nil {
		errorResponse := helper.SendErrorResponse(false, "No book id was found", err.Error(), helper.EmptyObject{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse)
		return
	}
	book.ID = bookID
	authHeaderToken := ctx.GetHeader("Authorization")
	token, errToken := c.JWTService.VerifyToken(authHeaderToken)
	if errToken != nil {
		errorResponse := helper.SendErrorResponse(false, "Token is not verified", errToken.Error(), helper.EmptyObject{})
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse)
		return
	}
	claims := token.Claims.(jwt.MapClaims)
	userID := fmt.Sprintf("%v", claims["user_id"])
	if !c.bookService.IsAllowedToEdit(userID, book.ID) {
		c.bookService.Delete(book)
		successResponse := helper.SendSuccessResponse(true, "Book deleted successfully", helper.EmptyObject{})
		ctx.JSON(http.StatusNoContent, successResponse)
	} else {
		errorResponse := helper.SendErrorResponse(false, "You dont have permission to deleted this", "You are not the owner", helper.EmptyObject{})
		ctx.JSON(http.StatusForbidden, errorResponse)
	}
}

func (c *bookController) getUserIDFromToken(token string) string {
	authToken, err := c.JWTService.VerifyToken(token)
	if err != nil {
		log.Println("Token is not verify ", err.Error())
	}
	claims := authToken.Claims.(jwt.MapClaims)
	id := fmt.Sprintf("User ID is ", claims["user_id"])
	return id
}
