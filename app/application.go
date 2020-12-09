package app

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mmkader85/bookstore_users-api/datasources/mysql/users_db"
)

var (
	router = gin.Default()
)

func StartApplication() {
	loadEnv()
	users_db.InitializeApp()
	mapUrls()
	err := router.Run(":8000")
	if err != nil {
		fmt.Println("Server Error:", err)
	}
}
