package repository

import (
	"context"
	"polytech_timetable/internal/entity"
	"time"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) SetUserGroup(ctx context.Context, userID int64, group string) error {
	return r.db.WithContext(ctx).
		Model(&entity.User{}).
		Where("id = ?", userID).
		Update("group", group).Error
}

func (r *UserRepository) UpdateSettings(ctx context.Context, userID int64, enabled bool) error {
	return r.db.WithContext(ctx).
		Model(&entity.User{}).
		Where("id = ?", userID).
		Update("notifications_enabled", enabled).Error
}

func (r *UserRepository) FindByUserID(ctx context.Context, userID int64) (entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).
		Where("id = ?", userID).
		Take(&user).Error
	return user, err
}

func (r *UserRepository) FindByGroupWithNotifications(ctx context.Context, group string) ([]entity.User, error) {
	var users []entity.User
	err := r.db.WithContext(ctx).
		Where("\"group\" = ? AND notifications_enabled = true", group).
		Find(&users).Error
	return users, err
}

func (r *UserRepository) FindOrCreate(ctx context.Context, data *entity.User) (*entity.User, error) {
	var user entity.User
	tx := r.db.WithContext(ctx).
		Where(entity.User{TelegramID: data.TelegramID}).
		Attrs(entity.User{
			UserName:  data.UserName,
			FirstName: data.FirstName,
			LastName:  data.LastName,
			PhotoURL:  data.PhotoURL,
		}).
		FirstOrCreate(&user)

	if tx.Error != nil {
		return nil, tx.Error
	}

	if tx.RowsAffected == 0 {
		err := r.db.WithContext(ctx).
			Model(&user).
			Update("updated_at", time.Now().UnixMilli()).Error

		if err != nil {
			return nil, err
		}
	}

	return &user, nil
}
