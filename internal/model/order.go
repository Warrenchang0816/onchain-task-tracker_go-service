package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Order struct {
	ID                  int64
	OrderID             string           // UUID
	ProductID           int64
	CustomerAddress     string
	MerchantAddress     string

	Quantity            int
	UnitPrice           decimal.Decimal
	TotalPrice          decimal.Decimal

	Status              string           // PLACED|CONFIRMED|COMPLETED|CANCELLED
	PaymentTxHash       string

	CreatedAt           time.Time
	CompletedAt         *time.Time
}