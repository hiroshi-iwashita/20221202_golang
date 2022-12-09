package main

import (
	"io"
	"net/http"

	"github.com/hiroshi-iwashita/20221202_golang/internal/models"
)

// jsonResponse is the type used for generic JSON responses
type jsonResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type envelope map[string]interface{}

func (app *applicationConfig) User(w http.ResponseWriter, _ *http.Request) {
	io.WriteString(w, "Users")
}

func (app *applicationConfig) AllUsers(w http.ResponseWriter, r *http.Request) {
	var users models.User
	all, err := users.GetAll()
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "success",
		Data:    envelope{"users": all},
	}

	app.writeJSON(w, http.StatusOK, payload)
}
