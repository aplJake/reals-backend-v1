package routers

import (
		"github.com/aplJake/reals-course/server/controllers"
		"github.com/aplJake/reals-course/server/routers/middleware"
		"github.com/go-chi/chi"
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
		router.Get("/users", controllers.GetUsers)
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
		router.Get("/countries", controllers.GetCountries)
		return router
}

func AdminPageHandler() *chi.Mux {
		router := chi.NewRouter()
		router.Use(middleware.AdminOnly)
		router.Get("/users", controllers.GetUsers)
		router.Post("/users", controllers.CreateNewAdminUser)
		return router
}




