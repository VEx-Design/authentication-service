package models

type User struct {
	ID      string `gorm:"primaryKey"`
	Email   string `gorm:"unique"`
	Name    string
	Picture string
}
