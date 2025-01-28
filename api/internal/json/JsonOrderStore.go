package json

import (
	"bookstore/api/api/internal/model"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"
)

type JsonOrderStore struct {
	filename string
	mutex    sync.RWMutex
	lastID   int
	orders   []model.Order
}

type OrdersData struct {
	Orders []model.Order `json:"orders"`
}

func NewJsonOrderStore() *JsonOrderStore {
	store := &JsonOrderStore{
		filename: "../data/orders.json",
		orders:   make([]model.Order, 0),
	}

	if err := store.loadFromFile(); err != nil {
		panic(err)
	}

	return store
}

func (s *JsonOrderStore) loadFromFile() error {
	if err := os.MkdirAll("../data", 0755); err != nil {
		return err
	}

	if _, err := os.Stat(s.filename); os.IsNotExist(err) {
		initialData := OrdersData{Orders: []model.Order{}}
		data, _ := json.MarshalIndent(initialData, "", "  ")
		if err := os.WriteFile(s.filename, data, 0644); err != nil {
			return err
		}
	}

	data, err := os.ReadFile(s.filename)
	if err != nil {
		return err
	}

	var ordersData OrdersData
	if err := json.Unmarshal(data, &ordersData); err != nil {
		return err
	}

	s.orders = ordersData.Orders

	for _, order := range s.orders {
		if order.ID > s.lastID {
			s.lastID = order.ID
		}
	}

	return nil
}

func (s *JsonOrderStore) SaveToFile() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	data, err := json.MarshalIndent(OrdersData{Orders: s.orders}, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.filename, data, 0644)
}

func (s *JsonOrderStore) getNextID() int {
	s.lastID++
	return s.lastID
}

func (s *JsonOrderStore) CreateOrder(ctx context.Context, order model.Order) (model.Order, error) {
	select {
	case <-ctx.Done():
		return model.Order{}, ctx.Err()
	default:
		s.mutex.Lock()
		defer s.mutex.Unlock()

		order.ID = s.getNextID()

		s.orders = append(s.orders, order)
		return order, nil
	}
}

func (s *JsonOrderStore) GetOrder(ctx context.Context, id int) (model.Order, error) {
	select {
	case <-ctx.Done():
		return model.Order{}, ctx.Err()
	default:
		s.mutex.RLock()
		defer s.mutex.RUnlock()

		for _, order := range s.orders {
			if order.ID == id {
				return order, nil
			}
		}
		return model.Order{}, fmt.Errorf("order with id %d not found", id)
	}
}

func (s *JsonOrderStore) UpdateOrder(ctx context.Context, id int, updatedOrder model.Order) (model.Order, error) {
	select {
	case <-ctx.Done():
		return model.Order{}, ctx.Err()
	default:
		s.mutex.Lock()
		defer s.mutex.Unlock()

		for i, order := range s.orders {
			if order.ID == id {
				s.orders[i] = updatedOrder
				return updatedOrder, nil
			}
		}
		return model.Order{}, fmt.Errorf("order with id %d not found", id)
	}
}

func (s *JsonOrderStore) DeleteOrder(ctx context.Context, id int) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		s.mutex.Lock()
		defer s.mutex.Unlock()

		for i, order := range s.orders {
			if order.ID == id {
				s.orders = append(s.orders[:i], s.orders[i+1:]...)
				return nil
			}
		}
		return fmt.Errorf("order with id %d not found", id)
	}
}

func (s *JsonOrderStore) SearchOrders(ctx context.Context, params map[string]string) ([]model.Order, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		s.mutex.RLock()
		defer s.mutex.RUnlock()

		if params == nil {
			return s.orders, nil
		}

		result := []model.Order{}
		for _, order := range s.orders {
			matches := true
			for key, value := range params {
				switch key {
				case "customer_id":
					if strconv.Itoa(order.CustomerId) != value {
						matches = false
					}
				case "status":
					if order.Status != value {
						matches = false
					}
				}
			}
			if matches {
				result = append(result, order)
			}
		}
		return result, nil
	}
}
