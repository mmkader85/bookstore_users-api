package app

import (
	"github.com/mmkader85/bookstore_users-api/controllers/ping"
	"github.com/mmkader85/bookstore_users-api/controllers/users"
)

func mapUrls() {
	router.GET("/ping", ping.Ping)
	router.POST("/user", users.Create)
	router.GET("/user/:user_id", users.Get)
	router.PUT("/user/:user_id", users.Update)
	router.PATCH("/user/:user_id", users.Update)
	router.DELETE("/user/:user_id", users.Delete)
	router.GET("/internal/user/search", users.Search)
}
