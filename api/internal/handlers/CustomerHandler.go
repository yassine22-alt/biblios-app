package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/yassine22-alt/biblios-app/api/internal/errors"
	"github.com/yassine22-alt/biblios-app/api/internal/model"
	"github.com/yassine22-alt/biblios-app/api/internal/service"
)

type CustomerHandler struct {
	customerService *service.CustomerService
}

func NewCustomerHandler(customerService *service.CustomerService) *CustomerHandler {
	return &CustomerHandler{
		customerService: customerService,
	}
}

func (h *CustomerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		h.CreateCustomer(w, r)
	} else if r.Method == http.MethodGet {
		h.GetCustomers(w, r)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors.Error{Message: "Request not allowed"})
	}
}

func (h *CustomerHandler) ServeHTTPById(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPut {
		h.UpdateCustomer(w, r)
	} else if r.Method == http.MethodDelete {
		h.DeleteCustomer(w, r)
	} else if r.Method == http.MethodGet {
		h.GetCustomer(w, r)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors.Error{Message: "Request not allowed"})
	}
}

func (h *CustomerHandler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	decoder := json.NewDecoder(r.Body)
	var customerInput model.CustomerInput
	err := decoder.Decode(&customerInput)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors.Error{Message: "Invalid customer payload"})
		return
	}

	customer, err := h.customerService.CreateCustomer(ctx, customerInput)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors.Error{Message: err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(customer)
}

func (h *CustomerHandler) GetCustomer(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	customer, err := h.customerService.GetCustomer(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errors.Error{Message: "Customer not found"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(customer)
}

func (h *CustomerHandler) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var customerInput model.CustomerInput
	err = decoder.Decode(&customerInput)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors.Error{Message: "Invalid customer payload"})
		return
	}

	customer, err := h.customerService.UpdateCustomer(ctx, id, customerInput)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors.Error{Message: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(customer)
}

func (h *CustomerHandler) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	err = h.customerService.DeleteCustomer(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errors.Error{Message: "Customer not found"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *CustomerHandler) GetCustomers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	params := r.URL.Query()
	searchParams := make(map[string]string)
	for key, value := range params {
		if len(value) > 0 && value[0] != "" {
			searchParams[key] = value[0]
		}
	}

	customers, err := h.customerService.SearchCustomers(ctx, searchParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(customers)
}
