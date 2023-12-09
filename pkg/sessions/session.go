package sessions

import (
	"os"
	"time"
)

type Session struct {
	Id        int
	UserId    int
	Name      string
	Token     string
	ExpiredAt time.Time
}

func NewSession() Session {
	return Session{
		Name: os.Getenv("COOKIE_NAME"),
	}
}

func (s *Session) IsExpired() bool {
	return s.ExpiredAt.Before(time.Now())
}
