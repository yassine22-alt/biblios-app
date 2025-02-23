package json

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/yassine22-alt/biblios-app/api/internal/model"
)

type JsonBookStore struct {
	filename string
	mutex    sync.RWMutex
	lastID   int
	books    []model.Book
}

type BooksData struct {
	Books []model.Book `json:"books"`
}

func NewJsonBookStore() *JsonBookStore {

	store := &JsonBookStore{
		filename: "../data/books.json",
		books:    make([]model.Book, 0),
	}

	if err := store.loadFromFile(); err != nil {
		panic(err)
	}

	return store

}

func (s *JsonBookStore) loadFromFile() error {
	if err := os.MkdirAll("../data", 0755); err != nil {
		return err
	}

	if _, err := os.Stat(s.filename); os.IsNotExist(err) {
		initialData := BooksData{Books: []model.Book{}}
		data, _ := json.MarshalIndent(initialData, "", "  ")
		if err := os.WriteFile(s.filename, data, 0644); err != nil {
			return err
		}
	}

	data, err := os.ReadFile(s.filename)
	if err != nil {
		return err
	}

	var booksData BooksData
	if err := json.Unmarshal(data, &booksData); err != nil {
		return err
	}

	s.books = booksData.Books

	for _, book := range s.books {
		if book.ID > s.lastID {
			s.lastID = book.ID
		}
	}

	return nil
}

func (s *JsonBookStore) SaveToFile() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	data, err := json.MarshalIndent(BooksData{Books: s.books}, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.filename, data, 0644)
}

func (s *JsonBookStore) getNextID() int {
	s.lastID++
	return s.lastID
}

func (s *JsonBookStore) CreateBook(ctx context.Context, book model.Book) (model.Book, error) {
	select {
	case <-ctx.Done():
		return model.Book{}, ctx.Err()
	default:
		s.mutex.Lock()
		defer s.mutex.Unlock()

		book.ID = s.getNextID()
		s.books = append(s.books, book)

		return book, nil
	}
}

func (s *JsonBookStore) GetBook(ctx context.Context, id int) (model.Book, error) {
	select {
	case <-ctx.Done():
		return model.Book{}, ctx.Err()
	default:
		s.mutex.RLock()
		defer s.mutex.RUnlock()
		for _, book := range s.books {
			if book.ID == id {
				return book, nil
			}
		}
		return model.Book{}, fmt.Errorf("book with id %d not found", id)
	}
}

func (s *JsonBookStore) UpdateBook(ctx context.Context, id int, updatedBook model.Book) (model.Book, error) {
	select {
	case <-ctx.Done():
		return model.Book{}, ctx.Err()
	default:

		s.mutex.Lock()
		defer s.mutex.Unlock()

		for i, book := range s.books {
			if book.ID == id {
				s.books[i] = updatedBook
				return updatedBook, nil
			}
		}
		return model.Book{}, fmt.Errorf("book with id %d not found", id)
	}
}

func (s *JsonBookStore) DeleteBook(ctx context.Context, id int) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		s.mutex.Lock()
		defer s.mutex.Unlock()

		for i, book := range s.books {
			if book.ID == id {
				s.books = append(s.books[:i], s.books[i+1:]...)
				return nil
			}
		}
		return fmt.Errorf("book with id %d not found", id)
	}
}
func (s *JsonBookStore) SearchBooks(ctx context.Context, params map[string]string) ([]model.Book, error) {
	// Uncomment the line below if you wanna test context timeout
	// time.Sleep(6 * time.Second)
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		s.mutex.RLock()
		defer s.mutex.RUnlock()

		if params == nil {
			return s.books, nil
		}

		filteredBooks := []model.Book{}

		for _, book := range s.books {
			matches := true

			for key, value := range params {

				switch key {
				case "title":
					if !strings.EqualFold(book.Title, value) {
						matches = false
					}
				case "author":
					authorID, err := strconv.Atoi(value)
					if err != nil || book.AuthorID != authorID {
						matches = false
					}
				case "genre":
					value := strings.ToLower(value)
					genreMatches := false

					for _, genre := range book.Genres {
						if strings.ToLower(genre) == value {
							genreMatches = true
							break
						}
					}
					if !genreMatches {
						matches = false
					}
				case "year":
					if book.PublishedAt.Format("2006-01-02") != value {
						matches = false
					}
				default:
					matches = true
				}
			}

			if matches {
				filteredBooks = append(filteredBooks, book)
			}
		}

		return filteredBooks, nil
	}
}
