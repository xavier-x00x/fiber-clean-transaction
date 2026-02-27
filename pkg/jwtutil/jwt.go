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
		"id":   user.ID,
		"role": user.Role,
		"exp":  time.Now().Add(time.Duration(minute) * time.Minute).Unix(), // Token expires in minutes
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
		return nil, err
	}

	auth := &dto.UserJwt{
		ID:   uint(claims["id"].(float64)),
		Role: claims["role"].(string),
	}
	return auth, nil
}
