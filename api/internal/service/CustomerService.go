package service

import (
	"context"
	"errors"
	"time"

	"github.com/yassine22-alt/biblios-app/api/internal/model"
	"github.com/yassine22-alt/biblios-app/api/internal/repository"
)

type CustomerService struct {
	repo repository.CustomerStore
}

func NewCustomerService(repo repository.CustomerStore) *CustomerService {
	return &CustomerService{
		repo: repo,
	}
}

func (s *CustomerService) CreateCustomer(ctx context.Context, customerInput model.CustomerInput) (model.Customer, error) {
	if err := ctx.Err(); err != nil {
		return model.Customer{}, err
	}

	customer := model.Customer{
		Name:      customerInput.Name,
		Email:     customerInput.Email,
		Address:   customerInput.Address,
		CreatedAt: time.Now(),
	}

	if customer.Name == "" {
		return model.Customer{}, errors.New("customer name is mandatory")
	}
	if customer.Email == "" {
		return model.Customer{}, errors.New("customer email is mandatory")
	}

	return s.repo.CreateCustomer(ctx, customer)
}

func (s *CustomerService) GetCustomer(ctx context.Context, id int) (model.Customer, error) {
	if err := ctx.Err(); err != nil {
		return model.Customer{}, err
	}
	return s.repo.GetCustomer(ctx, id)
}

func (s *CustomerService) UpdateCustomer(ctx context.Context, id int, customerInput model.CustomerInput) (model.Customer, error) {
	if err := ctx.Err(); err != nil {
		return model.Customer{}, err
	}

	existingCustomer, err := s.repo.GetCustomer(ctx, id)
	if err != nil {
		return model.Customer{}, err
	}

	updatedCustomer := model.Customer{
		ID:        existingCustomer.ID,
		Name:      customerInput.Name,
		Email:     customerInput.Email,
		Address:   customerInput.Address,
		CreatedAt: existingCustomer.CreatedAt,
	}

	if updatedCustomer.Name == "" {
		return model.Customer{}, errors.New("customer name is mandatory")
	}
	if updatedCustomer.Email == "" {
		return model.Customer{}, errors.New("customer email is mandatory")
	}

	return s.repo.UpdateCustomer(ctx, id, updatedCustomer)
}

func (s *CustomerService) DeleteCustomer(ctx context.Context, id int) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return s.repo.DeleteCustomer(ctx, id)
}

func (s *CustomerService) SearchCustomers(ctx context.Context, params map[string]string) ([]model.Customer, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return s.repo.SearchCustomers(ctx, params)
}
