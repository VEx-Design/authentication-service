package repository

import (
	ports "authentication-service/external/_ports"
	models "authentication-service/external/repository/adaptors/postgres/model"
	"authentication-service/internal/core/entities"

	"gorm.io/gorm"
)

type userRepositoryPQ struct {
	client *gorm.DB
}

func NewUserRepositoryPQ(client *gorm.DB) ports.UserRepository {
	return &userRepositoryPQ{client: client}
}

func (r *userRepositoryPQ) AddUser(userData entities.User) (res string) {
	var user models.User
	if err := r.client.FirstOrCreate(&user, models.User{
		ID:      userData.ID,
		Email:   userData.Email,
		Name:    userData.Name,
		Picture: userData.Picture,
	}).Error; err != nil {
		return "have error"
	}
	return "success"
}
