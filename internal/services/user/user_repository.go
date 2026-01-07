package user

import (
	"context"
	"go-chat/internal/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetById(context.Context, uint64) (*model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (ur *userRepository) GetById(ctx context.Context, userId uint64) (*model.User, error) {
	var usr model.User
	if err := ur.db.WithContext(ctx).First(&usr, userId).Error; err != nil {
		return nil, err
	}

	return &usr, nil
}
