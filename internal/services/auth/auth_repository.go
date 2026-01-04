package auth

import "gorm.io/gorm"

type AuthRepository interface{}

type authRepository struct {
	db *gorm.DB
}

func NewAuthReposeitory(db *gorm.DB) AuthRepository {
	return &authRepository{
		db: db,
	}
}
