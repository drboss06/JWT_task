package service

import (
	"JWTService/internal/repository"
)

type Authorization interface {
	GenerateToken(guid string, clientIp string) (string, string, error)
	RefreshToken(refreshToken []byte, guid string) (string, string, error)
}

type TodoItem interface {
}

type Service struct {
	Authorization
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthServices(repos.Authorization),
	}
}
