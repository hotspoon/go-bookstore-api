package health

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, db *sql.DB) {
	handler := NewHandler(db)
	router.GET("/health", handler.Check)
}
