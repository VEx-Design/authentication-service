package logic

import (
	"authentication-service/internal/core/entities"
)

type UserService interface {
	GenerateJWT(ID string, Email string) (string, error)
	CreateUser(userData entities.User) (string, error)
	CheckUser(email string) (bool, string, error)
}
