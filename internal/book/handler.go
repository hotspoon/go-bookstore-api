package book

import (
	"net/http"

	appmiddleware "bookstore-api/internal/shared/middleware"
	"bookstore-api/internal/shared/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return Handler{service: service}
}

// FindAll godoc
// @Summary Get all books
// @Description Returns all books ordered by ID.
// @Tags books
// @Produce json
// @Success 200 {object} response.Envelope{data=[]Book}
// @Failure 500 {object} response.Error
// @Router /books [get]
func (h Handler) FindAll(c *gin.Context) {
	books, err := h.service.FindAll(c.Request.Context())
	if err != nil {
		appmiddleware.AddError(c, err, "failed to get books")
		return
	}

	response.OK(c, http.StatusOK, books)
}

// FindOne godoc
// @Summary Get a book by ID
// @Description Returns a single book by its ID.
// @Tags books
// @Produce json
// @Param id path string true "Book ID"
// @Success 200 {object} response.Envelope{data=Book}
// @Failure 400 {object} response.Error
// @Failure 404 {object} response.Error
// @Failure 500 {object} response.Error
// @Router /books/{id} [get]
func (h Handler) FindOne(c *gin.Context) {
	bookID := c.Param("id")
	book, err := h.service.FindOne(c.Request.Context(), bookID)
	if err != nil {
		appmiddleware.AddError(c, err, "failed to get book")
		return
	}

	response.OK(c, http.StatusOK, book)
}

// Create godoc
// @Summary Create a new book
// @Description Creates a new book with a generated ID.
// @Tags books
// @Accept json
// @Produce json
// @Param book body CreateRequest true "Book data"
// @Success 201 {object} response.Mutation
// @Failure 400 {object} response.Error
// @Failure 409 {object} response.Error
// @Failure 500 {object} response.Error
// @Router /books [post]
func (h Handler) Create(c *gin.Context) {
	var request CreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid book data")
		return
	}

	book, err := h.service.Create(c.Request.Context(), request.Book())
	if err != nil {
		appmiddleware.AddError(c, err, "failed to create book")
		return
	}

	response.Success(c, http.StatusCreated, "book created successfully", book.ID)
}

// Update godoc
// @Summary Update an existing book
// @Description Updates an existing book by its ID with the provided data.
// @Tags books
// @Accept json
// @Produce json
// @Param id path string true "Book ID"
// @Param book body UpdateRequest true "Updated book data"
// @Success 200 {object} response.Mutation
// @Failure 400 {object} response.Error
// @Failure 404 {object} response.Error
// @Failure 409 {object} response.Error
// @Failure 500 {object} response.Error
// @Router /books/{id} [put]
func (h Handler) Update(c *gin.Context) {
	var request UpdateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid book data")
		return
	}

	if err := h.service.Update(c.Request.Context(), c.Param("id"), request.Book()); err != nil {
		appmiddleware.AddError(c, err, "failed to update book")
		return
	}

	response.Success(c, http.StatusOK, "book updated successfully", c.Param("id"))
}

// Delete godoc
// @Summary Delete a book by ID
// @Description Deletes a book by its ID.
// @Tags books
// @Produce json
// @Param id path string true "Book ID"
// @Success 200 {object} response.Mutation
// @Failure 400 {object} response.Error
// @Failure 404 {object} response.Error
// @Failure 500 {object} response.Error
// @Router /books/{id} [delete]
func (h Handler) Delete(c *gin.Context) {
	if err := h.service.Delete(c.Request.Context(), c.Param("id")); err != nil {
		appmiddleware.AddError(c, err, "failed to delete book")
		return
	}

	response.Success(c, http.StatusOK, "book deleted successfully", c.Param("id"))
}
