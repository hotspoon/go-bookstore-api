package main

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"bookstore-api/internal/config"

	_ "modernc.org/sqlite"
)

func TestGetAllBooksRoute(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	t.Cleanup(func() {
		db.Close()
	})

	_, err = db.Exec(`
		CREATE TABLE Book (
			id TEXT PRIMARY KEY,
			pub_id INTEGER NOT NULL,
			title TEXT NOT NULL,
			price REAL NOT NULL,
			category TEXT,
			quantity INTEGER NOT NULL,
			b_format TEXT,
			prod_year INTEGER NOT NULL,
			filesize INTEGER
		);

		INSERT INTO Book (
			id, pub_id, title, price, category, quantity, b_format, prod_year, filesize
		) VALUES (
			'book-1', 1, 'First Book', 10, 'Fiction', 5, 'Paperback', 2024, 512
		);
	`)
	if err != nil {
		t.Fatalf("prepare database: %v", err)
	}

	router := setupRouter(config.Config{
		APIVersion:     "v1",
		AllowedOrigins: []string{"*"},
	}, db)

	request := httptest.NewRequest(http.MethodGet, "/api/v1/books", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}

	if !strings.Contains(recorder.Body.String(), `"title":"First Book"`) {
		t.Fatalf("body = %s, want seeded book", recorder.Body.String())
	}
}
