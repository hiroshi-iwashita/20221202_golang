package driver

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
)

type DB struct {
	SQL *sql.DB
	dns string
}

var dbConn = &DB{}

const maxOpenDbConn = 5 // might be changed in production
const maxIdleDbConn = 5 // might be changed in production
const maxDbLifeTime = 5 * time.Minute

func ConnectDB() (*DB, error) {
	setDns()

	d, err := open(dbConn.dns)
	if err != nil {
		return nil, err
	}

	d.SetMaxOpenConns(maxOpenDbConn)
	d.SetMaxIdleConns(maxIdleDbConn)
	d.SetConnMaxLifetime(maxDbLifeTime)

	err = testDB(d)
	if err != nil {
		return nil, err
	}

	dbConn.SQL = d
	return dbConn, nil
}

// DNS
func setDns() {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		// handle error
	}
	address := os.Getenv("MYSQL_HOST") + ":" + os.Getenv("DB_PORT")
	c := mysql.Config{
		DBName:    os.Getenv("MYSQL_DATABASE"),
		User:      os.Getenv("MYSQL_USER"),
		Passwd:    os.Getenv("MYSQL_PASSWORD"),
		Addr:      address,
		Net:       "tcp",
		ParseTime: true,
		Collation: "utf8mb4_unicode_ci",
		Loc:       jst,
	}

	dbConn.dns = c.FormatDSN()
}

// NewDatabase creates a new database for the application
func open(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func testDB(d *sql.DB) error {
	err := d.Ping()
	if err != nil {
		fmt.Println("Error!", err)
		return err
	}
	return nil
}
