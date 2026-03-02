package main

import (
	"context"
	"fmt"
	"os/signal"
	"polytech_timetable/internal/app"
	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	cfg := app.NewConfig()
	logger := app.NewLogger(cfg)
	db := app.NewDatabase(cfg, logger)
	fiber := app.NewFiber(cfg)
	validator := app.NewValidator()
	bot := app.NewBot(cfg, logger)

	app.Start(bot, cfg, logger)

	app.Bootstrap(app.BootstrapConfig{
		Ctx:       ctx,
		DB:        db,
		App:       fiber,
		Log:       logger,
		Cfg:       cfg,
		Bot:       bot,
		Validator: validator,
	})

	go func() {
		if err := fiber.Listen(fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port)); err != nil {
			logger.Errorf("Failed to start server: %v", err)
			stop()
		}
	}()

	<-ctx.Done()
	logger.Info("Shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := fiber.ShutdownWithContext(shutdownCtx); err != nil {
		logger.Errorf("Shutdown error: %+v", err)
	}

	sqlDB, _ := db.DB()
	if sqlDB != nil {
		sqlDB.Close()
	}

	logger.Info("Shutdown complete")

}
