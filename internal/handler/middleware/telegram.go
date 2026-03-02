package middleware

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
	initdata "github.com/telegram-mini-apps/init-data-golang"
)

type TelegramMiddleware struct {
	BotToken string
	Log      *logrus.Logger
}

func NewTelegramMiddleware(botToken string, log *logrus.Logger) func(fiber.Ctx) error {
	mw := &TelegramMiddleware{
		BotToken: botToken,
		Log:      log,
	}
	return mw.Handle
}

func (m *TelegramMiddleware) Handle(c fiber.Ctx) error {
	authData := c.Get("Authorization")

	if authData == "" {
		return c.Next()
	}

	authData = strings.TrimPrefix(authData, "tma ")

	if err := initdata.Validate(authData, m.BotToken, 24*time.Hour); err != nil {
		m.Log.Warnf("Auth error: %+v", err)
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	data, err := initdata.Parse(authData)
	if err != nil {
		m.Log.Warnf("AuthData parse error: %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	c.Locals("telegram_id", data.User.ID)
	c.Locals("telegram_username", data.User.Username)
	c.Locals("telegram_firstname", data.User.FirstName)
	c.Locals("telegram_lastname", data.User.LastName)
	c.Locals("telegram_photoURL", data.User.PhotoURL)

	return c.Next()
}

func GetTelegramID(c fiber.Ctx) int64 {
	id, ok := c.Locals("telegram_id").(int64)
	if !ok {
		return 0
	}
	return id
}
