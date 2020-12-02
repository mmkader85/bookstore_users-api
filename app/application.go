package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

func StartApplication() {
	mapUrls()
	err := router.Run(":8000")
	if err != nil {
		fmt.Println("Server Error:", err)
	}
}
