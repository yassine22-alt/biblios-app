package json

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/yassine22-alt/biblios-app/api/internal/model"
)

type JsonAuthorStore struct {
	filename string
	mutex    sync.RWMutex
	lastID   int
	authors  []model.Author
}

type AuthorsData struct {
	Authors []model.Author `json:"authors"`
}

func NewJsonAuthorStore() *JsonAuthorStore {
	store := &JsonAuthorStore{
		filename: "../data/authors.json",
		authors:  make([]model.Author, 0),
	}

	if err := store.loadFromFile(); err != nil {
		panic(err)
	}

	return store
}

func (s *JsonAuthorStore) getNextID() int {
	s.lastID++
	return s.lastID
}

func (s *JsonAuthorStore) loadFromFile() error {
	if err := os.MkdirAll("../data", 0755); err != nil {
		return err
	}

	if _, err := os.Stat(s.filename); os.IsNotExist(err) {
		initialData := AuthorsData{Authors: []model.Author{}}
		data, _ := json.MarshalIndent(initialData, "", "  ")
		if err := os.WriteFile(s.filename, data, 0644); err != nil {
			return err
		}
	}

	data, err := os.ReadFile(s.filename)
	if err != nil {
		return err
	}

	var authorsData AuthorsData
	if err := json.Unmarshal(data, &authorsData); err != nil {
		return err
	}

	s.authors = authorsData.Authors

	for _, author := range s.authors {
		if author.ID > s.lastID {
			s.lastID = author.ID
		}
	}

	return nil
}

func (s *JsonAuthorStore) SaveToFile() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	data, err := json.MarshalIndent(AuthorsData{Authors: s.authors}, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.filename, data, 0644)
}

func (s *JsonAuthorStore) CreateAuthor(ctx context.Context, author model.Author) (model.Author, error) {
	select {
	case <-ctx.Done():
		return model.Author{}, ctx.Err()
	default:
		s.mutex.Lock()
		defer s.mutex.Unlock()

		author.ID = s.getNextID()
		s.authors = append(s.authors, author)
		return author, nil
	}
}

func (s *JsonAuthorStore) GetAuthor(ctx context.Context, id int) (model.Author, error) {
	select {
	case <-ctx.Done():
		return model.Author{}, ctx.Err()
	default:
		s.mutex.RLock()
		defer s.mutex.RUnlock()

		for _, author := range s.authors {
			if author.ID == id {
				return author, nil
			}
		}
		return model.Author{}, fmt.Errorf("author with id %d not found", id)
	}
}

func (s *JsonAuthorStore) UpdateAuthor(ctx context.Context, id int, updatedAuthor model.Author) (model.Author, error) {
	select {
	case <-ctx.Done():
		return model.Author{}, ctx.Err()
	default:
		s.mutex.Lock()
		defer s.mutex.Unlock()

		for i, author := range s.authors {
			if author.ID == id {
				s.authors[i] = updatedAuthor
				return updatedAuthor, nil
			}
		}
		return model.Author{}, fmt.Errorf("author with id %d not found", id)
	}
}

func (s *JsonAuthorStore) DeleteAuthor(ctx context.Context, id int) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		s.mutex.Lock()
		defer s.mutex.Unlock()

		for i, author := range s.authors {
			if author.ID == id {
				s.authors = append(s.authors[:i], s.authors[i+1:]...)
				return nil
			}
		}
		return fmt.Errorf("author with id %d not found", id)
	}
}

func (s *JsonAuthorStore) SearchAuthors(ctx context.Context, params map[string]string) ([]model.Author, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		s.mutex.Lock()
		defer s.mutex.Unlock()

		if params == nil {
			return s.authors, nil
		}

		result := []model.Author{}
		for _, author := range s.authors {
			matches := true
			for key, value := range params {
				switch key {
				case "firstName":
					if !strings.EqualFold(author.FirstName, value) {
						matches = false
					}
				case "lastName":
					if !strings.EqualFold(author.LastName, value) {
						matches = false
					}
				}
			}
			if matches {
				result = append(result, author)
			}
		}
		return result, nil
	}
}
