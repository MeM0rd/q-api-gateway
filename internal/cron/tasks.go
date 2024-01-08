package cron

import (
	"github.com/MeM0rd/q-api-gateway/pkg/client/postgres"
	"github.com/MeM0rd/q-api-gateway/pkg/logger"
	"time"
)

type Task interface {
	Prepare() (func(), time.Duration)
}

type SessionCleaner struct {
	Logger *logger.Logger
}

func NewSessionCleaner() Task {
	return &SessionCleaner{
		Logger: logger.New(),
	}
}

func (sc *SessionCleaner) Prepare() (func(), time.Duration) {
	f := func() {
		q := `DELETE FROM sessions WHERE expired_at < $1`

		_, err := postgres.DB.Exec(q, time.Now())
		if err != nil {
			sc.Logger.Infof("session cleaner error: %v", err)
			return
		}
	}

	dur := time.Second * 60 * 10 // 10 min

	return f, dur
}
