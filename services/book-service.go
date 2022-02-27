package services

import (
	"fmt"
	"github.com/mashingan/smapping"
	"golang-api-jwt/dto"
	"golang-api-jwt/entity"
	"golang-api-jwt/repository"
	"log"
)

type BookService interface {
	Insert(bookDTO dto.BookCreateDTO) entity.Book
	Update(bookDTO dto.BookUpdateDTO) entity.Book
	FindByID(bookID uint64) entity.Book
	Delete(book entity.Book)
	All() []entity.Book
	IsAllowedToEdit(userID string, bookID uint64) bool
}

type bookService struct {
	bookRepository repository.BookRepository
}

func NewBookService(bookRepository repository.BookRepository) BookService {
	return &bookService{bookRepository: bookRepository}
}

func (s *bookService) Insert(bookDTO dto.BookCreateDTO) entity.Book {
	// Instantiate the model to take
	book := entity.Book{}
	// Take each field from the model to map in smapping to BookDTO
	err := smapping.FillStruct(&book, smapping.MapFields(&bookDTO))
	if err != nil {
		log.Println("Error on BookService Insert ", err)
	}
	// Save
	result := s.bookRepository.InsertBook(book)
	return result
}

func (s *bookService) Update(bookDTO dto.BookUpdateDTO) entity.Book {
	book := entity.Book{}
	err := smapping.FillStruct(&book, smapping.MapFields(bookDTO))
	if err != nil {
		log.Println("Error on BookService Update ", err)
	}
	result := s.bookRepository.UpdateBook(book)
	return result
}

func (s *bookService) Delete(book entity.Book) {
	s.bookRepository.DeleteBook(book)
}

func (s *bookService) FindByID(bookID uint64) entity.Book {
	book := s.bookRepository.FindBookByID(bookID)
	return book
}

func (s *bookService) All() []entity.Book {
	return s.bookRepository.AllBooks()
}

func (s *bookService) IsAllowedToEdit(userID string, bookID uint64) bool {
	book := s.bookRepository.FindBookByID(bookID)
	id := fmt.Sprintf("User ID is ", book.UserID)
	return userID == id
}
