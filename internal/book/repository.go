package book

import (
	"context"
	"database/sql"
	"fmt"
)

// Repository defines the data access operations available for books.
type Repository interface {
	FindAll(ctx context.Context) ([]Book, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAll(ctx context.Context) ([]Book, error) {
	const query = `
		SELECT id, pub_id, title, price, category, quantity, b_format, prod_year, filesize
		FROM Book
		ORDER BY id
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("find all books: %w", err)
	}
	defer rows.Close()

	books := make([]Book, 0)
	for rows.Next() {
		var book Book
		if err := rows.Scan(
			&book.ID,
			&book.PubID,
			&book.Title,
			&book.Price,
			&book.Category,
			&book.Quantity,
			&book.Format,
			&book.ProdYear,
			&book.FileSize,
		); err != nil {
			return nil, fmt.Errorf("scan book: %w", err)
		}

		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate books: %w", err)
	}

	return books, nil
}
