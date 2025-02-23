package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/yassine22-alt/biblios-app/api/internal/model"
)

type PostgresOrderStore struct {
	db *pgx.Conn
}

func NewPostgresOrderStore(db *pgx.Conn) *PostgresOrderStore {
	return &PostgresOrderStore{db: db}
}

func (s *PostgresOrderStore) CreateOrder(ctx context.Context, order model.Order) (model.Order, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return model.Order{}, fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	// Insert order
	query := `
    INSERT INTO orders (customer_id, total_price, created_at, status)
    VALUES ($1, $2, $3, $4)
    RETURNING id, customer_id, total_price, created_at, status`

	err = tx.QueryRow(ctx, query,
		order.CustomerId,
		order.TotalPrice,
		order.CreatedAt,
		order.Status,
	).Scan(
		&order.ID,
		&order.CustomerId,
		&order.TotalPrice,
		&order.CreatedAt,
		&order.Status,
	)
	if err != nil {
		return model.Order{}, fmt.Errorf("error creating order: %v", err)
	}

	// Insert order items
	for _, item := range order.Items {
		query = `
        INSERT INTO order_items (order_id, book_id, quantity)
        VALUES ($1, $2, $3)`

		_, err = tx.Exec(ctx, query, order.ID, item.BookID, item.Quantity)
		if err != nil {
			return model.Order{}, fmt.Errorf("error creating order item: %v", err)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return model.Order{}, fmt.Errorf("error committing transaction: %v", err)
	}

	return order, nil
}

func (s *PostgresOrderStore) GetOrder(ctx context.Context, id int) (model.Order, error) {
	orderQuery := `
    SELECT id, customer_id, total_price, created_at, status
    FROM orders 
    WHERE id = $1`

	var order model.Order
	err := s.db.QueryRow(ctx, orderQuery, id).Scan(
		&order.ID,
		&order.CustomerId,
		&order.TotalPrice,
		&order.CreatedAt,
		&order.Status,
	)
	if err != nil {
		return model.Order{}, fmt.Errorf("error getting order: %v", err)
	}

	// Get order items
	itemsQuery := `
    SELECT book_id, quantity
    FROM order_items
    WHERE order_id = $1`

	rows, err := s.db.Query(ctx, itemsQuery, id)
	if err != nil {
		return model.Order{}, fmt.Errorf("error getting order items: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item model.OrderItem
		err := rows.Scan(&item.BookID, &item.Quantity)
		if err != nil {
			return model.Order{}, fmt.Errorf("error scanning order item: %v", err)
		}
		order.Items = append(order.Items, item)
	}

	return order, nil
}

func (s *PostgresOrderStore) UpdateOrder(ctx context.Context, id int, order model.Order) (model.Order, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return model.Order{}, fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	// Update order
	query := `
    UPDATE orders 
    SET customer_id = $1, total_price = $2, status = $3
    WHERE id = $4
    RETURNING id, customer_id, total_price, created_at, status`

	err = tx.QueryRow(ctx, query,
		order.CustomerId,
		order.TotalPrice,
		order.Status,
		id,
	).Scan(
		&order.ID,
		&order.CustomerId,
		&order.TotalPrice,
		&order.CreatedAt,
		&order.Status,
	)
	if err != nil {
		return model.Order{}, fmt.Errorf("error updating order: %v", err)
	}

	// Delete existing items
	_, err = tx.Exec(ctx, "DELETE FROM order_items WHERE order_id = $1", id)
	if err != nil {
		return model.Order{}, fmt.Errorf("error deleting order items: %v", err)
	}

	// Insert new items
	for _, item := range order.Items {
		_, err = tx.Exec(ctx,
			"INSERT INTO order_items (order_id, book_id, quantity) VALUES ($1, $2, $3)",
			id, item.BookID, item.Quantity)
		if err != nil {
			return model.Order{}, fmt.Errorf("error inserting order item: %v", err)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return model.Order{}, fmt.Errorf("error committing transaction: %v", err)
	}

	return order, nil
}

func (s *PostgresOrderStore) DeleteOrder(ctx context.Context, id int) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	// Delete order items first (due to foreign key constraint)
	_, err = tx.Exec(ctx, "DELETE FROM order_items WHERE order_id = $1", id)
	if err != nil {
		return fmt.Errorf("error deleting order items: %v", err)
	}

	// Delete order
	_, err = tx.Exec(ctx, "DELETE FROM orders WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("error deleting order: %v", err)
	}

	return tx.Commit(ctx)
}

func (s *PostgresOrderStore) SearchOrders(ctx context.Context, params map[string]string) ([]model.Order, error) {
	query := `
    SELECT id, customer_id, total_price, created_at, status 
    FROM orders`

	var conditions []string
	var args []interface{}
	argPosition := 1

	if CustomerId, ok := params["customer_id"]; ok {
		conditions = append(conditions, fmt.Sprintf("customer_id = $%d", argPosition))
		args = append(args, CustomerId)
		argPosition++
	}

	if status, ok := params["status"]; ok {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argPosition))
		args = append(args, status)
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
		return nil, fmt.Errorf("error searching orders: %v", err)
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var order model.Order
		err := rows.Scan(
			&order.ID,
			&order.CustomerId,
			&order.TotalPrice,
			&order.CreatedAt,
			&order.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning order row: %v", err)
		}
		orders = append(orders, order)
	}

	// Get items for each order
	for i := range orders {
		itemRows, err := s.db.Query(ctx,
			"SELECT book_id, quantity FROM order_items WHERE order_id = $1",
			orders[i].ID)
		if err != nil {
			return nil, fmt.Errorf("error getting order items: %v", err)
		}
		defer itemRows.Close()

		for itemRows.Next() {
			var item model.OrderItem
			if err := itemRows.Scan(&item.BookID, &item.Quantity); err != nil {
				return nil, fmt.Errorf("error scanning order item: %v", err)
			}
			orders[i].Items = append(orders[i].Items, item)
		}
	}

	return orders, nil
}
