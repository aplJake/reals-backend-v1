package models

import (
		"github.com/dgrijalva/jwt-go"
		"net/http"
		"time"
)

type CredentialsToken struct {
		UserId  uint
		IsAdmin bool
		UserType string
		jwt.StandardClaims
}

var jwtKey = []byte("7278AC6970BEC14904CFC8CB7D4F933C")

type TokenI interface {
		CreateToken()
		RefreshToken() error
}

func (token *CredentialsToken) CreateJWTToken(w http.ResponseWriter, expirationTime time.Time) string {
		jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), token)

		tokenString, err := jwtToken.SignedString([]byte(jwtKey))
		if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
		}

		http.SetCookie(w, &http.Cookie{
				Name:    "token",
				Value:   tokenString,
				Expires: expirationTime,
		})

		return tokenString
}
