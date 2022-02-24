package helper

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

func HashPassword(pass string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), 12)
	if err != nil {
		log.Fatal(err)
	}
	return string(bytes), err
}
