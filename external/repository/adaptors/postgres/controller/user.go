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

func (r *userRepositoryPQ) AddUser(userData entities.User) error {
	var user models.User
	if err := r.client.FirstOrCreate(&user, models.User{
		ID:      userData.ID,
		Email:   userData.Email,
		Name:    userData.Name,
		Picture: userData.Picture,
	}).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepositoryPQ) GetUserByEmail(email string) (isFound bool, res entities.User, err error) {
	var user models.User
	if err := r.client.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, entities.User{}, nil
		} else {
			return false, entities.User{}, err
		}
	}
	return true, entities.User{
		ID:      user.ID,
		Email:   user.Email,
		Name:    user.Name,
		Picture: user.Picture,
	}, nil
}
