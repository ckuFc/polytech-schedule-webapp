package entity

type Teacher struct {
	ID        int64  `gorm:"column:id;primaryKey;autoIncrement"`
	Name      string `gorm:"column:name;uniqueIndex;size:255;not null"`
	CreatedAt int64  `gorm:"column:created_at;autoCreateTime:milli"`
}

type TeacherWithReviewCount struct {
	Teacher
	ReviewCount int64 `gorm:"column:review_count"`
}
