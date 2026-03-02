package worker

import (
	"context"
	"polytech_timetable/internal/metrics"
	"polytech_timetable/internal/usecase"
	"time"

	"github.com/sirupsen/logrus"
)

type ScheduleWorker struct {
	uc  *usecase.ScheduleUseCase
	log *logrus.Logger
}

func NewScheduleWorker(uc *usecase.ScheduleUseCase, log *logrus.Logger) *ScheduleWorker {
	return &ScheduleWorker{uc: uc, log: log}
}

func (w *ScheduleWorker) Start(ctx context.Context, interval time.Duration) {
	w.log.Infof("Start schedule worker (interval: %v)", interval)
	go func() {
		w.run()

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				w.run()
			}
		}
	}()
}

func (w *ScheduleWorker) run() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	if err := w.uc.SyncAllSchedules(ctx); err != nil {
		w.log.Errorf("Error in worker: %v", err)
		metrics.ParserSyncError.Inc()
	} else {
		metrics.ParserSyncSuccess.Inc()
	}
}
