package handler

import (
	"polytech_timetable/internal/handler/middleware"
	"polytech_timetable/internal/model"
	"polytech_timetable/internal/usecase"

	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	uc  *usecase.UserUseCase
	log *logrus.Logger
}

func NewUserController(uc *usecase.UserUseCase, log *logrus.Logger) *UserController {
	return &UserController{
		uc:  uc,
		log: log,
	}
}

func (h *UserController) GetMe(c fiber.Ctx) error {
	user := middleware.GetUser(c)
	if user == nil {
		return fiber.ErrUnauthorized
	}

	return c.JSON(model.WebResponse[fiber.Map]{
		Data: fiber.Map{
			"group":                 user.Group,
			"notifications_enabled": user.NotificationsEnabled,
		},
	})
}

func (h *UserController) SetGroup(c fiber.Ctx) error {
	user := middleware.GetUser(c)
	if user == nil {
		return fiber.ErrNotFound
	}

	var request model.SetGroupRequest
	if err := c.Bind().JSON(&request); err != nil {
		return fiber.ErrBadRequest
	}

	request.UserID = user.ID

	if err := h.uc.SetGroup(c.Context(), &request); err != nil {
		return MapDomainError(err)
	}

	return c.JSON(model.WebResponse[string]{Data: "success"})
}

func (h *UserController) ToggleNotifications(c fiber.Ctx) error {
	user := middleware.GetUser(c)
	if user == nil {
		return fiber.ErrNotFound
	}

	var request model.SetNotificationsSettingRequest
	if err := c.Bind().JSON(&request); err != nil {
		return fiber.ErrBadRequest
	}

	request.UserID = user.ID
	if err := h.uc.UpdateSettings(c.Context(), &request); err != nil {
		return MapDomainError(err)
	}
	return c.JSON(model.WebResponse[string]{Data: "success"})
}
