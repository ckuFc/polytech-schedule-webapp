package middleware

import (
	"polytech_timetable/internal/entity"
	"polytech_timetable/internal/repository"

	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
)

type RegisterMiddleware struct {
	log      *logrus.Logger
	userRepo *repository.UserRepository
}

func NewRegisterMiddleware(log *logrus.Logger, userRepo *repository.UserRepository) func(fiber.Ctx) error {
	mw := &RegisterMiddleware{
		log:      log,
		userRepo: userRepo,
	}
	return mw.Handle
}

func (m *RegisterMiddleware) Handle(c fiber.Ctx) error {
	tgID := GetTelegramID(c)

	if tgID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	username, _ := c.Locals("telegram_username").(string)
	firstname, _ := c.Locals("telegram_firstname").(string)
	lastname, _ := c.Locals("telegram_lastname").(string)
	photoURL, _ := c.Locals("telegram_photoURL").(string)

	userData := entity.User{
		TelegramID: tgID,
		UserName:   username,
		FirstName:  firstname,
		LastName:   lastname,
		PhotoURL:   photoURL,
	}

	user, err := m.userRepo.FindOrCreate(c.Context(), &userData)
	if err != nil {
		m.log.Errorf("DB error finding/creating user: %+v", err)
		return fiber.ErrInternalServerError
	}

	c.Locals("user", user)
	return c.Next()
}

func GetUser(c fiber.Ctx) *entity.User {
	user, ok := c.Locals("user").(*entity.User)
	if !ok {
		return nil
	}
	return user
}
