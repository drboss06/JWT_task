package repository

import (
	JWTServiceObjects "JWTService"
	"go.mongodb.org/mongo-driver/mongo"
)

type Authorization interface {
	SetSession(guid string, session JWTServiceObjects.Session) error
	GetSession(guid string) (JWTServiceObjects.Session, error)
	SetRefreshToken(refreshToken []byte, session JWTServiceObjects.Session) error
}

type Repository struct {
	Authorization
}

func NewRepository(db *mongo.Client) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
	}
}
