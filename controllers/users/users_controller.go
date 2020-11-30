package users

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateUser(ctx *gin.Context) {
	ctx.String(http.StatusNotImplemented, "Implement CreateUser!")
}

func GetUser(ctx *gin.Context) {
	ctx.String(http.StatusNotImplemented, "Implement GetUser!")
}