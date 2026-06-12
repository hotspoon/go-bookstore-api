package author

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	_ "modernc.org/sqlite"
)

func TestRepositoryCRUD(t *testing.T) {
	db := openTestDB(t)
	repository := NewRepository(db)
	ctx := context.Background()

	created, err := repository.Create(ctx, "New Author")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if created.ID != 1 || created.Name != "New Author" {
		t.Fatalf("Create() = %#v", created)
	}

	found, err := repository.FindOne(ctx, created.ID)
	if err != nil {
		t.Fatalf("FindOne() error = %v", err)
	}
	if found != created {
		t.Fatalf("FindOne() = %#v, want %#v", found, created)
	}

	if err := repository.Update(ctx, created.ID, "Updated Author"); err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	authors, err := repository.FindAll(ctx)
	if err != nil {
		t.Fatalf("FindAll() error = %v", err)
	}
	if len(authors) != 1 || authors[0].Name != "Updated Author" {
		t.Fatalf("FindAll() = %#v", authors)
	}

	if err := repository.Delete(ctx, created.ID); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	if _, err := repository.FindOne(ctx, created.ID); !errors.Is(err, ErrAuthorNotFound) {
		t.Fatalf("FindOne() error = %v, want ErrAuthorNotFound", err)
	}
}

func TestRepositoryRejectsDuplicateName(t *testing.T) {
	db := openTestDB(t)
	repository := NewRepository(db)
	ctx := context.Background()

	if _, err := repository.Create(ctx, "Same Author"); err != nil {
		t.Fatalf("first Create() error = %v", err)
	}
	if _, err := repository.Create(ctx, " same author "); !errors.Is(err, ErrAuthorAlreadyExists) {
		t.Fatalf("duplicate Create() error = %v, want ErrAuthorAlreadyExists", err)
	}
}

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { db.Close() })

	_, err = db.Exec(`
		CREATE TABLE Source (
			id INTEGER NOT NULL PRIMARY KEY,
			s_name TEXT NOT NULL
		);
		CREATE UNIQUE INDEX idx_source_unique_name
		ON Source (LOWER(TRIM(s_name)));
	`)
	if err != nil {
		t.Fatalf("prepare database: %v", err)
	}

	return db
}
