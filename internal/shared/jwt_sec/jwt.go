package jwtsec

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	srverr "protravel-finance/internal/shared/server_error"
	"time"
)

type JWT struct {
	jwtSecret []byte
}

func NewJWT(jwtSecret string) *JWT {
	return &JWT{
		jwtSecret: []byte(jwtSecret),
	}
}

func (m *JWT) ValidateAccessToken(tokenString string) (string, srverr.ServerError) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return m.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return "", srverr.NewServerError(srverr.ErrUnauthorized)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", srverr.NewServerError(srverr.ErrInternalServerError).
			SetError("invalid claims")
	}
	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", srverr.NewServerError(srverr.ErrInternalServerError).
			SetError("user_id not found in claims")
	}
	return userID, nil
}

func (m *JWT) GenerateAccessToken(userID string) (string, srverr.ServerError) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().UTC().Add(15 * time.Minute).Unix(),
		"iat":     time.Now().UTC().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	access, err := token.SignedString(m.jwtSecret)
	if err != nil {
		return "", srverr.NewServerError(srverr.ErrInternalServerError).
			SetError(err.Error()).SetDetails("failed to sign access token")
	}
	return access, nil
}

func (m *JWT) GenerateRefreshToken(userID string) (string, string, srverr.ServerError) {
	jti := uuid.New().String() // jti для ревок
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().UTC().Add(7 * 24 * time.Hour).Unix(),
		"iat":     time.Now().UTC().Unix(),
		"jti":     jti,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refresh, err := token.SignedString(m.jwtSecret)
	if err != nil {
		return "", "", srverr.NewServerError(srverr.ErrInternalServerError).
			SetError(err.Error()).SetDetails("failed to sign refresh token")
	}
	return refresh, jti, nil
}

// ValidateRefreshToken return userID and jti
func (m *JWT) ValidateRefreshToken(tokenString string) (string, string, srverr.ServerError) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return m.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return "", "", srverr.NewServerError(srverr.ErrUnauthorized)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", srverr.NewServerError(srverr.ErrInternalServerError).
			SetError("invalid claims")
	}
	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", "", srverr.NewServerError(srverr.ErrInternalServerError).
			SetError("user_id not found in claims")
	}
	jti, ok := claims["jti"].(string)
	if !ok {
		return "", "", srverr.NewServerError(srverr.ErrInternalServerError).
			SetError("jti not found in claims")
	}
	return userID, jti, nil
}
