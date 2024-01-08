package cron

import (
	"context"
)

var ctx context.Context
var cancel context.CancelFunc

func Start() {
	ctx, cancel = context.WithCancel(context.Background())

	sch := NewScheduler()

	go sch.Schedule(ctx, NewSessionCleaner())
}

func Stop() {
	cancel()
}
