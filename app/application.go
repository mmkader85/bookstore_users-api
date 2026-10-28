package app

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mmkader85/bookstore_users-api/datasources/mysql/users_db"
	"github.com/mmkader85/bookstore_users-api/logger"
)

const maxRequestBodyBytes = 1 << 20 // 1 MiB

var (
	router = gin.Default()
)

func limitBody(ctx *gin.Context) {
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, maxRequestBodyBytes)
	ctx.Next()
}

func StartApplication() {
	loadEnv()
	users_db.InitializeDB()
	router.Use(limitBody)
	mapUrls()
	logger.Info("About to start the application...")

	// Explicit timeouts protect against slowloris-style connection exhaustion.
	srv := &http.Server{
		Addr:              ":8000",
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalln("Server Error:", err)
	}
}
