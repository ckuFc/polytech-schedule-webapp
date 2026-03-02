package repository

import (
	"context"
	"polytech_timetable/internal/entity"

	"gorm.io/gorm"
)

type ReviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) *ReviewRepository {
	return &ReviewRepository{
		db: db,
	}
}

func (r *ReviewRepository) CreateReview(ctx context.Context, review *entity.TeacherReview) error {
	return r.db.WithContext(ctx).Create(review).Error
}

func (r *ReviewRepository) FindByTeacherID(ctx context.Context, teacherID int64) ([]entity.TeacherReview, error) {
	var reviews []entity.TeacherReview
	err := r.db.WithContext(ctx).
		Where("teacher_id = ?", teacherID).
		Order("created_at DESC").
		Find(&reviews).Error
	return reviews, err
}

func (r *ReviewRepository) FindByTeacherAndUser(ctx context.Context, teacherID int64, userID int64) (entity.TeacherReview, error) {
	var review entity.TeacherReview
	err := r.db.WithContext(ctx).
		Where("teacher_id = ? AND user_id = ?", teacherID, userID).
		Take(&review).Error
	return review, err
}

func (r *ReviewRepository) DeleteByTeacherAndUser(ctx context.Context, teacherID int64, userID int64) error {
	return r.db.WithContext(ctx).
		Where("teacher_id = ? AND user_id = ?", teacherID, userID).
		Delete(&entity.TeacherReview{}).Error
}
