package app

import (
	"github.com/mmkader85/bookstore_users-api/controllers/ping"
	"github.com/mmkader85/bookstore_users-api/controllers/users"
)

func mapUrls() {
	router.GET("/ping", ping.Ping)
	router.POST("/user", users.CreateUser)
	router.GET("/user/:user_id", users.GetUser)
}
