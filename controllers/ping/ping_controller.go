package ping

import (
	"github.com/gin-gonic/gin"
	"github.com/mmkader85/bookstore_users-api/datasources/mysql/users_db"
	"net/http"
)

func Ping(ctx *gin.Context) {
	err := users_db.Client.Ping()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.String(http.StatusOK, "pong")
}
