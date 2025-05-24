package infrastructure

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

type PasswordService interface {
	PasswordHasher(password string) (string, error)
	PasswordComparator(hash, password string) bool
}

type passwordService struct{}

func NewPasswordService() PasswordService {
	return &passwordService{}
}

func (s *passwordService) PasswordHasher(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (s *passwordService) PasswordComparator(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	log.Printf("Password comparison error: %v", err)

	return err != nil
}
