package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/hiroshi-iwashita/20221202_golang/internal/driver"
	"github.com/hiroshi-iwashita/20221202_golang/internal/models"
	"github.com/jmoiron/sqlx"
)

// // application is the type for all data we want to share with the
// // various parts of our application. We will share this information
// // in most cases by using this type as the receiver for functions.
type applicationConfig struct {
	port         int
	infoLog      *log.Logger
	errorLog     *log.Logger
	models       models.Models
	environment  string
	inProduction bool
}

var port int
var dbConnectRetryTimes int
var infoLog *log.Logger
var errorLog *log.Logger
var environment string
var inProduction bool

func init() {
	// fmt.Println("main.init")

	// set port number
	eap := os.Getenv("API_PORT")
	p, _ := strconv.Atoi(eap)
	port = p

	// set DB connect retry times
	dbcrt := os.Getenv("DB_CONNECT_RETRY")
	dc, _ := strconv.Atoi(dbcrt)
	dbConnectRetryTimes = dc

	// setup infoLog
	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	// setup errorLog
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// set environment
	environment = os.Getenv("ENV")

	// set inProduction
	eip := os.Getenv("INPRODUCTION")
	ip, _ := strconv.ParseBool(eip)
	inProduction = ip
}

func main() {
	dbPool, _ := runDB()

	app := &applicationConfig{
		port:         port,
		infoLog:      infoLog,
		errorLog:     errorLog,
		models:       models.New(dbPool),
		environment:  environment,
		inProduction: inProduction,
	}

	err := app.serveAPIPort()
	if err != nil {
		log.Fatal(err)
	}
}

// runDB connects to database
func runDB() (*sqlx.DB, error) {
	db, err := driver.ConnectDB(dbConnectRetryTimes)
	if err != nil {
		log.Fatal("Cannot connect to database")
	}

	return db.SQL, nil
}

// serveAPIPort starts the API server
func (app *applicationConfig) serveAPIPort() error {
	app.infoLog.Println("API listening on port", app.port)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.port),
		Handler: app.routes(),
	}

	return srv.ListenAndServe()
}
