package utils

import (
	"errors"
	"time"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserId int    `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken generates access and refresh tokens
func GenerateToken(userId int, email, role string) (string, string, error) {
	jwtSecret, _ := beego.AppConfig.String("jwtsecret")
	jwtExpiration, _ := beego.AppConfig.Int("jwtexpiration")

	// Access token
	accessClaims := Claims{
		UserId: userId,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(jwtExpiration) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", "", err
	}

	// Refresh token (valid for 7 days)
	refreshClaims := Claims{
		UserId: userId,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

// ParseToken parses and validates a JWT token
func ParseToken(tokenString string) (*Claims, error) {
	jwtSecret, _ := beego.AppConfig.String("jwtsecret")

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
