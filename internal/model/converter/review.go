package converter

import (
	"polytech_timetable/internal/entity"
	"polytech_timetable/internal/model"
)

func TeacherToResponse(teacher *entity.Teacher, reviewCount int64) model.TeacherResponse {
	return model.TeacherResponse{
		ID:          teacher.ID,
		Name:        teacher.Name,
		ReviewCount: reviewCount,
	}
}

func ReviewToResponse(review *entity.TeacherReview, currentUserID int64) model.ReviewResponse {
	return model.ReviewResponse{
		ID:        review.ID,
		Comment:   review.Comment,
		CreatedAt: review.CreatedAt,
		IsMine:    review.UserID == currentUserID,
	}
}
