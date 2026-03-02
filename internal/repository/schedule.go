package repository

import (
	"context"
	"polytech_timetable/internal/entity"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ScheduleRepository struct {
	db *gorm.DB
}

func NewScheduleRepository(db *gorm.DB) *ScheduleRepository {
	return &ScheduleRepository{
		db: db,
	}
}

func (r *ScheduleRepository) ReplaceAll(ctx context.Context, lessons []entity.Lesson) error {
	if len(lessons) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "external_id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"subject", "teacher", "group", "date",
				"lesson_num", "sub_group", "zam", "updated_at",
			}),
		}).
		CreateInBatches(lessons, 1000).Error
}

func (r *ScheduleRepository) GetLessonsByGroup(ctx context.Context, group string, from, to time.Time) ([]entity.Lesson, error) {
	var lessons []entity.Lesson
	err := r.db.WithContext(ctx).
		Where("\"group\" = ? AND date >= ? AND date <= ?", group, from, to).
		Order("date ASC, lesson_num ASC").
		Find(&lessons).Error
	return lessons, err
}

func (r *ScheduleRepository) SearchTeacher(ctx context.Context, teacherName string) ([]entity.Lesson, error) {
	var lessons []entity.Lesson
	err := r.db.WithContext(ctx).
		Where("teacher ILIKE ?", "%"+teacherName+"%").
		Order("date ASC").
		Limit(50).
		Find(&lessons).Error
	return lessons, err
}

func (r *ScheduleRepository) GetAllGroups(ctx context.Context) ([]string, error) {
	var groups []string
	err := r.db.WithContext(ctx).
		Model(&entity.Lesson{}).
		Distinct("\"group\"").
		Order("\"group\" ASC").
		Pluck("\"group\"", &groups).Error
	return groups, err
}

func (r *ScheduleRepository) GetFutureLessons(ctx context.Context) ([]entity.Lesson, error) {
	var lessons []entity.Lesson
	err := r.db.WithContext(ctx).
		Where("date >= ?", time.Now().Truncate(24*time.Hour)).
		Find(&lessons).Error
	return lessons, err
}
