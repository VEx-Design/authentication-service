package ports

import "authentication-service/internal/core/entities"

type UserRepository interface {
	AddUser(userData entities.User) (err error)
	GetUserByEmail(email string) (isFound bool, res entities.User, err error)
}
