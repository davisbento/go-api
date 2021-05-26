package jwtManager

import (
	"net/http"
	"strings"

	"github.com/codegangsta/negroni"
	"github.com/davisbento/go-api/core/utils"
	jwt "github.com/dgrijalva/jwt-go"
)

func GenerateToken(userId int64) (string, error) {
	atClaims := jwt.MapClaims{}
	atClaims["userId"] = userId

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte("secret"))

	return token, err
}

func JwtAuth() negroni.Handler {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		if r.URL.Path == "/v1/user" || r.URL.Path == "/v1/users/login" {
			next(w, r)
			return
		}

		bearToken := r.Header.Get("Authorization")
		splitToken := strings.Split(bearToken, " ")

		if len(splitToken) < 2 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(utils.FormatJSONError("Token missing"))
			return
		}

		token := splitToken[1]

		_, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(utils.FormatJSONError("Token Inválido"))
			return
		}

		next(w, r)
	})
}
