package ports

import "project-management-service/internal/core/entities"

type UserRepository interface {
	AddUser(userData entities.User) (res string)
}
