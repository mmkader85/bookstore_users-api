package ping

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mmkader85/bookstore_users-api/datasources/mysql/users_db"
	"github.com/mmkader85/bookstore_users-api/logger"
)

func Ping(ctx *gin.Context) {
	err := users_db.Client.Ping()
	if err != nil {
		logger.Error("Database ping failed.", err)
		ctx.String(http.StatusInternalServerError, "database unavailable")
		return
	}

	ctx.String(http.StatusOK, "pong")
}
