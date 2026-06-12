package author

import (
	"context"
	"fmt"
	"strings"
)

type Service interface {
	FindAll(ctx context.Context) ([]Author, error)
	FindOne(ctx context.Context, id int) (Author, error)
	Create(ctx context.Context, name string) (Author, error)
	Update(ctx context.Context, id int, name string) error
	Delete(ctx context.Context, id int) error
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{repository: repository}
}

func (s *service) FindAll(ctx context.Context) ([]Author, error) {
	authors, err := s.repository.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all authors: %w", err)
	}
	return authors, nil
}

func (s *service) FindOne(ctx context.Context, id int) (Author, error) {
	author, err := s.repository.FindOne(ctx, id)
	if err != nil {
		return Author{}, fmt.Errorf("get author by id: %w", err)
	}
	return author, nil
}

func (s *service) Create(ctx context.Context, name string) (Author, error) {
	author, err := s.repository.Create(ctx, strings.TrimSpace(name))
	if err != nil {
		return Author{}, fmt.Errorf("create author: %w", err)
	}
	return author, nil
}

func (s *service) Update(ctx context.Context, id int, name string) error {
	if err := s.repository.Update(ctx, id, strings.TrimSpace(name)); err != nil {
		return fmt.Errorf("update author: %w", err)
	}
	return nil
}

func (s *service) Delete(ctx context.Context, id int) error {
	if err := s.repository.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete author: %w", err)
	}
	return nil
}
