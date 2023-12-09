package sessions

import (
	"errors"
	"github.com/MeM0rd/q-api-gateway/pkg/client/postgres"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func CreateSession(userId *int) (Session, error) {
	session := NewSession()
	var err error

	sessionToken := uuid.NewString()
	expiredAt := time.Now().Add(5 * time.Hour)

	q := `INSERT INTO sessions(user_id, token, expired_at) VALUES ($1, $2, $3) RETURNING id, token, user_id expired_at`

	err = postgres.DB.QueryRow(q, userId, sessionToken, expiredAt).Scan(&session.Id, &session.Token, &session.UserId, &session.ExpiredAt)
	if err != nil {
		return Session{}, err
	}

	session.ExpiredAt = expiredAt

	return session, nil
}

func CheckSession(r *http.Request) (Session, error) {
	session := NewSession()

	cookie, err := r.Cookie("q-svc-token")
	if err != nil {
		return Session{}, err
	}

	q := `SELECT id, token, expired_at, user_id FROM sessions WHERE token = $1`

	err = postgres.DB.QueryRow(q, cookie.Value).Scan(&session.Id, &session.Token, &session.ExpiredAt, &session.UserId)
	if err != nil {
		return Session{}, err
	}

	if session.IsExpired() {
		q = `DELETE FROM sessions WHERE token = $1`

		postgres.DB.Exec(q, session.Token)
		return Session{}, errors.New("token is expired")
	}

	return session, nil
}

func DeleteSession(token string) error {

	q := `DELETE FROM sessions WHERE token = $1`

	err := postgres.DB.QueryRow(q, token).Err()
	if err != nil {
		return err
	}

	return nil
}

func GetToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie("q-svc-token")
	if err != nil {
		return "", err
	}

	return cookie.Value, nil
}
