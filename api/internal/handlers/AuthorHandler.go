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

type AuthorHandler struct {
	authorService *service.AuthorService
}

func NewAuthorHandler(authorService *service.AuthorService) *AuthorHandler {
	return &AuthorHandler{
		authorService: authorService,
	}
}

func (h *AuthorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		h.CreateAuthor(w, r)
	} else if r.Method == http.MethodGet {
		h.GetAuthors(w, r)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors.Error{Message: "Request not allowed"})
	}
}

func (h *AuthorHandler) ServeHTTPById(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.GetAuthor(w, r)
	} else if r.Method == http.MethodPut {
		h.UpdateAuthor(w, r)
	} else if r.Method == http.MethodDelete {
		h.DeleteAuthor(w, r)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors.Error{Message: "Request not allowed"})
	}
}

func (h *AuthorHandler) CreateAuthor(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	decoder := json.NewDecoder(r.Body)
	var authorInput model.AuthorInput
	err := decoder.Decode(&authorInput)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors.Error{Message: "Invalid author payload"})
		return
	}

	author, err := h.authorService.CreateAuthor(ctx, authorInput)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors.Error{Message: err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(author)
}

func (h *AuthorHandler) GetAuthor(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	author, err := h.authorService.GetAuthor(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errors.Error{Message: "Author not found"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(author)
}

func (h *AuthorHandler) GetAuthors(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	params := r.URL.Query()
	searchParams := make(map[string]string)
	for key, value := range params {
		if len(value) > 0 && value[0] != "" {
			searchParams[key] = value[0]
		}
	}

	authors, err := h.authorService.SearchAuthors(ctx, searchParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(authors)
}

func (h *AuthorHandler) UpdateAuthor(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var authorInput model.AuthorInput
	err = decoder.Decode(&authorInput)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors.Error{Message: "Invalid author payload"})
		return
	}

	author, err := h.authorService.UpdateAuthor(ctx, id, authorInput)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors.Error{Message: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(author)
}

func (h *AuthorHandler) DeleteAuthor(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	err = h.authorService.DeleteAuthor(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errors.Error{Message: "Author not found"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
