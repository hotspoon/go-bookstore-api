package author

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	appmiddleware "bookstore-api/internal/shared/middleware"

	"github.com/gin-gonic/gin"
)

type serviceStub struct {
	authors []Author
	author  Author
	err     error
}

func (s serviceStub) FindAll(context.Context) ([]Author, error) {
	return s.authors, s.err
}

func (s serviceStub) FindOne(context.Context, int) (Author, error) {
	return s.author, s.err
}

func (s serviceStub) Create(context.Context, string) (Author, error) {
	return s.author, s.err
}

func (s serviceStub) Update(context.Context, int, string) error {
	return s.err
}

func (s serviceStub) Delete(context.Context, int) error {
	return s.err
}

func TestHandlerFindAll(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewHandler(serviceStub{authors: []Author{{ID: 1, Name: "Ian H. Witten"}}})
	router := newAuthorRouter()
	router.GET("/authors", handler.FindAll)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/authors", nil))

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if !strings.Contains(recorder.Body.String(), `"name":"Ian H. Witten"`) {
		t.Fatalf("body = %s, want author", recorder.Body.String())
	}
}

func TestHandlerCreate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewHandler(serviceStub{author: Author{ID: 1001, Name: "New Author"}})
	router := newAuthorRouter()
	router.POST("/authors", handler.Create)

	request := httptest.NewRequest(http.MethodPost, "/authors", strings.NewReader(`{"name":"New Author"}`))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusCreated)
	}
	if recorder.Body.String() != `{"message":"author created successfully","id":"1001"}` {
		t.Fatalf("body = %s, want success response", recorder.Body.String())
	}
}

func TestHandlerRejectsInvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewHandler(serviceStub{})
	router := newAuthorRouter()
	router.GET("/authors/:id", handler.FindOne)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/authors/abc", nil))

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusBadRequest)
	}
}

func newAuthorRouter() *gin.Engine {
	router := gin.New()
	router.Use(appmiddleware.ErrorHandler(HTTPErrorMapper))
	return router
}
