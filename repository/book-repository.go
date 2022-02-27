package repository

import (
	"golang-api-jwt/entity"
	"gorm.io/gorm"
)

type BookRepository interface {
	InsertBook(book entity.Book) entity.Book
	UpdateBook(book entity.Book) entity.Book
	FindBookByID(bookID uint64) entity.Book
	AllBooks() []entity.Book
	DeleteBook(book entity.Book)
}

type bookConnection struct {
	connection *gorm.DB
}

func NewBookRepository(db *gorm.DB) BookRepository {
	return &bookConnection{connection: db}
}

func (db *bookConnection) InsertBook(book entity.Book) entity.Book {
	db.connection.Save(&book)
	db.connection.Preload("User").Find(&book)
	return book
}

func (db *bookConnection) UpdateBook(book entity.Book) entity.Book {
	db.connection.Save(&book)
	db.connection.Preload("User").Find(&book)
	return book
}

func (db *bookConnection) FindBookByID(bookID uint64) entity.Book {
	var book entity.Book
	db.connection.Preload("User").Find(&book, bookID)
	return book
}

func (db *bookConnection) AllBooks() []entity.Book {
	var books []entity.Book
	db.connection.Preload("User").Find(&books)
	return books
}

func (db *bookConnection) DeleteBook(book entity.Book) {
	db.connection.Delete(&book)
}
