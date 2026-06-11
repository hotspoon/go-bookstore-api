package book

import (
	"errors"
	"net/http"
)

var (
	ErrBookNotFound      = errors.New("book not found")
	ErrBookAlreadyExists = errors.New("book already exists")
)

func HTTPErrorMapper(err error) (int, string, bool) {
	switch {
	case errors.Is(err, ErrBookNotFound):
		return http.StatusNotFound, "book not found", true
	case errors.Is(err, ErrBookAlreadyExists):
		return http.StatusConflict, "book already exists", true
	default:
		return 0, "", false
	}
}
