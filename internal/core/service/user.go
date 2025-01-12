package service

import (
	ports "authentication-service/external/_ports"
	"authentication-service/internal/core/entities"
	"authentication-service/internal/core/logic"
	"crypto/sha256"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type userService struct {
	userRepo ports.UserRepository
}

func NewUserService(userRepo ports.UserRepository) logic.UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) CreateUser(userData entities.User) (string, error) {
	secret := os.Getenv("UUID_SECRET")
	uuid, err := generateUUID(userData.Email, secret)
	if err != nil {
		return "", err
	}

	user := entities.User{
		ID:      uuid.String(),
		Email:   userData.Email,
		Name:    userData.Name,
		Picture: userData.Picture,
	}
	s.userRepo.AddUser(user)

	return uuid.String(), nil
}

func (s *userService) CheckUser(email string) (bool, string, error) {
	isFound, user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return false, "", err
	} else {
		return isFound, user.ID, nil
	}
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

func generateUUID(email string, secret string) (uuid.UUID, error) {
	data := email + secret

	// Hash the combined string using SHA-256
	hash := sha256.New()
	hash.Write([]byte(data))
	hashBytes := hash.Sum(nil)

	// Generate a UUID based on the hash
	uuidFromHash, err := uuid.FromBytes(hashBytes[:16]) // Using the first 16 bytes of the hash
	if err != nil {
		return uuid.Nil, err
	}

	return uuidFromHash, nil
}
