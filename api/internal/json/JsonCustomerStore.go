package json

import (
	"bookstore/api/api/internal/model"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
)

type JsonCustomerStore struct {
	filename  string
	mutex     sync.RWMutex
	lastID    int
	customers []model.Customer
}

type CustomersData struct {
	Customers []model.Customer `json:"customers"`
}

func NewJsonCustomerStore() *JsonCustomerStore {
	store := &JsonCustomerStore{
		filename:  "../data/customers.json",
		customers: make([]model.Customer, 0),
	}

	if err := store.loadFromFile(); err != nil {
		panic(err)
	}

	return store
}

func (s *JsonCustomerStore) loadFromFile() error {
	if err := os.MkdirAll("../data", 0755); err != nil {
		return err
	}

	if _, err := os.Stat(s.filename); os.IsNotExist(err) {
		initialData := CustomersData{Customers: []model.Customer{}}
		data, _ := json.MarshalIndent(initialData, "", "  ")
		if err := os.WriteFile(s.filename, data, 0644); err != nil {
			return err
		}
	}

	data, err := os.ReadFile(s.filename)
	if err != nil {
		return err
	}

	var customersData CustomersData
	if err := json.Unmarshal(data, &customersData); err != nil {
		return err
	}

	s.customers = customersData.Customers

	for _, customer := range s.customers {
		if customer.ID > s.lastID {
			s.lastID = customer.ID
		}
	}
	return nil
}

func (s *JsonCustomerStore) SaveToFile() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	data, err := json.MarshalIndent(CustomersData{Customers: s.customers}, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.filename, data, 0644)
}

func (s *JsonCustomerStore) getNextID() int {
	s.lastID++
	return s.lastID
}

func (s *JsonCustomerStore) CreateCustomer(ctx context.Context, customer model.Customer) (model.Customer, error) {
	select {
	case <-ctx.Done():
		return model.Customer{}, ctx.Err()
	default:
		s.mutex.Lock()
		defer s.mutex.Unlock()

		customer.ID = s.getNextID()
		s.customers = append(s.customers, customer)
		return customer, nil
	}
}

func (s *JsonCustomerStore) GetCustomer(ctx context.Context, id int) (model.Customer, error) {
	select {
	case <-ctx.Done():
		return model.Customer{}, ctx.Err()
	default:
		s.mutex.RLock()
		defer s.mutex.RUnlock()

		for _, customer := range s.customers {
			if customer.ID == id {
				return customer, nil
			}
		}
		return model.Customer{}, fmt.Errorf("customer with id %d not found", id)
	}
}

func (s *JsonCustomerStore) UpdateCustomer(ctx context.Context, id int, updatedCustomer model.Customer) (model.Customer, error) {
	select {
	case <-ctx.Done():
		return model.Customer{}, ctx.Err()
	default:
		s.mutex.Lock()
		defer s.mutex.Unlock()

		for i, customer := range s.customers {
			if customer.ID == id {
				s.customers[i] = updatedCustomer
				return updatedCustomer, nil
			}
		}
		return model.Customer{}, fmt.Errorf("customer with id %d not found", id)
	}
}

func (s *JsonCustomerStore) DeleteCustomer(ctx context.Context, id int) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		s.mutex.Lock()
		defer s.mutex.Unlock()

		for i, customer := range s.customers {
			if customer.ID == id {
				s.customers = append(s.customers[:i], s.customers[i+1:]...)
				return nil
			}
		}
		return fmt.Errorf("customer with id %d not found", id)
	}
}

func (s *JsonCustomerStore) SearchCustomers(ctx context.Context, params map[string]string) ([]model.Customer, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		s.mutex.RLock()
		defer s.mutex.RUnlock()

		if params == nil {
			return s.customers, nil
		}

		result := []model.Customer{}
		for _, customer := range s.customers {
			matches := true
			for key, value := range params {
				switch key {
				case "name":
					if !strings.EqualFold(customer.Name, value) {
						matches = false
					}
				case "email":
					if !strings.EqualFold(customer.Email, value) {
						matches = false
					}
				}
			}
			if matches {
				result = append(result, customer)
			}
		}
		return result, nil
	}
}
