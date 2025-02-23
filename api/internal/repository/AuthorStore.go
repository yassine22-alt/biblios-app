package repository

import (
	"context"

	"github.com/yassine22-alt/biblios-app/api/internal/model"
)

type AuthorStore interface {
	CreateAuthor(ctx context.Context, author model.Author) (model.Author, error)
	GetAuthor(ctx context.Context, id int) (model.Author, error)
	UpdateAuthor(ctx context.Context, id int, author model.Author) (model.Author, error)
	DeleteAuthor(ctx context.Context, id int) error
	SearchAuthors(ctx context.Context, params map[string]string) ([]model.Author, error)
}
