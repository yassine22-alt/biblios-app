package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/yassine22-alt/biblios-app/api/internal/model"
)

type PostgresCustomerStore struct {
	db *pgx.Conn
}

func NewPostgresCustomerStore(db *pgx.Conn) *PostgresCustomerStore {
	return &PostgresCustomerStore{db: db}
}

func (s *PostgresCustomerStore) CreateCustomer(ctx context.Context, customer model.Customer) (model.Customer, error) {
	query := `
    INSERT INTO customers (name, email, created_at, street, city, state, postal_code, country)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    RETURNING id, name, email, created_at, street, city, state, postal_code, country`

	err := s.db.QueryRow(ctx, query,
		customer.Name,
		customer.Email,
		customer.CreatedAt,
		customer.Address.Street,
		customer.Address.City,
		customer.Address.State,
		customer.Address.PostalCode,
		customer.Address.Country,
	).Scan(
		&customer.ID,
		&customer.Name,
		&customer.Email,
		&customer.CreatedAt,
		&customer.Address.Street,
		&customer.Address.City,
		&customer.Address.State,
		&customer.Address.PostalCode,
		&customer.Address.Country,
	)

	if err != nil {
		return model.Customer{}, fmt.Errorf("error creating customer: %v", err)
	}

	return customer, nil
}

func (s *PostgresCustomerStore) GetCustomer(ctx context.Context, id int) (model.Customer, error) {
	query := `
    SELECT id, name, email, created_at, street, city, state, postal_code, country
    FROM customers 
    WHERE id = $1`

	var customer model.Customer
	err := s.db.QueryRow(ctx, query, id).Scan(
		&customer.ID,
		&customer.Name,
		&customer.Email,
		&customer.CreatedAt,
		&customer.Address.Street,
		&customer.Address.City,
		&customer.Address.State,
		&customer.Address.PostalCode,
		&customer.Address.Country,
	)

	if err != nil {
		return model.Customer{}, fmt.Errorf("error getting customer: %v", err)
	}

	return customer, nil
}

func (s *PostgresCustomerStore) UpdateCustomer(ctx context.Context, id int, customer model.Customer) (model.Customer, error) {
	query := `
    UPDATE customers 
    SET name = $1, email = $2, street = $3, city = $4, state = $5, postal_code = $6, country = $7
    WHERE id = $8
    RETURNING id, name, email, created_at, street, city, state, postal_code, country`

	err := s.db.QueryRow(ctx, query,
		customer.Name,
		customer.Email,
		customer.Address.Street,
		customer.Address.City,
		customer.Address.State,
		customer.Address.PostalCode,
		customer.Address.Country,
		id,
	).Scan(
		&customer.ID,
		&customer.Name,
		&customer.Email,
		&customer.CreatedAt,
		&customer.Address.Street,
		&customer.Address.City,
		&customer.Address.State,
		&customer.Address.PostalCode,
		&customer.Address.Country,
	)

	if err != nil {
		return model.Customer{}, fmt.Errorf("error updating customer: %v", err)
	}

	return customer, nil
}

func (s *PostgresCustomerStore) DeleteCustomer(ctx context.Context, id int) error {
	query := `DELETE FROM customers WHERE id = $1`
	_, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting customer: %v", err)
	}
	return nil
}

func (s *PostgresCustomerStore) SearchCustomers(ctx context.Context, params map[string]string) ([]model.Customer, error) {
	query := `
    SELECT id, name, email, created_at, street, city, state, postal_code, country 
    FROM customers`

	var conditions []string
	var args []interface{}
	argPosition := 1

	if name, ok := params["name"]; ok {
		conditions = append(conditions, fmt.Sprintf("name ILIKE $%d", argPosition))
		args = append(args, "%"+name+"%")
		argPosition++
	}

	if email, ok := params["email"]; ok {
		conditions = append(conditions, fmt.Sprintf("email = $%d", argPosition))
		args = append(args, email)
		argPosition++
	}

	if len(conditions) > 0 {
		query += " WHERE " + conditions[0]
		for i := 1; i < len(conditions); i++ {
			query += " AND " + conditions[i]
		}
	}

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error searching customers: %v", err)
	}
	defer rows.Close()

	var customers []model.Customer
	for rows.Next() {
		var customer model.Customer
		err := rows.Scan(
			&customer.ID,
			&customer.Name,
			&customer.Email,
			&customer.CreatedAt,
			&customer.Address.Street,
			&customer.Address.City,
			&customer.Address.State,
			&customer.Address.PostalCode,
			&customer.Address.Country,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning customer row: %v", err)
		}
		customers = append(customers, customer)
	}

	return customers, nil
}
