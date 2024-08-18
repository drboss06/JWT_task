package service

import (
	JWTServiceObjects "JWTService"
	"JWTService/internal/repository"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	_ "strconv"
	"time"
)

const (
	salt       = "asdjkalhsd123123laksj"
	signingKey = "kaijdhOAS;KD'JJAKsjd"
	tokenTTL   = 12 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	UserID int `json:"user_id"`
}

type AuthService struct {
	repo repository.Authorization
}

func NewAuthServices(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) GenerateToken(guid string, clientIp string) (string, string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(tokenTTL).Unix(),
		Subject:   guid,
	})

	stringToken, err := token.SignedString([]byte(signingKey)) // Token
	if err != nil {

		logrus.Fatalf("failed to generate token: %s", err.Error())

		return "", "", err
	}

	b, hashTokenRefresh, err := generateHash()

	if err != nil {
		return "", "", err
	}

	s.repo.SetSession(guid, JWTServiceObjects.Session{
		RefreshToken: hashTokenRefresh,
		LiveTime:     time.Now().Add(tokenTTL),
	})

	if err != nil {
		return "", "", err
	}

	return stringToken, base64.StdEncoding.EncodeToString(b), nil
}

func (s *AuthService) RefreshToken(refreshToken []byte, guid string) (string, string, error) {
	session, err := s.repo.GetSession(guid)

	if err != nil {
		return "", "", err
	}

	err = bcrypt.CompareHashAndPassword(session.RefreshToken, refreshToken)

	if err != nil {
		return "", "", err

	}

	if session.LiveTime.Before(time.Now()) {
		return "", "", errors.New("refresh token expired")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(tokenTTL).Unix(),
		Subject:   session.Guid,
	})

	stringToken, err := token.SignedString([]byte(signingKey)) // Token
	if err != nil {
		return "", "", err
	}

	b, hashTokenRefresh, err := generateHash()

	if err != nil {
		return "", "", err
	}

	err = s.repo.SetRefreshToken(session.RefreshToken, JWTServiceObjects.Session{
		RefreshToken: hashTokenRefresh,
		LiveTime:     time.Now().Add(tokenTTL),
	})

	if err != nil {
		return "", "", err
	}
	return stringToken, base64.StdEncoding.EncodeToString(b), nil
}

func generateHash() ([]byte, []byte, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate random bytes: %w", err)
	}

	hashTokenRefresh, err := bcrypt.GenerateFromPassword(b, 10)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate refresh hash token: %w", err)
	}

	return b, hashTokenRefresh, nil
}
