package gorm

import (
	models "authentication-service/external/repository/adaptors/postgres/model"

	"gorm.io/gorm"
)

func SyncDB(DB *gorm.DB) {
	DB.AutoMigrate(&models.User{})
}
