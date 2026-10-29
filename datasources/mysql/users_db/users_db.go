package users_db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var Client *sql.DB

func InitializeDB() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&interpolateParams=true&timeout=5s&readTimeout=10s&writeTimeout=10s",
		os.Getenv("USERS_DB_USERNAME"),
		os.Getenv("USERS_DB_PASSWORD"),
		os.Getenv("USERS_DB_HOST"),
		os.Getenv("USERS_DB_PORT"),
		os.Getenv("USERS_DB_NAME"))

	var err error
	Client, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	// Bound the pool so a traffic spike can't exhaust MySQL connections,
	// and recycle connections before the server/proxy drops them.
	Client.SetMaxOpenConns(25)
	Client.SetMaxIdleConns(25)
	Client.SetConnMaxLifetime(5 * time.Minute)

	if err = Client.Ping(); err != nil {
		panic(err)
	}
	log.Println("Database connection successful.")
}
