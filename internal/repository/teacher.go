package repository

import (
	"context"
	"polytech_timetable/internal/entity"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TeacherRepository struct {
	db *gorm.DB
}

func NewTeacherRepository(db *gorm.DB) *TeacherRepository {
	return &TeacherRepository{
		db: db,
	}
}

func (r *TeacherRepository) FindOrCreate(ctx context.Context, name string) (*entity.Teacher, error) {
	var teacher entity.Teacher
	err := r.db.WithContext(ctx).
		Where("name = ?", name).
		FirstOrCreate(&teacher, entity.Teacher{Name: name}).Error
	return &teacher, err
}

func (r *TeacherRepository) FindByTeacherID(ctx context.Context, teacherID int64) (entity.Teacher, error) {
	var teacher entity.Teacher
	err := r.db.WithContext(ctx).
		Where("id = ?", teacherID).
		Take(&teacher).Error
	return teacher, err
}

func (r *TeacherRepository) FindAll(ctx context.Context) ([]entity.Teacher, error) {
	var teachers []entity.Teacher
	err := r.db.WithContext(ctx).
		Order("name ASC").
		Find(&teachers).Error
	return teachers, err
}

func (r *TeacherRepository) CountReviews(ctx context.Context, teacherID int64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.TeacherReview{}).
		Where("teacher_id", teacherID).
		Count(&count).Error
	return count, err
}

func (r *TeacherRepository) FindAllWithReviewCount(ctx context.Context) ([]entity.TeacherWithReviewCount, error) {
	var result []entity.TeacherWithReviewCount
	err := r.db.WithContext(ctx).
		Model(&entity.Teacher{}).
		Select("teachers.*, COUNT(teacher_reviews.id) as review_count").
		Joins("LEFT JOIN teacher_reviews ON teacher_reviews.teacher_id = teachers.id").
		Group("teachers.id").
		Order("teachers.name ASC").
		Scan(&result).Error
	return result, err
}

func (r *TeacherRepository) SaveBatchIgnore(ctx context.Context, teachers []entity.Teacher) error {
	if len(teachers) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		CreateInBatches(teachers, 300).Error
}
