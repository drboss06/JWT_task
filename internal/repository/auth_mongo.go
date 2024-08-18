package repository

import (
	JWTServiceObjects "JWTService"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) SetSession(guid string, session JWTServiceObjects.Session) error {
	query := `
		INSERT INTO sessions (guid, refresh_token, live_time, client_ip)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (guid)
		DO UPDATE SET
			refresh_token = EXCLUDED.refresh_token,
			live_time = EXCLUDED.live_time,
			client_ip = EXCLUDED.client_ip;
	`

	_, err := r.db.Exec(query, guid, session.RefreshToken, session.LiveTime, session.ClientIp)
	if err != nil {
		return fmt.Errorf("failed to set session: %w", err)
	}

	return nil
}

func (r *AuthPostgres) GetSession(guid string) (JWTServiceObjects.Session, error) {
	query := `
		SELECT *
		FROM sessions
		WHERE guid = $1
	`

	var session JWTServiceObjects.Session
	err := r.db.Get(&session, query, guid)
	if err != nil {
		return JWTServiceObjects.Session{}, fmt.Errorf("failed to get session: %w", err)
	}

	return session, nil
}

func (r *AuthPostgres) SetRefreshToken(refreshToken []byte, session JWTServiceObjects.Session) error {
	query := `
		UPDATE sessions
		SET refresh_token = $1, live_time = $2
		WHERE refresh_token = $3
	`
	_, err := r.db.Exec(query, session.RefreshToken, session.LiveTime, refreshToken)
	if err != nil {
		return fmt.Errorf("failed to update refresh token: %w", err)
	}

	return nil
}
