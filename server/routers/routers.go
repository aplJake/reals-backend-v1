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
	router.Get("/countries", controllers.GetUsers)
	return router
}

func UserProfile() *chi.Mux {
	router := chi.NewRouter()
	router.With(middleware.UserProfileCtx).Route("/{userId}", func(r chi.Router) {
		r.Get("/", controllers.GetProfile)
		r.Put("/", controllers.UpdateProfile)
		r.Post("/", controllers.AddAds)
		// Add bew property from user profile
		r.Post("/property/new", controllers.NewPropertyListing)
		r.Get("/property/new", controllers.GetSeller)
	})

	// Listing manipulations
	router.Route("/property", func(r chi.Router) {
		//r.Use(PropertyCtx);
		r.With(PropertyUpdateCtx).Get("/update/{propertyID}", controllers.GetPropertyListingUpdate)
		r.Put("/update", controllers.PropertyListingUpdate)
		r.With(PropertyUpdateCtx).Delete("/delete/{propertyID}", controllers.PropertiesListingDelete)
	})

	return router
}

func PropertyUpdateCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		propertyID := chi.URLParam(r, "propertyID")
		if propertyID == "" {
			fmt.Print("porpertyID is emty")
			return
		}

		ctx := context.WithValue(r.Context(), "propertyID", propertyID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ListingsPages() *chi.Mux {
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
			fmt.Println("No id ctx", propertyID)
			return
		}

		ctx := context.WithValue(r.Context(), "propertyID", propertyData)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func QueueCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var propCtxData models.PropertyCtxData
		var err error
		if propertyID := chi.URLParam(r, "propertyID"); propertyID != "" {
			propCtxData.Queue, err = models.GetProperyQueue(propertyID)
			if err != nil {
				panic(err.Error())
				return
			}
			propCtxData.Profile, err = models.GetPropertyProfileData(propertyID)
		} else {
			fmt.Println("No id ctx")
			return
		}

		ctx := context.WithValue(r.Context(), "propertyData", propCtxData)
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
	router.Get("/with-cities", controllers.GetCountriesWithCities)
	router.Post("/", controllers.AddNewCountry)
	router.Put("/", controllers.UpdateCountry)
	router.With(CountryDeleteCtx).Delete("/{countryID}", controllers.DeleteCountry)

	// Cities Country
	router.With(CountryCtx).Get("/{countryID}/cities", controllers.GetCitiesByCountry)
	return router
}

func CountryDeleteCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		countryID := chi.URLParam(r, "countryID")
		if countryID == "" {
			fmt.Print("countryid is emty")
			return
		}

		ctx := context.WithValue(r.Context(), "countryToDeleteID", countryID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
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

// Cities and CityRegion router
func CitiesAnonymousHandler() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/", controllers.GetCitiesList)
	router.Post("/", controllers.AddNewCity)
	router.Put("/", controllers.UpdateCity)
	router.With(CityIDCtx).Delete("/{cityID}", controllers.DeleteCity)

	router.With(CityIDCtx).Get("/{cityID}/regions", controllers.GetRegionsList)

	//router.Route("/{cityID}", func(r chi.Router) {
	//	r.Use(CountryCtx)
	//	r.Get("/cities", controllers.GetCitiesList)
	//})
	return router
}

func CityIDCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cityID := chi.URLParam(r, "cityID")
		if cityID == "" {
			fmt.Print("city is emty")
			return
		}

		ctx := context.WithValue(r.Context(), "cityToDeleteID", cityID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

//localhost:/api/regions/
func CityRegionHandler() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/", controllers.GetRegionsList)
	router.Post("/", controllers.CreateNewRegion)
	router.Put("/", controllers.UpdateRegion)
	router.With(CityRegionIDCtx).Delete("/{regionID}", controllers.DeleteRegionByID)
	return router
}

func CityRegionIDCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		regionID := chi.URLParam(r, "regionID")
		if regionID == "" {
			fmt.Print("region id is empty")
			return
		}

		ctx := context.WithValue(r.Context(), "regionID", regionID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func PropertyQueueHandler() *chi.Mux  {
	router := chi.NewRouter()
	router.With(QueuePropertyCtx).Delete("/{userID}/property/{propertyID}", controllers.DeleteQueueUser)
	router.Post("/", controllers.AddUserToQueue)
	return router
}

func QueuePropertyCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var queue models.Queue
		var propertyID, userID string

		if propertyID = chi.URLParam(r, "propertyID"); propertyID != "" {
			queue.PropertyID = propertyID
		}

		if userID = chi.URLParam(r, "userID"); userID != "" {
			queue.UserID = userID
		}

		ctx := context.WithValue(r.Context(), "queueData", queue)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AdminPageHandler() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.AdminOnly)
	// Users info
	router.Get("/users", controllers.GetUsers)
	router.Post("/users", controllers.CreateNewAdminUser)
	// Admins info
	router.Get("/admins", controllers.GetAdmins)
	router.With(AdminDeleteCtx).Delete("/admins/{adminID}", controllers.DeleteAdminUser)
	return router
}

func AdminDeleteCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		adminID := chi.URLParam(r, "adminID")
		if adminID == "" {
			return
		}

		ctx := context.WithValue(r.Context(), "adminToDeleteID", adminID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func NotificationHandler() *chi.Mux {
	router := chi.NewRouter()
	// GET notifications in user profile
	router.With(NotificationUserIDCtx).Get("/{userID}", controllers.GetNotifications)
	router.With(NotificationCtx).Delete("/{notificationID}", controllers.DeleteNotification)
	return router
}

func NotificationUserIDCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userID")
		if userID == "" {
			fmt.Print("user id is emty")
			return
		}

		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func NotificationCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		notificationID := chi.URLParam(r, "notificationID")
		if notificationID == "" {
			fmt.Print("notification id is emty")
			return
		}

		ctx := context.WithValue(r.Context(), "notificationID", notificationID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}