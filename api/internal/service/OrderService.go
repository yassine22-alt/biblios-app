package service

import (
	"bookstore/api/api/internal/model"
	"bookstore/api/api/internal/repository"
	"context"
	"errors"
	"time"
)

type OrderService struct {
	repo         repository.OrderStore
	repoCustomer repository.CustomerStore
	repoBook     repository.BookStore
	currentID    int
}

func NewOrderService(repo repository.OrderStore, repoCustomer repository.CustomerStore, repoBook repository.BookStore) *OrderService {
	return &OrderService{
		repo:         repo,
		repoCustomer: repoCustomer,
		repoBook:     repoBook,
		currentID:    1,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, orderInput model.OrderInput) (model.Order, error) {
	if err := ctx.Err(); err != nil {
		return model.Order{}, err
	}

	totalPrice, err := s.calculateTotalPrice(ctx, orderInput.Items)
	if err != nil {
		return model.Order{}, err
	}

	order := model.Order{
		ID:         s.currentID,
		CustomerId: orderInput.CustomerId,
		Items:      orderInput.Items,
		TotalPrice: totalPrice,
		CreatedAt:  time.Now(),
		Status:     "Pending",
	}
	s.currentID++

	if order.CustomerId == 0 {
		return model.Order{}, errors.New("customer ID is mandatory")
	}
	if len(order.Items) == 0 {
		return model.Order{}, errors.New("order must have at least one item")
	}

	_, err2 := s.repoCustomer.GetCustomer(ctx, order.CustomerId)

	if err2 != nil {
		return model.Order{}, errors.New("customer non existant")
	}

	for _, item := range order.Items {
		_, err := s.repoBook.GetBook(ctx, item.BookID)
		if err != nil {
			return model.Order{}, errors.New("book non existant")
		}
	}

	return s.repo.CreateOrder(ctx, order)
}

func (s *OrderService) GetOrder(ctx context.Context, id int) (model.Order, error) {
	if err := ctx.Err(); err != nil {
		return model.Order{}, err
	}
	return s.repo.GetOrder(ctx, id)
}

func (s *OrderService) UpdateOrder(ctx context.Context, id int, orderInput model.OrderInput) (model.Order, error) {
	if err := ctx.Err(); err != nil {
		return model.Order{}, err
	}

	existingOrder, err := s.repo.GetOrder(ctx, id)
	if err != nil {
		return model.Order{}, err
	}

	totalPrice, err := s.calculateTotalPrice(ctx, orderInput.Items)
	if err != nil {
		return model.Order{}, err
	}

	updatedOrder := model.Order{
		ID:         existingOrder.ID,
		CustomerId: orderInput.CustomerId,
		Items:      orderInput.Items,
		TotalPrice: totalPrice,
		CreatedAt:  existingOrder.CreatedAt,
		Status:     existingOrder.Status,
	}

	if updatedOrder.CustomerId == 0 {
		return model.Order{}, errors.New("customer ID is mandatory")
	}
	if len(updatedOrder.Items) == 0 {
		return model.Order{}, errors.New("order must have at least one item")
	}

	return s.repo.UpdateOrder(ctx, id, updatedOrder)
}

func (s *OrderService) DeleteOrder(ctx context.Context, id int) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return s.repo.DeleteOrder(ctx, id)
}

func (s *OrderService) SearchOrders(ctx context.Context, params map[string]string) ([]model.Order, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return s.repo.SearchOrders(ctx, params)
}

func (s *OrderService) calculateTotalPrice(ctx context.Context, items []model.OrderItem) (float64, error) {
	var total float64
	for _, item := range items {
		book, err := s.repoBook.GetBook(ctx, item.BookID)
		if err != nil {
			return 0, err
		}
		total += book.Price * float64(item.Quantity)
	}
	return total, nil
}
