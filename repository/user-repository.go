package repository

import (
	"golang-api-jwt/entity"
	"golang-api-jwt/helper"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user entity.User) entity.User
	UpdateUser(user entity.User) entity.User
	VerifyCredentials(email string, password string) interface{}
	IsDuplicateEmail(email string) (tx *gorm.DB)
	FindByEmail(email string) entity.User
	ProfileUser(userID string) entity.User
}

type userConnection struct {
	connection *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userConnection {
	return &userConnection{
		connection: db,
	}
}

func (db *userConnection) CreateUser(user entity.User) entity.User {
	user.Password, _ = helper.HashPassword(user.Password)
	db.connection.Save(&user)
	return user
}

func (db *userConnection) UpdateUser(user entity.User) entity.User {
	if user.Password != "" {
		user.Password, _ = helper.HashPassword(user.Password)
	} else {
		var temporaryUser entity.User
		db.connection.Find(&temporaryUser, user.ID)
		user.Password = temporaryUser.Password
	}
	db.connection.Save(&user)
	return user
}

func (db *userConnection) VerifyCredentials(email string, password string) interface{} {
	var user entity.User
	result := db.connection.Where("email = ?", email).Take(&user)
	if result.Error == nil {
		return user
	}
	return nil
}

func (db *userConnection) IsDuplicateEmail(email string) (tx *gorm.DB) {
	var user entity.User
	return db.connection.Where("email = ?", user.Email).Take(&user)
}

func (db *userConnection) FindByEmail(email string) entity.User {
	var user entity.User
	db.connection.Where("email = ?", user.Email).Take(&user)
	return user
}

func (db *userConnection) ProfileUser(userID string) entity.User {
	var user entity.User
	db.connection.Preload("Books").Preload("Books.User").Find(&user, userID)
	return user
}
