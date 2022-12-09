package main

import (
	"io"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
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
	mux.Get("/about", func(w http.ResponseWriter, _ *http.Request) {
		io.WriteString(w, "about page")
	})
	mux.Get("/contact", func(w http.ResponseWriter, _ *http.Request) {
		io.WriteString(w, "contact page")
	})
	mux.Get("/user", app.User)
	mux.Get("/users", app.AllUsers)

	return mux
}
