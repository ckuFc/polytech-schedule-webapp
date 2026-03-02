package model

import (
	"time"
)

type GetTeacherScheduleResponse struct {
	Lessons []LessonResponse `json:"lessons"`
}

type GetGroupsResponse struct {
	Groups []string `json:"groups"`
}

type GetLessonsForGroupResponse struct {
	Group   string           `json:"group"`
	Lessons []LessonResponse `json:"lessons"`
}

type LessonResponse struct {
	ID         int64     `json:"id"`
	ExternalId string    `json:"external_id"`
	Subject    string    `json:"subject"`
	Teacher    string    `json:"teacher"`
	Group      string    `json:"group"`
	Date       time.Time `json:"date"`
	LessonNum  int       `json:"lesson_num"`
	Room       string    `json:"room"`
	SubGroup   int       `json:"sub_group"`
	Zam        int       `json:"zam"`
}
