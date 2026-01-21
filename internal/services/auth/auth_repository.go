package auth

import (
	"context"
	"errors"
	"go-chat/internal/model"

	"gorm.io/gorm"
)

type AuthRepository interface {
	SignUp(context.Context, *model.User) error
	FindByUsername(context.Context, string) (*model.User, error)
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthReposeitory(db *gorm.DB) AuthRepository {
	return &authRepository{
		db: db,
	}
}

func (r *authRepository) SignUp(ctx context.Context, user *model.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		if errors.Is(err, gorm.ErrCheckConstraintViolated) {
			return gorm.ErrCheckConstraintViolated
		}

		return err
	}
	return nil
}

func (r *authRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &user, nil
}
