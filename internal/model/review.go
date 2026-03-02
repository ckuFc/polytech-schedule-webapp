package model

type SearchTeacherRequest struct {
	Query string `json:"query" validate:"required,min=2,max=255"`
}

type CreateReviewRequest struct {
	UserID    int64  `json:"-"`
	TeacherID int64  `json:"teacher_id" validate:"required"`
	Comment   string `json:"comment" validate:"required,min=1,max=255"`
}

type DeleteReviewRequest struct {
	UserID    int64 `json:"-"`
	TeacherID int64 `json:"teacher_id" validate:"required"`
}

type GetReviewsRequest struct {
	TeacherID int64 `json:"teacher_id" validate:"required"`
}

type TeacherResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	ReviewCount int64  `json:"review_count"`
}

type ReviewResponse struct {
	ID        int64  `json:"id"`
	Comment   string `json:"comment"`
	CreatedAt int64  `json:"created_at"`
	IsMine    bool   `json:"is_mine"`
}

type GetReviewsResponse struct {
	Teacher TeacherResponse  `json:"teacher"`
	Reviews []ReviewResponse `json:"reviews"`
}
