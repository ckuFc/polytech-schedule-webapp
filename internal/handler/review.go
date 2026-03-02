package handler

import (
	"polytech_timetable/internal/handler/middleware"
	"polytech_timetable/internal/model"
	"polytech_timetable/internal/usecase"

	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
)

type ReviewController struct {
	uc  *usecase.ReviewUseCase
	log *logrus.Logger
}

func NewReviewController(uc *usecase.ReviewUseCase, log *logrus.Logger) *ReviewController {
	return &ReviewController{
		uc:  uc,
		log: log,
	}
}

func (h *ReviewController) ListTeachers(ctx fiber.Ctx) error {
	response, err := h.uc.ListAllTeachers(ctx.Context())
	if err != nil {
		return MapDomainError(err)
	}

	return ctx.JSON(model.WebResponse[[]model.TeacherResponse]{Data: response})
}

func (h *ReviewController) GetReviews(c fiber.Ctx) error {
	user := middleware.GetUser(c)
	if user == nil {
		return fiber.ErrUnauthorized
	}

	teacherID := fiber.Query[int64](c, "teacher_id", 0)

	if teacherID == 0 {
		return fiber.ErrBadRequest
	}

	request := &model.GetReviewsRequest{TeacherID: teacherID}
	response, err := h.uc.GetReviews(c.Context(), request, user.ID)
	if err != nil {
		return MapDomainError(err)
	}

	return c.JSON(model.WebResponse[*model.GetReviewsResponse]{Data: response})
}

func (h *ReviewController) CreateReview(c fiber.Ctx) error {
	user := middleware.GetUser(c)
	if user == nil {
		return fiber.ErrUnauthorized
	}

	request := new(model.CreateReviewRequest)

	if err := c.Bind().JSON(request); err != nil {
		h.log.Warnf("Failed to parse request body: %+v", err)
		return fiber.ErrBadRequest
	}

	request.UserID = user.ID

	if err := h.uc.CreateReview(c.Context(), request); err != nil {
		return MapDomainError(err)
	}

	return c.JSON(model.WebResponse[string]{Data: "success"})
}

func (h *ReviewController) DeleteReview(ctx fiber.Ctx) error {
	user := middleware.GetUser(ctx)
	if user == nil {
		return fiber.ErrUnauthorized
	}

	request := new(model.DeleteReviewRequest)

	if err := ctx.Bind().JSON(request); err != nil {
		h.log.Warnf("Failed to parse request body: %+v", err)
		return fiber.ErrBadRequest
	}

	request.UserID = user.ID

	if err := h.uc.DeleteReview(ctx.Context(), request); err != nil {
		return MapDomainError(err)
	}

	return ctx.JSON(model.WebResponse[string]{Data: "success"})
}
