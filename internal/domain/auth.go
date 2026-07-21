package domain

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepository interface {
	GetByEmail(ctx context.Context, email string) (*AuthUser, error)
}

type AuthService struct {
	repo      AuthRepository
	jwtSecret []byte
}

type AuthToken struct {
	Token string `json:"token"`
}

type AuthUser struct {
	ID           string
	Email        string
	Name         string
	PasswordHash string
}

type AuthContext struct {
	UserID string
	Email  string
	Name   string
}

func GetAuthContext(ctx context.Context) (*AuthContext, error) {
	v := ctx.Value(AuthContextKey)
	if v == nil {
		return nil, errors.New("no auth context")
	}
	ac, ok := v.(*AuthContext)
	if !ok {
		return nil, errors.New("invalid auth context")
	}
	return ac, nil
}

func NewAuthService(r AuthRepository, jwtSecret []byte) *AuthService {
	return &AuthService{repo: r, jwtSecret: jwtSecret}
}

func (s *AuthService) authenticate(ctx context.Context, email string, password string) (*AuthUser, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return user, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return user, errors.New("invalid email or password")
		}
		return user, err
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, email string, password string) (AuthToken, error) {
	var authToken AuthToken

	user, err := s.authenticate(ctx, email, password)
	if err != nil {
		return authToken, err
	}

	authToken.Token, err = s.GenerateToken(user)
	if err != nil {
		return authToken, err
	}

	return authToken, nil
}

type JWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
}

type contextKey string

const AuthContextKey = contextKey("auth")

func (s *AuthService) GenerateToken(authUser *AuthUser) (string, error) {
	claims := JWTClaims{
		UserID: authUser.ID,
		Email:  authUser.Email,
		Name:   authUser.Name,
	}
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour * 24))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) ParseToken(ctx context.Context, tokenString string) (*AuthContext, error) {
	claims := &JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		slog.ErrorContext(ctx, "invalid or expired token", slog.Any("token", token))
		return nil, err
	}

	return &AuthContext{
		UserID: claims.UserID,
		Email:  claims.Email,
		Name:   claims.Name,
	}, nil
}
