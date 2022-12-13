package driver

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	SQL *sqlx.DB
	dns string
}

var dbConn = &DB{}

const maxOpenDbConn = 5 // might be changed in production
const maxIdleDbConn = 5 // might be changed in production
const maxDbLifeTime = 5 * time.Minute

func ConnectDB(count int) (*DB, error) {
	setDns()

	d, err := open(dbConn.dns)
	if err != nil {
		return nil, err
	}

	d.SetMaxOpenConns(maxOpenDbConn)
	d.SetMaxIdleConns(maxIdleDbConn)
	d.SetConnMaxLifetime(maxDbLifeTime)

	err = testDB(d, count)
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
		log.Fatal("Load location failed:", err)
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
func open(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Open(os.Getenv("DB_DRIVER"), dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func testDB(d *sqlx.DB, count int) error {
	err := d.Ping()
	if err != nil {
		if count <= 0 {
			fmt.Println("Error!", err)
			return err
		}
		time.Sleep(time.Second * 2)
		count--
		fmt.Printf("Retry Connection... count:%v\n", count)
		return testDB(d, count)
	}
	return nil

	// use Ping.Context but not works
	// ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	// defer cancel()

	// if err := d.PingContext(ctx); err != nil {
	// 	log.Fatal(err)
	// 	return err
	// }

	// return nil
}

// RoundTrip を DB接続リトライで実装するほうがbetter

// type retryableRoundTripper struct {
// 	base     http.RoundTripper
// 	attempts int
// 	waitTime time.Duration
// }

// func (rt *retryableRoundTripper) shouldRetry(resp *http.Response, err error) bool {
// 	if err != nil {
// 		var netErr net.Error
// 		if errors.As(err, &netErr) && netErr.Temporary() {
// 			return true
// 		}
// 	}

// 	if resp != nil {
// 		if resp.StatusCode == 429 ||
// 			(500 <= resp.StatusCode && resp.StatusCode <= 504) {
// 			return true
// 		}
// 	}

// 	return false
// }

// func (rt *retryableRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
// 	var (
// 		resp *http.Response
// 		err  error
// 	)
// 	for count := 0; count < rt.attempts; count++ {
// 		resp, err = rt.base.RoundTrip(req)

// 		if !rt.shouldRetry(resp, err) {
// 			return resp, err
// 		}

// 		select {
// 		case <-req.Context().Done():
// 			return nil, req.Context().Err()
// 		case <-time.After(rt.waitTime):
// 		}
// 	}
// 	return resp, nil
// }
