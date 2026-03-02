package converter

import (
	"polytech_timetable/internal/entity"
	"polytech_timetable/internal/model"
)

func LessonsToResponse(lessons []entity.Lesson) []model.LessonResponse {
	response := make([]model.LessonResponse, 0, len(lessons))
	for i := range lessons {
		response = append(response, model.LessonResponse{
			ID:         lessons[i].ID,
			ExternalId: lessons[i].ExternalID,
			Subject:    lessons[i].Subject,
			Teacher:    lessons[i].Teacher,
			Group:      lessons[i].Group,
			Date:       lessons[i].Date,
			LessonNum:  lessons[i].LessonNum,
			Room:       lessons[i].Room,
			SubGroup:   lessons[i].SubGroup,
			Zam:        lessons[i].Zam,
		})
	}
	return response
}
