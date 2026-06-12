package author

import (
	"net/http"
	"strconv"

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
// @Summary Get all authors
// @Description Returns all authors stored in the Source table.
// @Tags authors
// @Produce json
// @Success 200 {object} response.Envelope{data=[]Author}
// @Failure 500 {object} response.Error
// @Router /authors [get]
func (h Handler) FindAll(c *gin.Context) {
	authors, err := h.service.FindAll(c.Request.Context())
	if err != nil {
		appmiddleware.AddError(c, err, "failed to get authors")
		return
	}
	response.OK(c, http.StatusOK, authors)
}

// FindOne godoc
// @Summary Get an author by ID
// @Tags authors
// @Produce json
// @Param id path int true "Author ID"
// @Success 200 {object} response.Envelope{data=Author}
// @Failure 400 {object} response.Error
// @Failure 404 {object} response.Error
// @Failure 500 {object} response.Error
// @Router /authors/{id} [get]
func (h Handler) FindOne(c *gin.Context) {
	id, ok := authorID(c)
	if !ok {
		return
	}

	author, err := h.service.FindOne(c.Request.Context(), id)
	if err != nil {
		appmiddleware.AddError(c, err, "failed to get author")
		return
	}
	response.OK(c, http.StatusOK, author)
}

// Create godoc
// @Summary Create an author
// @Description Creates an author with a generated integer ID.
// @Tags authors
// @Accept json
// @Produce json
// @Param author body Request true "Author data"
// @Success 201 {object} response.Mutation
// @Failure 400 {object} response.Error
// @Failure 409 {object} response.Error
// @Failure 500 {object} response.Error
// @Router /authors [post]
func (h Handler) Create(c *gin.Context) {
	var request Request
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid author data")
		return
	}

	author, err := h.service.Create(c.Request.Context(), request.Name)
	if err != nil {
		appmiddleware.AddError(c, err, "failed to create author")
		return
	}
	response.Success(c, http.StatusCreated, "author created successfully", strconv.Itoa(author.ID))
}

// Update godoc
// @Summary Update an author
// @Tags authors
// @Accept json
// @Produce json
// @Param id path int true "Author ID"
// @Param author body Request true "Author data"
// @Success 200 {object} response.Mutation
// @Failure 400 {object} response.Error
// @Failure 404 {object} response.Error
// @Failure 409 {object} response.Error
// @Failure 500 {object} response.Error
// @Router /authors/{id} [put]
func (h Handler) Update(c *gin.Context) {
	id, ok := authorID(c)
	if !ok {
		return
	}

	var request Request
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid author data")
		return
	}

	if err := h.service.Update(c.Request.Context(), id, request.Name); err != nil {
		appmiddleware.AddError(c, err, "failed to update author")
		return
	}
	response.Success(c, http.StatusOK, "author updated successfully", strconv.Itoa(id))
}

// Delete godoc
// @Summary Delete an author
// @Tags authors
// @Produce json
// @Param id path int true "Author ID"
// @Success 200 {object} response.Mutation
// @Failure 400 {object} response.Error
// @Failure 404 {object} response.Error
// @Failure 500 {object} response.Error
// @Router /authors/{id} [delete]
func (h Handler) Delete(c *gin.Context) {
	id, ok := authorID(c)
	if !ok {
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		appmiddleware.AddError(c, err, "failed to delete author")
		return
	}
	response.Success(c, http.StatusOK, "author deleted successfully", strconv.Itoa(id))
}

func authorID(c *gin.Context) (int, bool) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		response.Fail(c, http.StatusBadRequest, "invalid author id")
		return 0, false
	}
	return id, true
}
