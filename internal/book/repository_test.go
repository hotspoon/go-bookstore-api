package book

import (
	"context"
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func TestRepositoryFindAll(t *testing.T) {
	db := openTestDB(t)

	_, err := db.Exec(`
		INSERT INTO Book (
			id, pub_id, title, price, category, quantity, b_format, prod_year, filesize
		) VALUES
			('book-2', 2, 'Second Book', 20.5, NULL, 4, NULL, 2025, NULL),
			('book-1', 1, 'First Book', 10.0, 'Fiction', 7, 'Paperback', 2024, 512)
	`)
	if err != nil {
		t.Fatalf("seed books: %v", err)
	}

	books, err := NewRepository(db).FindAll(context.Background())
	if err != nil {
		t.Fatalf("FindAll() error = %v", err)
	}

	if len(books) != 2 {
		t.Fatalf("FindAll() returned %d books, want 2", len(books))
	}

	if books[0].ID != "book-1" || books[1].ID != "book-2" {
		t.Fatalf("FindAll() IDs = [%q, %q], want sorted IDs", books[0].ID, books[1].ID)
	}

	if books[0].Category == nil || *books[0].Category != "Fiction" {
		t.Fatalf("first book category = %v, want Fiction", books[0].Category)
	}

	if books[1].Category != nil || books[1].Format != nil || books[1].FileSize != nil {
		t.Fatalf("nullable fields on second book were not scanned as nil")
	}
}

func TestRepositoryFindAllReturnsEmptySlice(t *testing.T) {
	db := openTestDB(t)

	books, err := NewRepository(db).FindAll(context.Background())
	if err != nil {
		t.Fatalf("FindAll() error = %v", err)
	}

	if books == nil {
		t.Fatal("FindAll() returned nil, want an empty slice")
	}

	if len(books) != 0 {
		t.Fatalf("FindAll() returned %d books, want 0", len(books))
	}
}

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()

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
		)
	`)
	if err != nil {
		t.Fatalf("create Book table: %v", err)
	}

	return db
}
