package repository

import (
	"bookstore/api/api/internal/model"
	"context"
)

type CustomerStore interface {
	CreateCustomer(ctx context.Context, customer model.Customer) (model.Customer, error)
	GetCustomer(ctx context.Context, id int) (model.Customer, error)
	UpdateCustomer(ctx context.Context, id int, Customer model.Customer) (model.Customer, error)
	DeleteCustomer(ctx context.Context, id int) error
	SearchCustomers(ctx context.Context, params map[string]string) ([]model.Customer, error)
}
