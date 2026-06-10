package book

import (
	"context"
	"errors"
	"strings"
	"testing"
)

type repositoryStub struct {
	books []Book
	err   error
}

func (r repositoryStub) FindAll(context.Context) ([]Book, error) {
	return r.books, r.err
}

func (r repositoryStub) FindOne(context.Context, string) (Book, error) {
	return Book{}, r.err
}

func (r repositoryStub) Create(context.Context, Book) error {
	return r.err
}

func (r repositoryStub) Update(context.Context, string, Book) error {
	return r.err
}

func (r repositoryStub) Delete(context.Context, string) error {
	return r.err
}

func TestServiceFindAll(t *testing.T) {
	expected := []Book{{ID: "book-1", Title: "First Book"}}
	service := NewService(repositoryStub{books: expected})

	books, err := service.FindAll(context.Background())
	if err != nil {
		t.Fatalf("FindAll() error = %v", err)
	}

	if len(books) != 1 || books[0].ID != expected[0].ID {
		t.Fatalf("FindAll() books = %#v, want %#v", books, expected)
	}
}

func TestServiceFindAllWrapsRepositoryError(t *testing.T) {
	repositoryErr := errors.New("database unavailable")
	service := NewService(repositoryStub{err: repositoryErr})

	_, err := service.FindAll(context.Background())
	if err == nil {
		t.Fatal("FindAll() error = nil, want an error")
	}

	if !errors.Is(err, repositoryErr) {
		t.Fatalf("FindAll() error = %v, want wrapped repository error", err)
	}

	if !strings.Contains(err.Error(), "get all books") {
		t.Fatalf("FindAll() error = %q, want service context", err)
	}
}

func TestServiceFindOnePreservesNotFoundError(t *testing.T) {
	service := NewService(repositoryStub{err: ErrBookNotFound})

	_, err := service.FindOne(context.Background(), "unknown")
	if !errors.Is(err, ErrBookNotFound) {
		t.Fatalf("FindOne() error = %v, want ErrBookNotFound", err)
	}
}

func TestServiceCreateGeneratesID(t *testing.T) {
	service := &service{
		repository: repositoryStub{},
		idGenerator: func() (string, error) {
			return "0123456789abcdefabcd", nil
		},
	}

	book, err := service.Create(context.Background(), Book{Title: "New Book"})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if book.ID != "0123456789abcdefabcd" {
		t.Fatalf("Create() ID = %q, want generated ID", book.ID)
	}
}

func TestGenerateBookIDFitsDatabaseColumn(t *testing.T) {
	id, err := generateBookID()
	if err != nil {
		t.Fatalf("generateBookID() error = %v", err)
	}

	if len(id) != 20 {
		t.Fatalf("generateBookID() length = %d, want 20", len(id))
	}
}
