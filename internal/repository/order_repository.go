package repository

import (
	"database/sql"
	"fmt"
	"go-service/internal/model"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(order *model.Order) (int64, error) {
	query := `
		INSERT INTO orders (
			order_id, product_id, customer_address, merchant_address,
			quantity, unit_price, total_price, status, payment_tx_hash,
			created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, NOW()
		) RETURNING id
	`

	var id int64
	err := r.db.QueryRow(query,
		order.OrderID, order.ProductID, order.CustomerAddress, order.MerchantAddress,
		order.Quantity, order.UnitPrice, order.TotalPrice, order.Status, order.PaymentTxHash,
	).Scan(&id)

	return id, err
}

func (r *OrderRepository) FindByID(id int64) (*model.Order, error) {
	query := `
		SELECT id, order_id, product_id, customer_address, merchant_address,
			quantity, unit_price, total_price, status, payment_tx_hash,
			created_at, completed_at
		FROM orders WHERE id = $1
	`

	var order model.Order
	err := r.db.QueryRow(query, id).Scan(
		&order.ID, &order.OrderID, &order.ProductID, &order.CustomerAddress, &order.MerchantAddress,
		&order.Quantity, &order.UnitPrice, &order.TotalPrice, &order.Status, &order.PaymentTxHash,
		&order.CreatedAt, &order.CompletedAt,
	)

	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *OrderRepository) FindByOrderID(orderID string) (*model.Order, error) {
	// Try integer id first, then fall back to UUID order_id.
	query := `
		SELECT id, order_id, product_id, customer_address, merchant_address,
			quantity, unit_price, total_price, status, payment_tx_hash,
			created_at, completed_at
		FROM orders WHERE id = $1::bigint
	`

	var order model.Order
	err := r.db.QueryRow(query, orderID).Scan(
		&order.ID, &order.OrderID, &order.ProductID, &order.CustomerAddress, &order.MerchantAddress,
		&order.Quantity, &order.UnitPrice, &order.TotalPrice, &order.Status, &order.PaymentTxHash,
		&order.CreatedAt, &order.CompletedAt,
	)
	if err == nil {
		return &order, nil
	}

	// Fallback: treat as UUID order_id
	query2 := `
		SELECT id, order_id, product_id, customer_address, merchant_address,
			quantity, unit_price, total_price, status, payment_tx_hash,
			created_at, completed_at
		FROM orders WHERE order_id = $1
	`
	err = r.db.QueryRow(query2, orderID).Scan(
		&order.ID, &order.OrderID, &order.ProductID, &order.CustomerAddress, &order.MerchantAddress,
		&order.Quantity, &order.UnitPrice, &order.TotalPrice, &order.Status, &order.PaymentTxHash,
		&order.CreatedAt, &order.CompletedAt,
	)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) FindByCustomer(customerAddress string, limit, offset int) ([]model.Order, int64, error) {
	// Count query
	var total int64
	err := r.db.QueryRow("SELECT COUNT(*) FROM orders WHERE customer_address = $1", customerAddress).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Data query
	query := `
		SELECT id, order_id, product_id, customer_address, merchant_address,
			quantity, unit_price, total_price, status, payment_tx_hash,
			created_at, completed_at
		FROM orders
		WHERE customer_address = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, customerAddress, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	orders := make([]model.Order, 0)
	for rows.Next() {
		var order model.Order
		err := rows.Scan(
			&order.ID, &order.OrderID, &order.ProductID, &order.CustomerAddress, &order.MerchantAddress,
			&order.Quantity, &order.UnitPrice, &order.TotalPrice, &order.Status, &order.PaymentTxHash,
			&order.CreatedAt, &order.CompletedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		orders = append(orders, order)
	}

	return orders, total, nil
}

func (r *OrderRepository) FindByMerchant(merchantAddress string, limit, offset int) ([]model.Order, int64, error) {
	// Count query
	var total int64
	err := r.db.QueryRow("SELECT COUNT(*) FROM orders WHERE merchant_address = $1", merchantAddress).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Data query
	query := `
		SELECT id, order_id, product_id, customer_address, merchant_address,
			quantity, unit_price, total_price, status, payment_tx_hash,
			created_at, completed_at
		FROM orders
		WHERE merchant_address = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, merchantAddress, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	orders := make([]model.Order, 0)
	for rows.Next() {
		var order model.Order
		err := rows.Scan(
			&order.ID, &order.OrderID, &order.ProductID, &order.CustomerAddress, &order.MerchantAddress,
			&order.Quantity, &order.UnitPrice, &order.TotalPrice, &order.Status, &order.PaymentTxHash,
			&order.CreatedAt, &order.CompletedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		orders = append(orders, order)
	}

	return orders, total, nil
}

func (r *OrderRepository) UpdateStatus(orderID string, status string, paymentTxHash string, completedAt *string) error {
	// orderID may be an integer id (from REST path) or a UUID order_id (from internal logic).
	// Try integer id first; fall back to UUID order_id.
	query := `
		UPDATE orders
		SET status = $2, payment_tx_hash = $3, completed_at = $4
		WHERE id = $1::bigint
	`
	res, err := r.db.Exec(query, orderID, status, paymentTxHash, completedAt)
	if err != nil {
		// Fallback: treat as UUID order_id
		query2 := `
			UPDATE orders
			SET status = $2, payment_tx_hash = $3, completed_at = $4
			WHERE order_id = $1
		`
		_, err2 := r.db.Exec(query2, orderID, status, paymentTxHash, completedAt)
		return err2
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		// No row matched integer id — try UUID order_id
		query2 := `
			UPDATE orders
			SET status = $2, payment_tx_hash = $3, completed_at = $4
			WHERE order_id = $1
		`
		_, err2 := r.db.Exec(query2, orderID, status, paymentTxHash, completedAt)
		return err2
	}
	return nil
}

func (r *OrderRepository) FindByStatusAndAge(status string, maxAgeHours int) ([]model.Order, error) {
	query := `
		SELECT id, order_id, product_id, customer_address, merchant_address,
			quantity, unit_price, total_price, status, payment_tx_hash,
			created_at, completed_at
		FROM orders
		WHERE status = $1 AND created_at < NOW() - INTERVAL '%d hours'
	`

	query = fmt.Sprintf(query, maxAgeHours)

	rows, err := r.db.Query(query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make([]model.Order, 0)
	for rows.Next() {
		var order model.Order
		err := rows.Scan(
			&order.ID, &order.OrderID, &order.ProductID, &order.CustomerAddress, &order.MerchantAddress,
			&order.Quantity, &order.UnitPrice, &order.TotalPrice, &order.Status, &order.PaymentTxHash,
			&order.CreatedAt, &order.CompletedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}