package ports

import "authentication-service/internal/core/entities"

type UserRepository interface {
	AddUser(userData entities.User) (res string)
}
