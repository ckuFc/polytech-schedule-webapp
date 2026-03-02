package middleware

import (
	"sync"

	"github.com/gofiber/fiber/v3"
	lru "github.com/hashicorp/golang-lru"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

type LimiterMiddleware struct {
	Log      *logrus.Logger
	limiters *lru.Cache
	rps      float64
	burst    int
	mu       sync.Mutex
}

func NewLimiterMiddleware(log *logrus.Logger, rps float64, burst int, cacheSize int) func(fiber.Ctx) error {
	cache, err := lru.New(cacheSize)
	if err != nil {
		log.Fatalf("Failed to create rate limiter cache: %+v", err)
	}
	mw := &LimiterMiddleware{
		Log:      log,
		limiters: cache,
		rps:      rps,
		burst:    burst,
		mu:       sync.Mutex{},
	}

	return mw.Handle
}

func (m *LimiterMiddleware) Handle(c fiber.Ctx) error {
	tgID := GetTelegramID(c)
	limiter := m.GetLimiter(tgID)

	if !limiter.Allow() {
		return fiber.ErrTooManyRequests
	}
	return c.Next()
}

func (m *LimiterMiddleware) GetLimiter(tgID int64) *rate.Limiter {
	m.mu.Lock()
	defer m.mu.Unlock()

	if limiter, ok := m.limiters.Get(tgID); ok {
		return limiter.(*rate.Limiter)
	}

	limiter := rate.NewLimiter(rate.Limit(m.rps), m.burst)
	m.limiters.Add(tgID, limiter)
	return limiter
}
