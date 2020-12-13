package app

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mmkader85/bookstore_users-api/datasources/mysql/users_db"
	"github.com/mmkader85/bookstore_users-api/logger"
)

var (
	router = gin.Default()
)

func StartApplication() {
	loadEnv()
	users_db.InitializeDB()
	mapUrls()
	logger.Info("About to start the application...")
	err := router.Run(":8000")
	if err != nil {
		log.Println("Server Error:", err)
	}
}
