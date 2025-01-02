package service

import (
	"os"
	ports "project-management-service/external/_ports"
	"project-management-service/internal/core/entities"
	"project-management-service/internal/core/logic"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type userService struct {
	userRepo ports.UserRepository
}

func NewUserService(userRepo ports.UserRepository) logic.UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) CreateUser(userData entities.User) error {
	s.userRepo.AddUser(userData)
	return nil
}

func (s *userService) GenerateJWT(id string, email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"sub": email,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})

	secret := os.Getenv("SECRET")
	tokenString, err := token.SignedString([]byte(secret))
	return tokenString, err
}
