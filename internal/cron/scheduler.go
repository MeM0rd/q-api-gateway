package cron

import (
	"context"
	"github.com/MeM0rd/q-api-gateway/pkg/logger"
	"time"
)

type Scheduler interface {
	Schedule(context.Context, Task)
}

type scheduler struct {
	Logger *logger.Logger
}

func NewScheduler() Scheduler {
	return &scheduler{
		Logger: logger.New(),
	}
}

func (sch *scheduler) Schedule(ctx context.Context, task Task) {
	f, dur := task.Prepare()

	ticker := time.NewTicker(dur)

	for {
		select {
		case <-ticker.C:
			f()
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}
