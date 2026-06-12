package author

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, db *sql.DB) {
	repository := NewRepository(db)
	service := NewService(repository)
	handler := NewHandler(service)

	router.GET("/authors", handler.FindAll)
	router.GET("/authors/:id", handler.FindOne)
	router.POST("/authors", handler.Create)
	router.PUT("/authors/:id", handler.Update)
	router.DELETE("/authors/:id", handler.Delete)
}
