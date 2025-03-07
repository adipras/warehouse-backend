package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"-"`
}

type UserSwagger struct {
	Name     string `json:"name" example:"John Doe"`
	Email    string `json:"email" gorm:"unique" example:"mail@example.com"`
	Password string `json:"password" example:"password123"`
}

type LoginCredentials struct {
	Email    string `json:"email" example:"admin1@email.com"`
	Password string `json:"password" example:"bismillah"`
}
