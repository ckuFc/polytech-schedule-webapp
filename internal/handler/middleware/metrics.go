package middleware

import (
	"errors"
	"polytech_timetable/internal/metrics"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
)

type MetricsMiddleware struct {
	Log *logrus.Logger
}

func NewMetricsMiddleware(log *logrus.Logger) func(fiber.Ctx) error {
	mw := &MetricsMiddleware{
		Log: log,
	}
	return mw.Handle
}

func (m *MetricsMiddleware) Handle(c fiber.Ctx) error {
	start := time.Now()

	err := c.Next()

	duration := time.Since(start).Seconds()

	status := c.Response().StatusCode()
	if err != nil {
		var fiberErr *fiber.Error
		if errors.As(err, &fiberErr) {
			status = fiberErr.Code
		} else {
			status = fiber.StatusInternalServerError
		}
	}

	method := c.Method()
	path := c.Route().Path

	metrics.HTTPRequestTotal.WithLabelValues(method, path, strconv.Itoa(status)).Inc()
	metrics.HTTPRequestDuration.WithLabelValues(method, path).Observe(duration)

	return err
}
