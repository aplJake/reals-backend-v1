package main

import (
		"fmt"
		"github.com/aplJake/reals-course/server/controllers"
		"github.com/aplJake/reals-course/server/routers"
		"github.com/aplJake/reals-course/server/routers/middleware"
		"github.com/go-chi/chi"
		m "github.com/go-chi/chi/middleware"
		"log"
		"net/http"
)

func InitRouter() *chi.Mux {
		router := chi.NewRouter()
		// TODO: ADD MIDDLEWARE FOR JWT AUTHORIZATION
		router.Use(
				m.Logger,
				m.DefaultCompress,
				middleware.CorsMiddleware.Handler,
		)

		router.Route("/api", func(r chi.Router) {
				r.Post("/signup", controllers.UserSignUp)
				r.Get("/signup", controllers.GetUser)
				r.Post("/signin", controllers.UserSignIn)

				//r.Mount("/", routers.UserAuthentication())
				r.Mount("/{userId}", routers.UserProfile())
				// User addding
				//r.Mount("/{userId}", routers.PropertyAdding())
		})

		// Public routes
		router.Group(func(r chi.Router) {
				r.Get("/", func(w http.ResponseWriter, r *http.Request) {
						w.Write([]byte("Welcome to API"))
				})
		})
		return router
}

func main() {
		fmt.Println("Server Running...")
		router := InitRouter()
		log.Fatal(http.ListenAndServe(":2308", router))
}
