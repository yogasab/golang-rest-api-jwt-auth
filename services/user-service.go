package services

import (
	"github.com/mashingan/smapping"
	"golang-api-jwt/dto"
	"golang-api-jwt/entity"
	"golang-api-jwt/repository"
	"log"
)

type UserService interface {
	UpdateProfile(user dto.UserUpdateDTO) entity.User
	MyProfile(userID string) entity.User
}

type userService struct {
	userRepository repository.UserRepository
}

func NewSUserService(userRepository repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}

func (s *userService) UpdateProfile(user dto.UserUpdateDTO) entity.User {
	userToUpdate := entity.User{}
	err := smapping.FillStruct(&userToUpdate, smapping.MapFields(&user))
	if err != nil {
		log.Fatalln(err)
	}
	updateUser := s.userRepository.UpdateUser(userToUpdate)
	return updateUser
}

func (s *userService) MyProfile(userID string) entity.User {
	return s.userRepository.ProfileUser(userID)
}
