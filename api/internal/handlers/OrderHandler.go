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

type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

func (h *OrderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		h.CreateOrder(w, r)
	} else if r.Method == http.MethodGet {
		h.GetOrders(w, r)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors.Error{Message: "Request not allowed"})
	}
}

func (h *OrderHandler) ServeHTTPById(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPut {
		h.UpdateOrder(w, r)
	} else if r.Method == http.MethodDelete {
		h.DeleteOrder(w, r)
	} else if r.Method == http.MethodGet {
		h.GetOrder(w, r)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors.Error{Message: "Request not allowed"})
	}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	decoder := json.NewDecoder(r.Body)
	var orderInput model.OrderInput
	err := decoder.Decode(&orderInput)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors.Error{Message: "Invalid order payload"})
		return
	}

	order, err := h.orderService.CreateOrder(ctx, orderInput)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors.Error{Message: err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	order, err := h.orderService.GetOrder(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errors.Error{Message: "Order not found"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	params := r.URL.Query()
	searchParams := make(map[string]string)
	for key, value := range params {
		if len(value) > 0 && value[0] != "" {
			searchParams[key] = value[0]
		}
	}

	orders, err := h.orderService.SearchOrders(ctx, searchParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
}

func (h *OrderHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var orderInput model.OrderInput
	err = decoder.Decode(&orderInput)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors.Error{Message: "Invalid order payload"})
		return
	}

	order, err := h.orderService.UpdateOrder(ctx, id, orderInput)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors.Error{Message: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	err = h.orderService.DeleteOrder(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errors.Error{Message: "Order not found"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
