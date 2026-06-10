package book

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, db *sql.DB) {
	repository := NewRepository(db)
	service := NewService(repository)
	handler := NewHandler(service)

	router.GET("/books", handler.FindAll)
	router.GET("/books/:id", handler.FindOne)
	router.POST("/books", handler.Create)
	router.PUT("/books/:id", handler.Update)
	router.DELETE("/books/:id", handler.Delete)
}
