package entity

type TeacherReview struct {
	ID        int64  `gorm:"column:id;primaryKey;autoIncrement"`
	TeacherID int64  `gorm:"column:teacher_id;uniqueIndex:idx_teacher_user"`
	UserID    int64  `gorm:"column:user_id;uniqueIndex:idx_teacher_user"`
	Comment   string `gorm:"column:comment;size:255;not null"`
	Course    int    `gorm:"column:course;not null;default:0"`
	CreatedAt int64  `gorm:"column:created_at;autoCreateTime:milli"`

	Teacher Teacher `gorm:"foreignKey:TeacherID"`
	User    User    `gorm:"foreignKey:UserID"`
}
