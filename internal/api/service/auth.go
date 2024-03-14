package service

import (
	"errors"
	"fmt"
	"github.com/Warh40k/vk-intern-filmotecka/internal/api/repository"
	"github.com/Warh40k/vk-intern-filmotecka/internal/domain"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

type AuthService struct {
	repos repository.Authorization
	log   *slog.Logger
}

var (
	ErrUserNotFound = fmt.Errorf("specified user not found")
)

func NewAuthService(repos repository.Authorization, log *slog.Logger) *AuthService {
	return &AuthService{repos: repos, log: log}
}

const (
	salt       = "fjlsj2374slfjsd728vvnts"
	tokenTTL   = 12 * time.Hour
	signingKey = "j370sdfs34472fshvlruso043275fhka"
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

func GenerateJWT(user domain.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{jwt.StandardClaims{
		ExpiresAt: time.Now().Add(tokenTTL).Unix(),
		IssuedAt:  time.Now().Unix(),
	}, user.Id})
	return token.SignedString([]byte(signingKey))
}

func CheckJWT(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return -1, err
	}
	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return -1, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *AuthService) SignUp(user domain.User) error {
	hash, err := HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}
	user.PasswordHash = hash

	_, err = s.repos.SignUp(user)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) SignIn(username, password string) (string, error) {
	user, err := s.repos.GetUserByUsername(username)
	if err != nil {
		return "", ErrUserNotFound
	}
	hash := user.PasswordHash
	if CheckPassword(password, hash) {
		token, err := GenerateJWT(user)
		if err != nil {
			return "", ErrInternal
		}
		return token, nil
	}

	return "", ErrUnauthorized
}
