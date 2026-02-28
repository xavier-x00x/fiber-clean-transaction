package jwtutil

import (
	"errors"
	"fiber-clean-transaction/internal/dto"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// var secret = []byte("secret")

// func SetSecretKey(key []byte) {
//     secret = key
// }

func GenerateJWT(user *dto.UserJwt, minute int64) (string, error) {
	var secret = []byte(os.Getenv("JWT_SECRET"))
	claims := jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.Role,
		"store": user.Store,
		"exp":   time.Now().Add(time.Duration(minute) * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func ValidateJWT(tokenString string) (*dto.UserJwt, error) {
	var secret = []byte(os.Getenv("JWT_SECRET"))
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("token expired")
		}
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["id"] == nil {
		return nil, errors.New("invalid token claims")
	}

	auth := &dto.UserJwt{
		ID: uint(claims["id"].(float64)),
	}
	if email, ok := claims["email"].(string); ok {
		auth.Email = email
	}
	if role, ok := claims["role"].(string); ok {
		auth.Role = role
	}
	if store, ok := claims["store"].(string); ok {
		auth.Store = store
	}
	return auth, nil
}
