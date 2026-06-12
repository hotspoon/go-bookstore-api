package author

import (
	"errors"
	"net/http"
)

var (
	ErrAuthorNotFound      = errors.New("author not found")
	ErrAuthorAlreadyExists = errors.New("author already exists")
)

func HTTPErrorMapper(err error) (int, string, bool) {
	switch {
	case errors.Is(err, ErrAuthorNotFound):
		return http.StatusNotFound, "author not found", true
	case errors.Is(err, ErrAuthorAlreadyExists):
		return http.StatusConflict, "author already exists", true
	default:
		return 0, "", false
	}
}
