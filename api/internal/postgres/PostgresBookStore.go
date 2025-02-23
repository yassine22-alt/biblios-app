package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/yassine22-alt/biblios-app/api/internal/model"
)

type PostgresBookStore struct {
	db *pgx.Conn
}

func NewPostgresBookStore(db *pgx.Conn) *PostgresBookStore {
	return &PostgresBookStore{db: db}
}

func (s *PostgresBookStore) CreateBook(ctx context.Context, book model.Book) (model.Book, error) {
	query := `
    INSERT INTO books (title, author_id, published_at, price, stock)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING id, title, author_id, published_at, price, stock`

	err := s.db.QueryRow(ctx, query,
		book.Title,
		book.AuthorID,
		book.PublishedAt,
		book.Price,
		book.Stock,
	).Scan(
		&book.ID,
		&book.Title,
		&book.AuthorID,
		&book.PublishedAt,
		&book.Price,
		&book.Stock,
	)

	if err != nil {
		return model.Book{}, fmt.Errorf("error creating book")
	}
	return book, nil
}

func (s *PostgresBookStore) GetBook(ctx context.Context, id int) (model.Book, error) {
	query := `
    SELECT id, title, author_id, published_at, price, stock
    FROM books 
    WHERE id = $1`

	var book model.Book
	err := s.db.QueryRow(ctx, query, id).Scan(
		&book.ID,
		&book.Title,
		&book.AuthorID,
		&book.PublishedAt,
		&book.Price,
		&book.Stock,
	)

	if err != nil {
		return model.Book{}, fmt.Errorf("error getting book: %v", err)
	}

	return book, nil
}

func (s *PostgresBookStore) UpdateBook(ctx context.Context, id int, book model.Book) (model.Book, error) {
	query := `
    UPDATE books 
    SET title = $1, author_id = $2, published_at = $3, price = $4, stock = $5
    WHERE id = $6
    RETURNING id, title, author_id, published_at, price, stock`

	err := s.db.QueryRow(ctx, query,
		book.Title,
		book.AuthorID,
		book.PublishedAt,
		book.Price,
		book.Stock,
		id,
	).Scan(
		&book.ID,
		&book.Title,
		&book.AuthorID,
		&book.PublishedAt,
		&book.Price,
		&book.Stock,
	)

	if err != nil {
		return model.Book{}, fmt.Errorf("error updating book: %v", err)
	}

	return book, nil
}

func (s *PostgresBookStore) DeleteBook(ctx context.Context, id int) error {
	query := `DELETE FROM books WHERE id = $1`
	_, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting book: %v", err)
	}
	return nil
}

func (s *PostgresBookStore) SearchBooks(ctx context.Context, params map[string]string) ([]model.Book, error) {
	query := `
    SELECT id, title, author_id, published_at, price, stock 
    FROM books`

	var conditions []string
	var args []interface{}
	argPosition := 1

	if title, ok := params["title"]; ok {
		conditions = append(conditions, fmt.Sprintf("title ILIKE $%d", argPosition))
		args = append(args, "%"+title+"%")
		argPosition++
	}

	if authorID, ok := params["author_id"]; ok {
		conditions = append(conditions, fmt.Sprintf("author_id = $%d", argPosition))
		args = append(args, authorID)
		argPosition++
	}

	if len(conditions) > 0 {
		query += " WHERE " + conditions[0]
		for i := 1; i < len(conditions); i++ {
			query += " AND " + conditions[i]
		}
	}

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error searching books: %v", err)
	}
	defer rows.Close()

	var books []model.Book
	for rows.Next() {
		var book model.Book
		err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.AuthorID,
			&book.PublishedAt,
			&book.Price,
			&book.Stock,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning book row: %v", err)
		}
		books = append(books, book)
	}

	return books, nil
}
