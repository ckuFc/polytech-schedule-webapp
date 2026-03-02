package handler

import (
	"polytech_timetable/internal/handler/middleware"
	"polytech_timetable/internal/model"
	"polytech_timetable/internal/usecase"

	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
)

type ScheduleController struct {
	uc  *usecase.ScheduleUseCase
	log *logrus.Logger
}

func NewScheduleController(uc *usecase.ScheduleUseCase, log *logrus.Logger) *ScheduleController {
	return &ScheduleController{
		uc:  uc,
		log: log,
	}
}

func (h *ScheduleController) GetSchedule(c fiber.Ctx) error {
	user := middleware.GetUser(c)
	if user == nil {
		return fiber.ErrUnauthorized
	}

	group := c.Query("group")

	if group == "" {
		group = user.Group
	}

	if group == "" {
		return fiber.ErrBadRequest
	}

	response, err := h.uc.GetLessonsForGroup(c.Context(), group)
	if err != nil {
		return MapDomainError(err)
	}

	if response.Lessons == nil {
		response.Lessons = []model.LessonResponse{}
	}

	return c.JSON(model.WebResponse[*model.GetLessonsForGroupResponse]{
		Data: response,
	})
}

func (h *ScheduleController) GetGroups(c fiber.Ctx) error {
	response, err := h.uc.GetGroupsList(c.Context())
	if err != nil {
		return MapDomainError(err)
	}

	if response.Groups == nil {
		response.Groups = []string{}
	}

	return c.JSON(model.WebResponse[*model.GetGroupsResponse]{
		Data: response,
	})
}

func (h *ScheduleController) GetTeacherSchedule(c fiber.Ctx) error {
	teacherID := fiber.Query[int64](c, "teacher_id", 0)
	if teacherID == 0 {
		return fiber.ErrBadRequest
	}

	response, err := h.uc.GetTeacherSchedule(c.Context(), teacherID)
	if err != nil {
		return MapDomainError(err)
	}

	return c.JSON(model.WebResponse[*model.GetTeacherScheduleResponse]{
		Data: response,
	})
}
