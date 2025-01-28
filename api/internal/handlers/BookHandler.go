package handlers

import (
	"bookstore/api/api/internal/errors"
	"bookstore/api/api/internal/model"
	"bookstore/api/api/internal/service"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type BookHandler struct {
	bookService *service.BookService
}

func NewBookHandler(bookService *service.BookService) *BookHandler {
	return &BookHandler{
		bookService: bookService,
	}
}

func (h *BookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		h.CreateBook(w, r)
	} else if r.Method == http.MethodGet {
		h.GetBooks(w, r)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors.Error{Message: "Request not allowed"})
	}
}

func (h *BookHandler) ServeHTTPById(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPut {
		h.UpdateBook(w, r)
	} else if r.Method == http.MethodDelete {
		h.DeleteBook(w, r)
	} else if r.Method == http.MethodGet {
		h.GetBook(w, r)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors.Error{Message: "Request not allowed"})
	}
}

func (h *BookHandler) GetBooks(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)

	defer cancel()

	params := r.URL.Query()

	searchParams := make(map[string]string)
	// here I only take nonempty params
	for key, value := range params {
		if len(value) > 0 && value[0] != "" {
			searchParams[key] = value[0]
		}
	}
	books, err := h.bookService.SearchBooks(ctx, searchParams)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(books)

}

func (h *BookHandler) GetBook(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	idString := r.PathValue("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	book, err1 := h.bookService.GetBook(ctx, int(id))

	if err1 != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errors.Error{Message: "book not found"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(book)
}

func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	idString := r.PathValue("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	err1 := h.bookService.DeleteBook(ctx, int(id))
	if err1 != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errors.Error{Message: "Book not found"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)

	defer cancel()

	idString := r.PathValue("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var bookInput model.BookInput
	err1 := decoder.Decode(&bookInput)
	if err1 != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors.Error{Message: "Book not updated"})
		return
	}

	book, err2 := h.bookService.UpdateBook(ctx, int(id), bookInput)
	if err2 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.Error{Message: err2.Error()})
		return
	}
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(book)
}

func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	decoder := json.NewDecoder(r.Body)
	var bookInput model.BookInput
	err := decoder.Decode(&bookInput)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors.Error{Message: "Invalid book payload"})
		return
	}

	book, err := h.bookService.CreateBook(ctx, bookInput)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors.Error{Message: err.Error()})
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}
