package repository

import (
	JWTServiceObjects "JWTService"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Authorization interface {
	SetSession(guid string, session JWTServiceObjects.Session) error
	GetSession(guid string) (JWTServiceObjects.Session, error)
	SetRefreshToken(refreshToken []byte, session JWTServiceObjects.Session) error
}

type Repository struct {
	Authorization
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
	}
}
