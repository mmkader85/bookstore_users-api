package users_db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var Client *sql.DB

func InitializeDB() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
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

	if err = Client.Ping(); err != nil {
		panic(err)
	}
	log.Println("Database connection successful.")
}
