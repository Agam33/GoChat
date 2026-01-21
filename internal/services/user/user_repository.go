package user

import (
	"context"
	"go-chat/internal/model"
	"go-chat/internal/utils"
	"go-chat/internal/utils/types"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserRooms(context.Context, uint64, *types.Pagination) ([]model.UserRoom, error)
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

func (r *userRepository) GetUserRooms(ctx context.Context, userId uint64, pagination *types.Pagination) ([]model.UserRoom, error) {
	var rooms []model.UserRoom
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userId).
		Order("created_at DESC").
		Preload("Room").
		Limit(pagination.Limit).
		Offset(utils.PageOffset(pagination.Page, pagination.Limit)).
		Find(&rooms).Error
	if err != nil {
		return []model.UserRoom{}, err
	}

	return rooms, nil
}

func (ur *userRepository) GetById(ctx context.Context, userId uint64) (*model.User, error) {
	var usr model.User
	if err := ur.db.WithContext(ctx).First(&usr, userId).Error; err != nil {
		return nil, err
	}

	return &usr, nil
}
