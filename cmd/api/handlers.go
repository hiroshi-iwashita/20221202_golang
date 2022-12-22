package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/hiroshi-iwashita/20221202_golang/internal/models"
)

// jsonResponse is the type used for generic JSON responses
type jsonResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type envelope map[string]interface{}

var users models.User

// Login is the handler used to attempt to log a user into the api
func (app *applicationConfig) Login(w http.ResponseWriter, r *http.Request) {
	type credentials struct {
		UserName string `json:"email"`
		Password string `json:"password"`
	}

	var creds credentials
	var payload jsonResponse

	err := app.readJSON(w, r, &creds)
	if err != nil {
		app.errorLog.Println(err)
		payload.Error = true
		payload.Message = "invalid json supplied, or json missing entirely"
		_ = app.writeJSON(w, http.StatusBadRequest, payload)
	}

	// TODO authenticate
	app.infoLog.Println(creds.UserName, creds.Password)

	// look up the user by email
	user, err := app.models.User.ShowByEmail(creds.UserName)
	if err != nil {
		app.errorJSON(w, errors.New("invalid username / password"))
		return
	}

	// validate the user's password
	validPassword, err := user.PasswordMatches(creds.Password)
	if err != nil || !validPassword {
		app.errorJSON(w, errors.New("invalid username / password "))
	}

	// we have a valid user, so generaet a token
	token, err := app.models.Token.GenerateToken(user.UserID, 24*time.Hour)
	if err != nil {
		app.errorJSON(w, err)
	}

	// make sure user is active
	// if user.Active == 0 {
	// 	app.errorJSON(w, errors.New("user is not active"))
	// 	return
	// }

	fmt.Println(*token)
	fmt.Println(*user)
	// save it to the database
	err = app.models.Token.Insert(*token, *user)
	if err != nil {
		app.errorJSON(w, err)
	}

	// send back a response
	payload = jsonResponse{
		Error:   false,
		Message: "logged in",
		Data:    envelope{"token": token, "user": user},
	}

	err = app.writeJSON(w, http.StatusOK, payload)
	if err != nil {
		app.errorLog.Println(err)
	}
}

func (app *applicationConfig) User(w http.ResponseWriter, _ *http.Request) {
	io.WriteString(w, "user")
}

func (app *applicationConfig) AllUsers(w http.ResponseWriter, r *http.Request) {
	all, err := users.Index()
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "success",
		Data:    all,
	}

	app.writeJSON(w, http.StatusOK, payload)
}

func (app *applicationConfig) getUserByID(w http.ResponseWriter, r *http.Request) {
	var userID = chi.URLParam(r, "id")
	// if err != nil {
	// 	app.errorJSON(w, err)
	// 	return
	// }

	user, err := users.ShowByID(userID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, user)
}

func (app *applicationConfig) DeleteUserByID(w http.ResponseWriter, r *http.Request) {
	// How to user request payload (?)
	// var requestPayload struct {
	// 	ID int `json:"id"`
	// }

	// err := app.readJSON(w, r, &requestPayload)
	// if err != nil {
	// 	app.errorJSON(w, err)
	// 	return
	// }

	userID := chi.URLParam(r, "user_id")

	err := app.models.User.Delete(userID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "User deleted",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}
