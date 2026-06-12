package author

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type Repository interface {
	FindAll(ctx context.Context) ([]Author, error)
	FindOne(ctx context.Context, id int) (Author, error)
	Create(ctx context.Context, name string) (Author, error)
	Update(ctx context.Context, id int, name string) error
	Delete(ctx context.Context, id int) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAll(ctx context.Context) ([]Author, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, s_name FROM Source ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("find all authors: %w", err)
	}
	defer rows.Close()

	authors := make([]Author, 0)
	for rows.Next() {
		var author Author
		if err := rows.Scan(&author.ID, &author.Name); err != nil {
			return nil, fmt.Errorf("scan author: %w", err)
		}
		authors = append(authors, author)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate authors: %w", err)
	}

	return authors, nil
}

func (r *repository) FindOne(ctx context.Context, id int) (Author, error) {
	var author Author
	err := r.db.QueryRowContext(ctx, `
		SELECT id, s_name
		FROM Source
		WHERE id = ?
	`, id).Scan(&author.ID, &author.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return Author{}, fmt.Errorf("find author %d: %w", id, ErrAuthorNotFound)
		}
		return Author{}, fmt.Errorf("find author: %w", err)
	}

	return author, nil
}

func (r *repository) Create(ctx context.Context, name string) (Author, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return Author{}, fmt.Errorf("begin create author: %w", err)
	}
	defer tx.Rollback()

	var id int
	if err := tx.QueryRowContext(ctx, `SELECT COALESCE(MAX(id), 0) + 1 FROM Source`).Scan(&id); err != nil {
		return Author{}, fmt.Errorf("generate author id: %w", err)
	}

	if _, err := tx.ExecContext(ctx, `INSERT INTO Source (id, s_name) VALUES (?, ?)`, id, name); err != nil {
		if isDuplicateAuthorError(err) {
			return Author{}, fmt.Errorf("insert author: %w", ErrAuthorAlreadyExists)
		}
		return Author{}, fmt.Errorf("insert author: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return Author{}, fmt.Errorf("commit create author: %w", err)
	}

	return Author{ID: id, Name: name}, nil
}

func (r *repository) Update(ctx context.Context, id int, name string) error {
	result, err := r.db.ExecContext(ctx, `UPDATE Source SET s_name = ? WHERE id = ?`, name, id)
	if err != nil {
		if isDuplicateAuthorError(err) {
			return fmt.Errorf("update author: %w", ErrAuthorAlreadyExists)
		}
		return fmt.Errorf("update author: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("update author %d: %w", id, ErrAuthorNotFound)
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, id int) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM Source WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete author: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("delete author %d: %w", id, ErrAuthorNotFound)
	}
	return nil
}

func isDuplicateAuthorError(err error) bool {
	message := err.Error()
	return strings.Contains(message, "idx_source_unique_name") ||
		strings.Contains(message, "UNIQUE constraint failed: index")
}
