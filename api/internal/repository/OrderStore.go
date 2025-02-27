package repository

import (
	"bookstore/api/api/internal/model"
	"context"
)

type OrderStore interface {
	CreateOrder(ctx context.Context, order model.Order) (model.Order, error)
	GetOrder(ctx context.Context, id int) (model.Order, error)
	UpdateOrder(ctx context.Context, id int, order model.Order) (model.Order, error)
	DeleteOrder(ctx context.Context, id int) error
	SearchOrders(ctx context.Context, params map[string]string) ([]model.Order, error)
}
