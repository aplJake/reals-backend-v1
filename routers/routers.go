package routers

import (
		"context"
		"github.com/aplJake/reals-course/server/controllers"
		"github.com/aplJake/reals-course/server/models"
		"github.com/aplJake/reals-course/server/routers/middleware"
		"github.com/go-chi/chi"
		"net/http"
)

func UserAuthentication() *chi.Mux {
		router := chi.NewRouter()
		router.Post("/signup", controllers.UserSignUp)
		router.Get("/signup", controllers.GetUser)
		router.Post("/signin", controllers.UserSignIn)
		return router
}

func Users() *chi.Mux {
		router := chi.NewRouter()
		router.Get("/listings", controllers.GetUsers)
		return router
}



func UserProfile() *chi.Mux {
		router := chi.NewRouter()
		router.Use(middleware.UserProfileCtx)
		router.Get("/", controllers.GetProfile)
		router.Put("/", controllers.UpdateProfile)
		router.Post("/", controllers.AddAds)
		// Add bew property from user profile
		router.Post("/property/new", controllers.NewPropertyListing)
		router.Get("/property/new", controllers.GetSeller)

		return router
}

//func PropertyAdding() *chi.Mux {
//		router := chi.NewRouter()
//		router.Use(middleware.UserProfileCtx)
//		router.Post("/property/new", controllers.NewPropertyListing)
//		return router
//}

func CountriesAnonymousHandler() *chi.Mux {
		router := chi.NewRouter()
		router.Get("/", controllers.GetCountries)
		
		router.Route("/{countryID}", func(r chi.Router) {
				r.Use(CountryCtx)
				r.Get("/cities", controllers.GetCitiesByCountry)
		})
		return router
}

func CountryCtx(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var country models.Country

				if countryID := chi.URLParam(r, "countryID"); countryID != "" {
						country = models.DBGetCountry(countryID)
				} else {
						return
				}

				ctx := context.WithValue(r.Context(), "coutryID", country)
				next.ServeHTTP(w, r.WithContext(ctx))
		})
}

func AdminPageHandler() *chi.Mux {
		router := chi.NewRouter()
		router.Use(middleware.AdminOnly)
		router.Get("/listings", controllers.GetUsers)
		router.Post("/listings", controllers.CreateNewAdminUser)
		return router
}




