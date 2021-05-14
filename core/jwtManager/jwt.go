package jwtManager

import (
	jwt "github.com/dgrijalva/jwt-go"
)

func GenerateToken(userId int64) (string, error) {
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["userId"] = userId

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte("secret"))

	return token, err
}
