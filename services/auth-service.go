package services

import (
	"github.com/mashingan/smapping"
	"golang-api-jwt/dto"
	"golang-api-jwt/entity"
	"golang-api-jwt/repository"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type AuthService interface {
	VerifyCredentials(email string, password string) interface{}
	CreateUser(user dto.RegisterDTO) entity.User
	FindByEmail(email string) entity.User
	IsDuplicateEmail(email string) bool
}

type authService struct {
	userRepository repository.UserRepository
}

func NewAuthService(userRepository repository.UserRepository) AuthService {
	return &authService{
		userRepository: userRepository,
	}
}

func (service *authService) VerifyCredentials(email string, password string) interface{} {
	result := service.userRepository.VerifyCredentials(email, password)
	if value, ok := result.(entity.User); ok {
		isPasswordMatched := comparePassword(value.Password, []byte(password))
		if value.Email == email && isPasswordMatched {
			return result
		}
		return false
	}
	return false
}

func (service *authService) CreateUser(user dto.RegisterDTO) entity.User {
	userToCreate := entity.User{}
	err := smapping.FillStruct(&userToCreate, smapping.MapFields(&user))
	if err != nil {
		log.Println(err)
	}
	result := service.userRepository.CreateUser(userToCreate)
	return result

}

func (service *authService) FindByEmail(email string) entity.User {
	user := service.FindByEmail(email)
	return user
}

func (service *authService) IsDuplicateEmail(email string) bool {
	result := service.userRepository.IsDuplicateEmail(email)
	return !(result.Error == nil)
}

func comparePassword(hashedPassword string, plainPassword []byte) bool {
	byteHash := []byte(hashedPassword)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPassword)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
