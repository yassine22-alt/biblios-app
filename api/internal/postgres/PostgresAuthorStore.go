package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/yassine22-alt/biblios-app/api/internal/model"
)

type PostgresAuthorStore struct {
	db *pgx.Conn
}

func NewPostgresAuthorStore(db *pgx.Conn) *PostgresAuthorStore {
	return &PostgresAuthorStore{db: db}
}

func (s *PostgresAuthorStore) CreateAuthor(ctx context.Context, author model.Author) (model.Author, error) {
	query := `
    INSERT INTO authors (first_name, last_name, bio)
    VALUES ($1, $2, $3)
    RETURNING id, first_name, last_name, bio`

	err := s.db.QueryRow(ctx, query,
		author.FirstName,
		author.LastName,
		author.Bio,
	).Scan(&author.ID, &author.FirstName, &author.LastName, &author.Bio)

	if err != nil {
		return model.Author{}, fmt.Errorf("error creating author: %v", err)
	}

	return author, nil
}

func (s *PostgresAuthorStore) GetAuthor(ctx context.Context, id int) (model.Author, error) {
	query := `
        SELECT id, first_name, last_name, bio 
        FROM authors 
        WHERE id = $1`

	var author model.Author
	err := s.db.QueryRow(ctx, query, id).Scan(
		&author.ID,
		&author.FirstName,
		&author.LastName,
		&author.Bio,
	)

	return author, err
}

func (s *PostgresAuthorStore) UpdateAuthor(ctx context.Context, id int, author model.Author) (model.Author, error) {
	query := `
	UPDATE authors
	SET first_name = $1, last_name = $2, bio = $3
	WHERE id = $4
	RETURNING id`

	err := s.db.QueryRow(ctx, query,
		author.FirstName,
		author.LastName,
		author.Bio,
		id,
	).Scan(&author.ID)
	return author, err
}

func (s *PostgresAuthorStore) DeleteAuthor(ctx context.Context, id int) error {
	query := `DELETE FROM authors WHERE id = $1`
	_, err := s.db.Exec(ctx, query, id)
	return err
}

func (s *PostgresAuthorStore) SearchAuthors(ctx context.Context, params map[string]string) ([]model.Author, error) {
	query := `SELECT id, first_name, last_name, bio FROM authors`

	// Add WHERE clauses based on params
	args := []interface{}{}
	if name, ok := params["firstName"]; ok {
		query += ` WHERE first_name ILIKE $1`
		args = append(args, "%"+name+"%")
	}

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var authors []model.Author
	for rows.Next() {
		var author model.Author
		err := rows.Scan(
			&author.ID,
			&author.FirstName,
			&author.LastName,
			&author.Bio,
		)
		if err != nil {
			return nil, err
		}
		authors = append(authors, author)
	}
	return authors, nil
}
