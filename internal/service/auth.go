package service

import (
	JWTServiceObjects "JWTService"
	"JWTService/internal/repository"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	_ "strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

const (
	salt       = "asdjkalhsd123123laksj"
	signingKey = "kaijdhOAS;KD'JJAKsjd"
	tokenTTL   = 12 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	UserID int    `json:"user_id"`
	IP     string `json:"ip"`
}

type AuthService struct {
	repo repository.Authorization
}

func NewAuthServices(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

// GenerateToken generates a new JWT token and refresh token for a given user GUID and client IP.
//
// Parameters:
// - guid: the user's GUID.
// - clientIp: the client's IP address.
//
// Returns:
// - string: the generated JWT token.
// - string: the refresh token (encoded in base64).
// - error: an error if any occurred.
func (s *AuthService) GenerateToken(guid string, clientIp string) (string, string, error) {

	claims := tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			Subject:   guid,
		},
		IP: clientIp,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

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
		ClientIp:     clientIp,
	})

	if err != nil {
		return "", "", err
	}

	return stringToken, base64.StdEncoding.EncodeToString(b), nil
}

// sendEmailWarning sends a warning email to the user when their IP address changes.
//
// Parameters:
// - userEmail: the email address of the user.
// - oldIP: the old IP address of the user.
// - newIP: the new IP address of the user.
// Return type: none
func sendEmailWarning(userEmail string, oldIP, newIP string) {
	logrus.Infof("Warning email sent to %s: IP address changed from %s to %s", userEmail, oldIP, newIP)
}

// RefreshToken refreshes the JWT and refresh tokens for a given user GUID and current IP.
//
// Parameters:
// - refreshToken: the refresh token bytes.
// - guid: the user's GUID.
// - curIP: the current IP address.
//
// Returns:
// - string: the new JWT token.
// - string: the new refresh token (encoded in base64).
// - error: an error if any occurred.
func (s *AuthService) RefreshToken(refreshToken []byte, guid string, curIP string) (string, string, error) {
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

	if session.ClientIp != curIP {
		sendEmailWarning("user@example.com", session.ClientIp, curIP)
	}

	newClaims := tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			Subject:   session.Guid,
		},
		IP: curIP,
	}
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS512, newClaims)

	stringToken, err := newToken.SignedString([]byte(signingKey)) // Token
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

// generateHash generates a hash token and a refresh token from a random byte slice.
//
// It takes no parameters.
// Returns a byte slice representing the random bytes, a byte slice representing the refresh hash token, and an error if any occurred.
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

// parseToken parses a JSON Web Token (JWT) and extracts the IP address from its claims.
//
// token is the JSON Web Token to be parsed.
// Returns the IP address extracted from the token claims as a string, and an error if any occurred.
func parseToken(token string) (string, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := parsedToken.Claims.(*tokenClaims)

	if !ok {
		return "", errors.New("token claims are not of type *tokenClaims")
	}

	return claims.IP, nil
}
