package book

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

type serviceStub struct {
	books []Book
	book  Book
	err   error
}

func (s serviceStub) FindAll(context.Context) ([]Book, error) {
	return s.books, s.err
}

func (s serviceStub) FindOne(context.Context, string) (Book, error) {
	return s.book, s.err
}

func (s serviceStub) Create(context.Context, Book) (Book, error) {
	return s.book, s.err
}

func (s serviceStub) Update(context.Context, string, Book) error {
	return s.err
}

func (s serviceStub) Delete(context.Context, string) error {
	return s.err
}

func TestHandlerFindAll(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHandler(serviceStub{
		books: []Book{{ID: "book-1", Title: "First Book"}},
	})
	router := gin.New()
	router.GET("/books", handler.FindAll)

	request := httptest.NewRequest(http.MethodGet, "/books", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}

	if !strings.Contains(recorder.Body.String(), `"data":[{"id":"book-1"`) {
		t.Fatalf("body = %s, want book data envelope", recorder.Body.String())
	}
}

func TestHandlerFindAllReturnsInternalServerError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHandler(serviceStub{err: errors.New("service failed")})
	router := gin.New()
	router.GET("/books", handler.FindAll)

	request := httptest.NewRequest(http.MethodGet, "/books", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusInternalServerError)
	}

	if recorder.Body.String() != `{"message":"failed to get books"}` {
		t.Fatalf("body = %s, want generic error response", recorder.Body.String())
	}
}

func TestHandlerFindOneAcceptsStringID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHandler(serviceStub{
		book: Book{ID: "0-7475-3269-9", Title: "Harry Potter"},
	})
	router := gin.New()
	router.GET("/books/:id", handler.FindOne)

	request := httptest.NewRequest(http.MethodGet, "/books/0-7475-3269-9", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}

	if !strings.Contains(recorder.Body.String(), `"id":"0-7475-3269-9"`) {
		t.Fatalf("body = %s, want string book ID", recorder.Body.String())
	}
}

func TestHandlerFindOneReturnsNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHandler(serviceStub{err: ErrBookNotFound})
	router := gin.New()
	router.GET("/books/:id", handler.FindOne)

	request := httptest.NewRequest(http.MethodGet, "/books/unknown", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusNotFound)
	}

	if recorder.Body.String() != `{"message":"book not found"}` {
		t.Fatalf("body = %s, want not found response", recorder.Body.String())
	}
}

func TestHandlerFindOneReturnsInternalServerError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHandler(serviceStub{err: errors.New("database unavailable")})
	router := gin.New()
	router.GET("/books/:id", handler.FindOne)

	request := httptest.NewRequest(http.MethodGet, "/books/book-1", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusInternalServerError)
	}

	if recorder.Body.String() != `{"message":"failed to get book"}` {
		t.Fatalf("body = %s, want generic error response", recorder.Body.String())
	}
}

func TestHandlerCreateReturnsSuccessMessage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHandler(serviceStub{
		book: Book{ID: "0123456789abcdefabcd", PubID: 1, Title: "New Book"},
	})
	router := gin.New()
	router.POST("/books", handler.Create)

	body := `{"pub_id":1,"title":"New Book","price":10,"quantity":5,"prod_year":2026}`
	request := httptest.NewRequest(http.MethodPost, "/books", strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusCreated)
	}

	if recorder.Body.String() != `{"message":"book created successfully","id":"0123456789abcdefabcd"}` {
		t.Fatalf("body = %s, want success message", recorder.Body.String())
	}
}

func TestHandlerCreateRejectsIDOnlyInput(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHandler(serviceStub{})
	router := gin.New()
	router.POST("/books", handler.Create)

	request := httptest.NewRequest(http.MethodPost, "/books", strings.NewReader(`{"id":"client-id"}`))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusBadRequest)
	}
}

func TestHandlerUpdateReturnsSuccessMessage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHandler(serviceStub{})
	router := gin.New()
	router.PUT("/books/:id", handler.Update)

	body := `{"pub_id":1,"title":"Updated Book","price":20,"quantity":4,"prod_year":2026}`
	request := httptest.NewRequest(http.MethodPut, "/books/book-1", strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}

	if recorder.Body.String() != `{"message":"book updated successfully","id":"book-1"}` {
		t.Fatalf("body = %s, want success message", recorder.Body.String())
	}
}

func TestHandlerDeleteReturnsSuccessMessage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHandler(serviceStub{})
	router := gin.New()
	router.DELETE("/books/:id", handler.Delete)

	request := httptest.NewRequest(http.MethodDelete, "/books/book-1", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}

	if recorder.Body.String() != `{"message":"book deleted successfully","id":"book-1"}` {
		t.Fatalf("body = %s, want success message", recorder.Body.String())
	}
}

func TestHandlerCreateAndUpdateRejectDuplicateBook(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHandler(serviceStub{err: ErrBookAlreadyExists})
	router := gin.New()
	router.POST("/books", handler.Create)
	router.PUT("/books/:id", handler.Update)

	body := `{"pub_id":1,"title":"Duplicate","price":20,"quantity":4,"prod_year":2026}`
	for _, test := range []struct {
		method string
		path   string
	}{
		{method: http.MethodPost, path: "/books"},
		{method: http.MethodPut, path: "/books/book-1"},
	} {
		request := httptest.NewRequest(test.method, test.path, strings.NewReader(body))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, request)

		if recorder.Code != http.StatusConflict {
			t.Fatalf("%s status = %d, want %d", test.method, recorder.Code, http.StatusConflict)
		}

		if recorder.Body.String() != `{"message":"book already exists"}` {
			t.Fatalf("%s body = %s, want duplicate response", test.method, recorder.Body.String())
		}
	}
}

func TestHandlerUpdateAndDeleteReturnNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHandler(serviceStub{err: ErrBookNotFound})
	router := gin.New()
	router.PUT("/books/:id", handler.Update)
	router.DELETE("/books/:id", handler.Delete)

	body := `{"pub_id":1,"title":"Updated Book","price":20,"quantity":4,"prod_year":2026}`
	updateRequest := httptest.NewRequest(http.MethodPut, "/books/unknown", strings.NewReader(body))
	updateRequest.Header.Set("Content-Type", "application/json")
	updateRecorder := httptest.NewRecorder()
	router.ServeHTTP(updateRecorder, updateRequest)

	if updateRecorder.Code != http.StatusNotFound {
		t.Fatalf("update status = %d, want %d", updateRecorder.Code, http.StatusNotFound)
	}

	deleteRequest := httptest.NewRequest(http.MethodDelete, "/books/unknown", nil)
	deleteRecorder := httptest.NewRecorder()
	router.ServeHTTP(deleteRecorder, deleteRequest)

	if deleteRecorder.Code != http.StatusNotFound {
		t.Fatalf("delete status = %d, want %d", deleteRecorder.Code, http.StatusNotFound)
	}
}
