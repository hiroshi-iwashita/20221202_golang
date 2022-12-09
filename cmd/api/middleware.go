package main

import (
	"net/http"

	"github.com/justinas/nosurf"
)

// NoSurf is the csrf protection middleware
func (app *applicationConfig) noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.inProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

// SessionLoad loads and saves session data for current request
// func SessionLoad(next http.Handler) http.Handler {
// 	return session.LoadAndSave(next)
// }

// func (app *application) AuthTokenMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		_, err := app.models.Token.AuthenticateToken(r)
// 		if err != nil {
// 			payload := jsonResponse{
// 				Error:   true,
// 				Message: "invalid authentication credentials",
// 			}

// 			_ = app.writeJSON(w, http.StatusUnauthorized, payload)
// 			return
// 		}
// 		next.ServeHTTP(w, r)
// 	})
// }
