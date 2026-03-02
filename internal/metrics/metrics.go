package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"gorm.io/gorm"
)

var (
	GroupChangesTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "polytech_group_changes_total",
		Help: "Total changes count",
	}, []string{"group_name"})

	ParserSyncSuccess = promauto.NewCounter(prometheus.CounterOpts{
		Name: "polytech_parser_sync_success_total",
		Help: "Total count of successfull parser sync",
	})

	ParserSyncError = promauto.NewCounter(prometheus.CounterOpts{
		Name: "polytech_parser_sync_errors_total",
		Help: "Total count of failed parser sync",
	})

	HTTPRequestTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "polytech_http_requests_total",
		Help: "Total count of http requests to API",
	}, []string{"method", "path", "status"})

	HTTPRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "polytech_http_request_duration_seconds",
		Help:    "Duration of http requests",
		Buckets: []float64{0.01, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5},
	}, []string{"method", "path"})
)

func InitDBMetrics(db *gorm.DB) {
	promauto.NewGaugeFunc(prometheus.GaugeOpts{
		Name: "polytech_users_total",
		Help: "Total registered users in DB",
	}, func() float64 {
		var count int64
		db.Table("users").Count(&count)
		return float64(count)
	})

	promauto.NewGaugeFunc(prometheus.GaugeOpts{
		Name: "polytech_users_active_today",
		Help: "Users who opened the app today",
	}, func() float64 {
		now := time.Now()
		startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		startMillis := startOfDay.UnixMilli()

		var count int64
		db.Table("users").Where("updated_at >= ?", startMillis).Count(&count)
		return float64(count)
	})
}
