package users_db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var Client *sql.DB

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("unable to load environment config file")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
		os.Getenv("USERS_DB_USERNAME"),
		os.Getenv("USERS_DB_PASSWORD"),
		os.Getenv("USERS_DB_HOST"),
		os.Getenv("USERS_DB_NAME"))

	Client, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	if err = Client.Ping(); err != nil {
		panic(err)
	}
	log.Println("database connection successful")
}
