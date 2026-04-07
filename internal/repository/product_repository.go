package repository

import (
	"database/sql"
	"fmt"
	"go-service/internal/model"
	"strings"

	"github.com/shopspring/decimal"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(product *model.Product) (int64, error) {
	query := `
		INSERT INTO products (
			product_id, merchant_address, customer_address,
			name, description, category, image_url,
			original_price, discounted_price, discount_rate,
			quantity, available_qty, status, expiry_time,
			sales_target, sold_quantity,
			chain_id, vault_address, contract_product_id, payment_status,
			created_at, updated_at, listed_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
			NOW(), NOW(), NOW()
		) RETURNING id
	`

	var id int64
	err := r.db.QueryRow(query,
		product.ProductID, product.MerchantAddress, product.CustomerAddress,
		product.Name, product.Description, product.Category, product.ImageURL,
		product.OriginalPrice, product.DiscountedPrice, product.DiscountRate,
		product.Quantity, product.AvailableQty, product.Status, product.ExpiryTime,
		product.SalesTarget, product.SoldQuantity,
		product.ChainID, product.VaultAddress, product.ContractProductID, product.PaymentStatus,
	).Scan(&id)

	return id, err
}

func (r *ProductRepository) FindByID(id int64) (*model.Product, error) {
	query := `
		SELECT id, product_id, merchant_address, customer_address,
			name, description, category, image_url,
			original_price, discounted_price, discount_rate,
			quantity, available_qty, status, expiry_time,
			sales_target, sold_quantity,
			chain_id, vault_address, contract_product_id, payment_status,
			created_at, updated_at, listed_at, sold_out_at
		FROM products WHERE id = $1
	`

	var product model.Product
	err := r.db.QueryRow(query, id).Scan(
		&product.ID, &product.ProductID, &product.MerchantAddress, &product.CustomerAddress,
		&product.Name, &product.Description, &product.Category, &product.ImageURL,
		&product.OriginalPrice, &product.DiscountedPrice, &product.DiscountRate,
		&product.Quantity, &product.AvailableQty, &product.Status, &product.ExpiryTime,
		&product.SalesTarget, &product.SoldQuantity,
		&product.ChainID, &product.VaultAddress, &product.ContractProductID, &product.PaymentStatus,
		&product.CreatedAt, &product.UpdatedAt, &product.ListedAt, &product.SoldOutAt,
	)

	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepository) FindAvailable(filters map[string]interface{}, limit, offset int) ([]model.Product, int64, error) {
	whereConditions := []string{"status = 'AVAILABLE'"}
	args := []interface{}{}
	argCount := 0

	if category, ok := filters["category"].(string); ok && category != "" {
		argCount++
		whereConditions = append(whereConditions, fmt.Sprintf("category = $%d", argCount))
		args = append(args, category)
	}

	if maxPrice, ok := filters["maxPrice"].(decimal.Decimal); ok && !maxPrice.IsZero() {
		argCount++
		whereConditions = append(whereConditions, fmt.Sprintf("discounted_price <= $%d", argCount))
		args = append(args, maxPrice)
	}

	if minDiscount, ok := filters["minDiscount"].(int); ok && minDiscount > 0 {
		argCount++
		whereConditions = append(whereConditions, fmt.Sprintf("discount_rate >= $%d", argCount))
		args = append(args, minDiscount)
	}

	if merchant, ok := filters["merchant"].(string); ok && merchant != "" {
		argCount++
		whereConditions = append(whereConditions, fmt.Sprintf("merchant_address = $%d", argCount))
		args = append(args, merchant)
	}

	whereClause := strings.Join(whereConditions, " AND ")

	sortBy := "listed_at DESC"
	if sort, ok := filters["sortBy"].(string); ok {
		switch sort {
		case "price":
			sortBy = "discounted_price ASC"
		case "expiry":
			sortBy = "expiry_time ASC NULLS LAST"
		case "discount":
			sortBy = "discount_rate DESC"
		}
	}

	if sortOrder, ok := filters["sortOrder"].(string); ok && sortOrder == "desc" {
		sortBy = strings.Replace(sortBy, " ASC", " DESC", 1)
		sortBy = strings.Replace(sortBy, " DESC NULLS LAST", " DESC NULLS FIRST", 1)
	}

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM products WHERE %s", whereClause)
	var total int64
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Data query
	query := fmt.Sprintf(`
		SELECT id, product_id, merchant_address, customer_address,
			name, description, category, image_url,
			original_price, discounted_price, discount_rate,
			quantity, available_qty, status, expiry_time,
			sales_target, sold_quantity,
			chain_id, vault_address, contract_product_id, payment_status,
			created_at, updated_at, listed_at, sold_out_at
		FROM products
		WHERE %s
		ORDER BY %s
		LIMIT $%d OFFSET $%d
	`, whereClause, sortBy, argCount+1, argCount+2)

	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	products := make([]model.Product, 0)
	for rows.Next() {
		var product model.Product
		err := rows.Scan(
			&product.ID, &product.ProductID, &product.MerchantAddress, &product.CustomerAddress,
			&product.Name, &product.Description, &product.Category, &product.ImageURL,
			&product.OriginalPrice, &product.DiscountedPrice, &product.DiscountRate,
			&product.Quantity, &product.AvailableQty, &product.Status, &product.ExpiryTime,
			&product.SalesTarget, &product.SoldQuantity,
			&product.ChainID, &product.VaultAddress, &product.ContractProductID, &product.PaymentStatus,
			&product.CreatedAt, &product.UpdatedAt, &product.ListedAt, &product.SoldOutAt,
		)
		if err != nil {
			return nil, 0, err
		}
		products = append(products, product)
	}

	return products, total, nil
}

func (r *ProductRepository) Update(product *model.Product) error {
	query := `
		UPDATE products SET
			customer_address = $2, name = $3, description = $4, category = $5, image_url = $6,
			original_price = $7, discounted_price = $8, discount_rate = $9,
			quantity = $10, available_qty = $11, status = $12, expiry_time = $13,
			sales_target = $14, sold_quantity = $15,
			chain_id = $16, vault_address = $17, contract_product_id = $18, payment_status = $19,
			updated_at = NOW(), sold_out_at = $20
		WHERE id = $1
	`

	_, err := r.db.Exec(query,
		product.ID, product.CustomerAddress, product.Name, product.Description, product.Category, product.ImageURL,
		product.OriginalPrice, product.DiscountedPrice, product.DiscountRate,
		product.Quantity, product.AvailableQty, product.Status, product.ExpiryTime,
		product.SalesTarget, product.SoldQuantity,
		product.ChainID, product.VaultAddress, product.ContractProductID, product.PaymentStatus,
		product.SoldOutAt,
	)

	return err
}

func (r *ProductRepository) ReserveQuantity(productID int64, quantity int) error {
	query := `
		UPDATE products
		SET available_qty = available_qty - $2, sold_quantity = sold_quantity + $2, updated_at = NOW()
		WHERE id = $1 AND available_qty >= $2 AND status = 'AVAILABLE'
	`

	result, err := r.db.Exec(query, productID, quantity)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("insufficient stock or product not available")
	}

	return nil
}

func (r *ProductRepository) FindExpired() ([]model.Product, error) {
	query := `
		SELECT id, product_id, merchant_address
		FROM products
		WHERE status = 'AVAILABLE' AND expiry_time < NOW()
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]model.Product, 0)
	for rows.Next() {
		var product model.Product
		err := rows.Scan(&product.ID, &product.ProductID, &product.MerchantAddress)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (r *ProductRepository) UpdateExpiredStatus(productIDs []int64) error {
	if len(productIDs) == 0 {
		return nil
	}

	placeholders := make([]string, len(productIDs))
	args := make([]interface{}, len(productIDs))
	for i, id := range productIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf(`
		UPDATE products
		SET status = 'EXPIRED', available_qty = 0, updated_at = NOW()
		WHERE id IN (%s) AND status = 'AVAILABLE'
	`, strings.Join(placeholders, ","))

	_, err := r.db.Exec(query, args...)
	return err
}