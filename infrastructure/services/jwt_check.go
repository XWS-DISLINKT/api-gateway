package services

import (
	"github.com/dgrijalva/jwt-go"
	"net/http"
)

var jwtKey = []byte("secret_key")
var loggedUser = ""

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func JWTValid(w http.ResponseWriter, r *http.Request) bool {
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return false
		}
		w.WriteHeader(http.StatusBadRequest)
		return false
	}

	tknStr := c.Value

	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return false
		}
		w.WriteHeader(http.StatusBadRequest)
		return false
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return false
	}
	loggedUser = claims.Username
	return true
}
