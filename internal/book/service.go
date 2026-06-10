package book

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// Service defines the business operations available for books.
type Service interface {
	FindAll(ctx context.Context) ([]Book, error)
	FindOne(ctx context.Context, id string) (Book, error)
	Create(ctx context.Context, book Book) (Book, error)
	Update(ctx context.Context, id string, book Book) error
	Delete(ctx context.Context, id string) error
}

type service struct {
	repository  Repository
	idGenerator func() (string, error)
}

func NewService(repository Repository) Service {
	return &service{
		repository:  repository,
		idGenerator: generateBookID,
	}
}

func (s *service) FindAll(ctx context.Context) ([]Book, error) {
	books, err := s.repository.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all books: %w", err)
	}

	return books, nil
}

func (s *service) FindOne(ctx context.Context, id string) (Book, error) {
	book, err := s.repository.FindOne(ctx, id)
	if err != nil {
		return Book{}, fmt.Errorf("get book by id: %w", err)
	}

	return book, nil
}

func (s *service) Create(ctx context.Context, book Book) (Book, error) {
	id, err := s.idGenerator()
	if err != nil {
		return Book{}, fmt.Errorf("generate book id: %w", err)
	}

	book.ID = id
	if err := s.repository.Create(ctx, book); err != nil {
		return Book{}, fmt.Errorf("create book: %w", err)
	}

	return book, nil
}

func (s *service) Update(ctx context.Context, id string, book Book) error {
	if err := s.repository.Update(ctx, id, book); err != nil {
		return fmt.Errorf("update book: %w", err)
	}

	return nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	if err := s.repository.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete book: %w", err)
	}

	return nil
}

func generateBookID() (string, error) {
	bytes := make([]byte, 10)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}
