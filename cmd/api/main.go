// @title Bookstore API
// @version 1.0
// @description RESTful API for the bookstore application.
// @BasePath /api/v1
package main

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "bookstore-api/docs"
	"bookstore-api/internal/book"
	"bookstore-api/internal/config"
	"bookstore-api/internal/health"
	"bookstore-api/internal/platform/database"
	"bookstore-api/internal/platform/logging"
	appmiddleware "bookstore-api/internal/shared/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	cfg := config.Load()

	logFiles, err := logging.Init(cfg.LogPath, cfg.AccessLogPath)
	if err != nil {
		panic(err)
	}
	defer logFiles.Close()

	db, err := database.Open(cfg.DatabasePath)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	router := setupRouter(cfg, db, logFiles.AccessWriter)

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		logging.Logger.Info().Str("address", server.Addr).Msg("server started")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logging.Logger.Fatal().Err(err).Msg("server failed")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logging.Logger.Info().Msg("shutting down server")
	if err := server.Shutdown(ctx); err != nil {
		logging.Logger.Error().Err(err).Msg("graceful shutdown failed")
	}
}

func setupRouter(cfg config.Config, db *sql.DB, accessWriter io.Writer) *gin.Engine {
	router := gin.New()
	router.Use(gin.LoggerWithWriter(accessWriter))
	router.Use(gin.RecoveryWithWriter(accessWriter))
	router.Use(appmiddleware.RequestID())
	router.Use(appmiddleware.ErrorHandler(book.HTTPErrorMapper))
	router.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.AllowedOrigins,
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Request-ID"},
		ExposeHeaders:    []string{"X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	if err := router.SetTrustedProxies(nil); err != nil {
		panic(err)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api/" + cfg.APIVersion)
	health.RegisterRoutes(api, db)
	book.RegisterRoutes(api, db)

	return router
}
