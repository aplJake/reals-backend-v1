package middleware

import (
		"context"
		"github.com/aplJake/reals-course/server/models"
		"github.com/aplJake/reals-course/server/utils"
		"github.com/go-chi/chi"
		"net/http"
		"strconv"
)

// UserProfile middleware used to load the UserProfile
// from URL userId param as a request
func UserProfileCtx(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var profile *models.UserProfileRespond

				userIdString := chi.URLParam(r, "userId")
				if userIdString != "" {
						// Decode userId
						userId, err := strconv.Atoi(userIdString)
						if err != nil {
								response := utils.Message(false, "Invalid Url id param.")
								w.WriteHeader(http.StatusForbidden)
								w.Header().Add("Content-Type", "application/json")
								utils.Respond(w, response)
								return
						}
						profile, err = models.GetUserProfile(uint(userId))

				} else {
						response := utils.Message(false, "Unknown Error.")
						w.WriteHeader(http.StatusForbidden)
						w.Header().Add("Content-Type", "application/json")
						utils.Respond(w, response)
						return
				}

				ctx := context.WithValue(r.Context(), "profile", profile)
				next.ServeHTTP(w, r.WithContext(ctx))
		})
}

// Admin middleware used to verify the userId in
// admins database table
func AdminOnly(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var admin *models.Admin
				//isAdmin, ok := r.Context().Value("")
				q := chi.URLParam(r, "userId")
				if q != "" {
						userId, err := strconv.Atoi(q)
						if err != nil {
								response := utils.Message(false, "Invalid Url id param for admin user.")
								w.WriteHeader(http.StatusForbidden)
								w.Header().Add("Content-Type", "application/json")
								utils.Respond(w, response)
								return
						}
						// Search the Admin user in the database
						admin, _ := models.GetAdmin(uint(userId))
						if admin == nil {
								response := utils.Message(false, "Such admin user doesn`t exists.")
								w.WriteHeader(http.StatusForbidden)
								w.Header().Add("Content-Type", "application/json")
								utils.Respond(w, response)
								return
						}
				} else {
						response := utils.Message(false, "Unknown Error.")
						w.WriteHeader(http.StatusForbidden)
						w.Header().Add("Content-Type", "application/json")
						utils.Respond(w, response)
						return
				}

				ctx := context.WithValue(r.Context(), "admin", admin)
				next.ServeHTTP(w, r.WithContext(ctx))
		})
}