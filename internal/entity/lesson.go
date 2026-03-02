package entity

import "time"

type Lesson struct {
	ID         int64  `gorm:"primaryKey"`
	ExternalID string `gorm:"unique"`
	Subject    string
	Teacher    string
	Group      string    `gorm:"column:group;index"`
	Date       time.Time `gorm:"index"`
	LessonNum  int
	Room       string
	SubGroup   int
	Zam        int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
