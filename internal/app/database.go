package app

import (
	"fmt"
	"polytech_timetable/internal/entity"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/prometheus"
)

func NewDatabase(cfg Config, log *logrus.Logger) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Europe/Moscow",
		cfg.PG.DatabaseHost,
		cfg.PG.DatabaseUser,
		cfg.PG.DatabasePassword,
		cfg.PG.DatabaseDb,
		cfg.PG.DatabasePort,
	)

	var db *gorm.DB
	var err error

	maxRetries := 10
	for i := range maxRetries {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.New(&logrusWriter{Logger: log}, logger.Config{
				SlowThreshold:             time.Second * 5,
				Colorful:                  true,
				IgnoreRecordNotFoundError: true,
				ParameterizedQueries:      true,
				LogLevel:                  logger.Info,
			}),
		})
		if err == nil {
			break
		}

		log.Warnf("Attempt %d/%d to connect to database failed: %v", i+1, maxRetries, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Info("Successfull connect to database")

	if err := db.AutoMigrate(&entity.Lesson{}, &entity.User{}, &entity.Teacher{}, &entity.TeacherReview{}); err != nil {
		log.Fatalf("Failed auto migration : %v", err)
	}

	log.Info("Successfull auto-migration")

	db.Use(prometheus.New(prometheus.Config{
		DBName:          cfg.PG.DatabaseDb,
		RefreshInterval: 15,
		StartServer:     false,
	}))

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to set sql configuration: %+v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db
}

type logrusWriter struct {
	Logger *logrus.Logger
}

func (l *logrusWriter) Printf(message string, args ...any) {
	l.Logger.Tracef(message, args...)
}
