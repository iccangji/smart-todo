package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(user User) (string, error) {
	expire := time.Now().Add(60 * time.Minute) // Expiration 60 minutes

	claims := JWTClaims{
		UserID: user.ID.Hex(),
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expire),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(
		[]byte(os.Getenv("JWT_ACCESS_SECRET")),
	)
}

func GenerateRefreshToken(user User) (string, error) {
	expire := time.Now().Add(7 * 24 * time.Hour) // Expiration 7 days

	claims := JWTClaims{
		UserID: user.ID.Hex(),
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expire),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(
		[]byte(os.Getenv("JWT_REFRESH_SECRET")),
	)
}

func ParseRefreshToken(tokenStr string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&JWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_REFRESH_SECRET")), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)

	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
