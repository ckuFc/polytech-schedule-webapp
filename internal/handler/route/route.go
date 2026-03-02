package route

import (
	"polytech_timetable/internal/handler"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type RouteConfig struct {
	App                *fiber.App
	ScheduleController *handler.ScheduleController
	UserController     *handler.UserController
	ReviewController   *handler.ReviewController

	LimiterMiddleware  func(fiber.Ctx) error
	TelegramMiddleware func(fiber.Ctx) error
	RegisterMiddleware func(fiber.Ctx) error
	MetricsMiddleware  func(fiber.Ctx) error

	CORSOrigin string
}

func corsMiddleware(origin string) fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins: []string{origin},
		AllowHeaders: []string{"Origin, Content-Type, Accept"},
		AllowMethods: []string{"GET, POST, HEAD, PUT, DELETE, PATCH"},
	})
}

func prometheusHandler() fiber.Handler {
	handler := promhttp.Handler()
	return func(c fiber.Ctx) error {
		return adaptor.HTTPHandler(handler)(c)
	}
}

func (c *RouteConfig) Setup() {
	c.App.Use(corsMiddleware(c.CORSOrigin))

	c.App.Use(c.MetricsMiddleware)

	c.App.Get("/metrics", prometheusHandler())

	c.App.Use(c.TelegramMiddleware)
	c.App.Use(c.LimiterMiddleware)

	c.SetupAuthRoute()
}

func (c *RouteConfig) SetupAuthRoute() {
	api := c.App.Group("/api/v1")
	api.Get("/groups", c.ScheduleController.GetGroups)

	api.Use(c.RegisterMiddleware)

	api.Get("/schedule", c.ScheduleController.GetSchedule)
	api.Get("/schedule/teacher", c.ScheduleController.GetTeacherSchedule)

	api.Get("/user/me", c.UserController.GetMe)
	api.Post("/user/group", c.UserController.SetGroup)
	api.Post("/user/notifications", c.UserController.ToggleNotifications)

	api.Get("/teachers", c.ReviewController.ListTeachers)
	api.Get("/reviews", c.ReviewController.GetReviews)
	api.Post("/reviews", c.ReviewController.CreateReview)
	api.Delete("/reviews", c.ReviewController.DeleteReview)
}
