package app

import (
	"context"
	"polytech_timetable/internal/handler"
	"polytech_timetable/internal/handler/middleware"
	"polytech_timetable/internal/handler/route"
	"polytech_timetable/internal/metrics"
	"polytech_timetable/internal/repository"
	"polytech_timetable/internal/usecase"
	"polytech_timetable/internal/worker"
	"polytech_timetable/pkg/polytech"
	"time"

	"github.com/go-playground/validator/v10"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	Ctx       context.Context
	DB        *gorm.DB
	App       *fiber.App
	Log       *logrus.Logger
	Cfg       Config
	Validator *validator.Validate
	Bot       *tgbotapi.BotAPI
}

func Bootstrap(boot BootstrapConfig) route.RouteConfig {
	polytechClient := polytech.NewClient(boot.Cfg.App.BaseURL)

	// setup repositories
	scheduleRepository := repository.NewScheduleRepository(boot.DB)
	userRepository := repository.NewUserRepository(boot.DB)
	teacherRepository := repository.NewTeacherRepository(boot.DB)
	reviewRepository := repository.NewReviewRepository(boot.DB)

	// setup usecases
	scheduleUseCase := usecase.NewScheduleUseCase(boot.Log, polytechClient, scheduleRepository, userRepository, teacherRepository, boot.Bot)
	reviewUseCase := usecase.NewReviewUseCase(boot.Log, boot.Validator, userRepository, reviewRepository, teacherRepository)
	userUseCase := usecase.NewUserUseCase(boot.Log, boot.Validator, userRepository)

	//setup controllers
	scheduleController := handler.NewScheduleController(scheduleUseCase, boot.Log)
	reviewController := handler.NewReviewController(reviewUseCase, boot.Log)
	userController := handler.NewUserController(userUseCase, boot.Log)

	//setup middlewares
	limMiddleware := middleware.NewLimiterMiddleware(boot.Log, boot.Cfg.LM.RPS, boot.Cfg.LM.Burst, boot.Cfg.LM.Size)
	tgMiddleware := middleware.NewTelegramMiddleware(boot.Cfg.TG.BotToken, boot.Log)
	regMiddleware := middleware.NewRegisterMiddleware(boot.Log, userRepository)
	metMiddleware := middleware.NewMetricsMiddleware(boot.Log)

	cronWorker := worker.NewScheduleWorker(scheduleUseCase, boot.Log)
	cronWorker.Start(boot.Ctx, 10*time.Minute)

	//setup metrics
	metrics.InitDBMetrics(boot.DB)

	routeConfig := route.RouteConfig{
		App:                boot.App,
		ScheduleController: scheduleController,
		ReviewController:   reviewController,
		UserController:     userController,
		LimiterMiddleware:  limMiddleware,
		TelegramMiddleware: tgMiddleware,
		RegisterMiddleware: regMiddleware,
		MetricsMiddleware:  metMiddleware,
		CORSOrigin:         boot.Cfg.HTTP.CORSOrigin,
	}
	routeConfig.Setup()

	return routeConfig
}
