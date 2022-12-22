package main

import (
	"io"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/hiroshi-iwashita/20221202_golang/internal/models"
)

// routes generates our routes and attaches them to handlers, using the chi router
// note that we return type http.Handler, and not *chi.Mux; since chi.Mux satisfies
// the interface requirements for http.Handler, it makes sense to return the type
// that is part of the standard library.
func (app *applicationConfig) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	// mux.Use(app.noSurf)

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"https://*",
			"http://*",
		},
		AllowedMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
		},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-CSRF-Token",
		},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Get("/", func(w http.ResponseWriter, _ *http.Request) {
		io.WriteString(w, "Hello world")
	})

	mux.Route("/auth", func(mux chi.Router) {
		// mux.Get("/login", app.Login)
		mux.Post("/login", app.Login)
	})

	mux.Get("/user", app.User)
	mux.Get("/users/all", app.AllUsers)
	mux.Get("/users/get/{id}", app.getUserByID)
	mux.Get("/users/add", func(w http.ResponseWriter, r *http.Request) {
		var u = models.User{
			FirstName: "You",
			LastName:  "There",
			Email:     "edsoif@dfj.com",
			Password:  "password",
		}

		app.infoLog.Println("Adding user...")

		userID, err := app.models.User.Insert(u)
		if err != nil {
			app.errorLog.Println(err)
			app.errorJSON(w, err, http.StatusForbidden)
			return
		}

		app.infoLog.Println("Got back user_id of", userID)
		newUser, _ := app.models.User.ShowByID(userID)
		app.writeJSON(w, http.StatusOK, newUser)
	})
	mux.Post("/users/delete/{user_id}", app.DeleteUserByID)

	return mux
}
