package handler

import (
	"errors"
	"polytech_timetable/pkg/customerrors"

	"github.com/gofiber/fiber/v3"
)

func MapDomainError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, customerrors.ErrValidation):
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	case errors.Is(err, customerrors.ErrNotFound):
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	case errors.Is(err, customerrors.ErrAlreadyExists):
		return fiber.NewError(fiber.StatusConflict, err.Error())
	default:
		return fiber.ErrInternalServerError
	}
}
