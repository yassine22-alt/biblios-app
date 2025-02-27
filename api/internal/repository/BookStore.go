package repository

import (
	"bookstore/api/api/internal/model"
	"context"
)

type BookStore interface {
	CreateBook(ctx context.Context, book model.Book) (model.Book, error)
	GetBook(ctx context.Context, id int) (model.Book, error)
	UpdateBook(ctx context.Context, id int, book model.Book) (model.Book, error)
	DeleteBook(ctx context.Context, id int) error
	SearchBooks(ctx context.Context, params map[string]string) ([]model.Book, error)
}
