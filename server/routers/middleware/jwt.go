package middleware

import (
	"context"
	"fmt"
	"github.com/aplJake/reals-course/server/models"
	"github.com/aplJake/reals-course/server/utils"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"strings"
)

var JwtAuthentication = func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Endpoints that doesnt require auth
				notAuth := []string{"/api/user/new", "/api/user/login"}
				// Get current request path
				requestPath := r.URL.Path

				// Check if request doesnt need authentication, serve the request if it
				// doesnt need it
				for _, value := range notAuth {
						if value == requestPath {
								next.ServeHTTP(w, r)
								return
						}
				}

				response := make(map[string]interface{})
				// Get the token from the header
				tokenHeader := r.Header.Get("Authorization")

				// If token is missing, return 403 Unauth
				if tokenHeader == "" {
						response = utils.Message(false, "Missing auth token")
						w.WriteHeader(http.StatusForbidden)
						w.Header().Add("Content-Type", "application/json")
				}

				// Check the token in the Header with retrieved token
				// Form: `Bearer {token-body}`
				splitted := strings.Split(tokenHeader, "")
				if len(splitted) != 2 {
						fmt.Println(tokenHeader)
						response = utils.Message(false, "Invalid/Malformed auth token")
						w.WriteHeader(http.StatusForbidden)
						w.Header().Add("Content-Type", "application/json")
						utils.Respond(w, response)
						return
				}

				// Grab the token part
				tokenPart := splitted[1]
				tk := models.Token{}

				token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
						return []byte(os.Getenv("token_password")), nil
				})

				// Malformed token, returns 403
				if err != nil {
						response = utils.Message(false, "Malformed authentication token")
						w.WriteHeader(http.StatusForbidden)
						w.Header().Add("Content-Type", "application/json")
						utils.Respond(w, response)
						return
				}

				// Token is invalid, may be not signed on this server
				if !token.Valid {
						response = utils.Message(false, "Token is not valid.")
						w.WriteHeader(http.StatusForbidden)
						w.Header().Add("Content-Type", "application/json")
						utils.Respond(w, response)
						return
				}

				// OK, Set the caller to the user retrieved from the parsed token

				// Uses for monitoring
				//fmt.Sprintf("User %", tk.Username)
				ctx := context.WithValue(r.Context(), "user", tk.UserId)
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
		})
}