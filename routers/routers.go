package routers

import (
		"context"
		"fmt"
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

func ListingsPages() *chi.Mux  {
		router := chi.NewRouter()
		router.Get("/all-listings", controllers.GetAllListings)
		router.Get("/apartments", controllers.GetApartmentListings)
		router.Get("/homes", controllers.GetHomeListings)
		router.Route("/data", func(r chi.Router) {
			//r.Use(PropertyCtx)
			//localhost/api/pages/data/{propertyID}
			r.With(PropertyCtx).Get("/{propertyID}", controllers.GetPropertyPageData)
			r.With(QueueCtx).Get("/{propertyID}/queue", controllers.GetPropertyQueue)
		})
		return router
}

type contextResponseWriter struct {
		http.ResponseWriter
		ctx context.Context
}

func PropertyCtx(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var propertyData models.PropertyPageData
				var err error
				if propertyID := chi.URLParam(r, "propertyID"); propertyID != "" {
						propertyData, err = models.GetPropertyPageData(propertyID)
						if err != nil {
								panic(err.Error())
								return
						}
				} else {
						fmt.Println("No id ctx")
						return
				}

				ctx := context.WithValue(r.Context(), "propertyID", propertyData)
				next.ServeHTTP(w, r.WithContext(ctx))
		})
}

func QueueCtx(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var qDataArr []models.PropertyQueueData
				var err error
				if propertyID := chi.URLParam(r, "propertyID"); propertyID != "" {
						qDataArr, err = models.GetProperyQueue(propertyID)
						if err != nil {
								panic(err.Error())
								return
						}
				} else {
						fmt.Println("No id ctx")
						return
				}

				ctx := context.WithValue(r.Context(), "propertyQueue", qDataArr)
				next.ServeHTTP(w, r.WithContext(ctx))
		})
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
//
//func (s *server) adminOnly(h http.HandlerFunc) http.HandlerFunc {
//		return func(w http.ResponseWriter, r *http.Request) {
//				if !currentUser(r).IsAdmin {
//						http.NotFound(w, r)
//						return
//				}
//				h(w, r)
//		}
//}





