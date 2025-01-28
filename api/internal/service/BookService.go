package service

import (
	"bookstore/api/api/internal/model"
	"bookstore/api/api/internal/repository"
	"context"
	"errors"
	"time"
)

type BookService struct {
	repo       repository.BookStore
	repoAuthor repository.AuthorStore
	currentID  int
}

func NewBookService(repo repository.BookStore, repoAuthor repository.AuthorStore) *BookService {
	return &BookService{
		repo:       repo,
		repoAuthor: repoAuthor,
		currentID:  1,
	}
}

func (s *BookService) CreateBook(ctx context.Context, bookInput model.BookInput) (model.Book, error) {

	if err := ctx.Err(); err != nil {
		return model.Book{}, err
	}

	book := model.Book{
		ID:          s.currentID,
		Title:       bookInput.Title,
		AuthorID:    bookInput.AuthorID,
		Genres:      bookInput.Genres,
		PublishedAt: time.Now(),
		Price:       bookInput.Price,
		Stock:       bookInput.Stock,
	}
	s.currentID++

	if book.Stock < 0 || book.Price < 0 {
		return model.Book{}, errors.New("book details are invalid")
	}
	if book.Title == "" {
		return model.Book{}, errors.New("book title is mandatory")
	}

	_, err := s.repoAuthor.GetAuthor(ctx, bookInput.AuthorID)
	if err != nil {
		return model.Book{}, errors.New("author not found")
	}

	return s.repo.CreateBook(ctx, book)
}

func (s *BookService) GetBook(ctx context.Context, id int) (model.Book, error) {
	if err := ctx.Err(); err != nil {
		return model.Book{}, err
	}
	return s.repo.GetBook(ctx, id)
}

func (s *BookService) DeleteBook(ctx context.Context, id int) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return s.repo.DeleteBook(ctx, id)
}

func (s *BookService) UpdateBook(ctx context.Context, id int, bookInput model.BookInput) (model.Book, error) {

	if err := ctx.Err(); err != nil {
		return model.Book{}, err
	}

	existingBook, err := s.repo.GetBook(ctx, id)

	updatedBook := model.Book{
		ID:          existingBook.ID,
		Title:       bookInput.Title,
		AuthorID:    bookInput.AuthorID,
		Genres:      bookInput.Genres,
		PublishedAt: existingBook.PublishedAt,
		Price:       bookInput.Price,
		Stock:       bookInput.Stock,
	}

	if err != nil {
		return model.Book{}, err
	}
	if updatedBook.Stock < 0 || bookInput.Price < 0 {
		return model.Book{}, errors.New("book details are invalid")
	}
	if updatedBook.Title == "" {
		return model.Book{}, errors.New("book title is mandatory")
	}

	_, err1 := s.repoAuthor.GetAuthor(ctx, bookInput.AuthorID)
	if err1 != nil {
		return model.Book{}, errors.New("author not found")
	}

	if err := ctx.Err(); err != nil {
		return model.Book{}, err
	}

	return s.repo.UpdateBook(ctx, id, updatedBook)
}

func (s *BookService) SearchBooks(ctx context.Context, params map[string]string) ([]model.Book, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return s.repo.SearchBooks(ctx, params)
}
