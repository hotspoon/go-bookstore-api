package book

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

// Repository defines the data access operations available for books.
type Repository interface {
	FindAll(ctx context.Context) ([]Book, error)
	FindOne(ctx context.Context, id string) (Book, error)
	Create(ctx context.Context, book Book) error
	Update(ctx context.Context, id string, book Book) error
	Delete(ctx context.Context, id string) error
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

func (r *repository) FindOne(ctx context.Context, id string) (Book, error) {
	const query = `
		SELECT id, pub_id, title, price, category, quantity, b_format, prod_year, filesize
		FROM Book
		WHERE id = ?
	`

	var book Book
	if err := r.db.QueryRowContext(ctx, query, id).Scan(
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
		if err == sql.ErrNoRows {
			return Book{}, fmt.Errorf("find book %q: %w", id, ErrBookNotFound)
		}

		return Book{}, fmt.Errorf("find book: %w", err)
	}

	return book, nil
}

func (r *repository) Create(ctx context.Context, book Book) error {
	result, err := r.db.ExecContext(ctx, `
		INSERT INTO Book (
			id, pub_id, title, price, category, quantity, b_format, prod_year, filesize
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		book.ID,
		book.PubID,
		book.Title,
		book.Price,
		book.Category,
		book.Quantity,
		book.Format,
		book.ProdYear,
		book.FileSize,
	)
	if err != nil {
		if isDuplicateBookError(err) {
			return fmt.Errorf("insert book: %w", ErrBookAlreadyExists)
		}

		return fmt.Errorf("insert book: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("insert book: no rows affected")
	}

	return nil
}

func (r *repository) Update(ctx context.Context, id string, book Book) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE Book
		SET pub_id = ?, title = ?, price = ?, category = ?, quantity = ?, b_format = ?, prod_year = ?, filesize = ?
		WHERE id = ?
	`,
		book.PubID,
		book.Title,
		book.Price,
		book.Category,
		book.Quantity,
		book.Format,
		book.ProdYear,
		book.FileSize,
		id,
	)
	if err != nil {
		if isDuplicateBookError(err) {
			return fmt.Errorf("update book: %w", ErrBookAlreadyExists)
		}

		return fmt.Errorf("update book: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("update book %q: %w", id, ErrBookNotFound)
	}

	return nil
}

func isDuplicateBookError(err error) bool {
	message := err.Error()
	return strings.Contains(message, "idx_book_unique_identity") ||
		strings.Contains(message, "UNIQUE constraint failed: index")
}

func (r *repository) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `
		DELETE FROM Book
		WHERE id = ?
	`, id)
	if err != nil {
		return fmt.Errorf("delete book: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("delete book %q: %w", id, ErrBookNotFound)
	}

	return nil
}
