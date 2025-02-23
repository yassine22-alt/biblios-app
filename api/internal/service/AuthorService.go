package service

import (
	"context"
	"errors"

	"github.com/yassine22-alt/biblios-app/api/internal/model"
	"github.com/yassine22-alt/biblios-app/api/internal/repository"
)

type AuthorService struct {
	repo repository.AuthorStore
}

func NewAuthorService(repo repository.AuthorStore) *AuthorService {
	return &AuthorService{
		repo: repo,
	}
}

func (s *AuthorService) CreateAuthor(ctx context.Context, authorInput model.AuthorInput) (model.Author, error) {
	if err := ctx.Err(); err != nil {
		return model.Author{}, err
	}
	author := model.Author{
		FirstName: authorInput.FirstName,
		LastName:  authorInput.LastName,
		Bio:       authorInput.Bio,
	}

	if author.FirstName == "" || author.LastName == "" {
		return model.Author{}, errors.New("author name is mandatory")
	}

	return s.repo.CreateAuthor(ctx, author)
}

func (s *AuthorService) GetAuthor(ctx context.Context, id int) (model.Author, error) {
	if err := ctx.Err(); err != nil {
		return model.Author{}, err
	}
	return s.repo.GetAuthor(ctx, id)
}

func (s *AuthorService) UpdateAuthor(ctx context.Context, id int, authorInput model.AuthorInput) (model.Author, error) {
	if err := ctx.Err(); err != nil {
		return model.Author{}, err
	}

	existingAuthor, err := s.repo.GetAuthor(ctx, id)
	if err != nil {
		return model.Author{}, err
	}

	updatedAuthor := model.Author{
		ID:        existingAuthor.ID,
		FirstName: authorInput.FirstName,
		LastName:  authorInput.LastName,
		Bio:       authorInput.Bio,
	}

	if updatedAuthor.FirstName == "" || updatedAuthor.LastName == "" {
		return model.Author{}, errors.New("author name is mandatory")
	}

	return s.repo.UpdateAuthor(ctx, id, updatedAuthor)
}

func (s *AuthorService) DeleteAuthor(ctx context.Context, id int) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return s.repo.DeleteAuthor(ctx, id)
}

func (s *AuthorService) SearchAuthors(ctx context.Context, params map[string]string) ([]model.Author, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return s.repo.SearchAuthors(ctx, params)
}
