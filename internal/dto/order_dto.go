package dto

import (
	"time"

	"github.com/shopspring/decimal"
)

type CreateOrderRequest struct {
	ProductID       int64  `json:"productId" binding:"required"`
	Quantity        int    `json:"quantity" binding:"required,min=1"`
	CustomerAddress string `json:"customerAddress" binding:"required"`
}

type UpdateOrderStatusRequest struct {
	Status        string `json:"status" binding:"required,oneof=PLACED CONFIRMED COMPLETED CANCELLED"`
	PaymentTxHash string `json:"paymentTxHash,omitempty"`
}

type OrderResponse struct {
	ID              int64           `json:"id"`
	OrderID         string          `json:"orderId"`
	ProductID       int64           `json:"productId"`
	CustomerAddress string          `json:"customerAddress"`
	MerchantAddress string          `json:"merchantAddress"`
	Quantity        int             `json:"quantity"`
	UnitPrice       decimal.Decimal `json:"unitPrice"`
	TotalPrice      decimal.Decimal `json:"totalPrice"`
	Status          string          `json:"status"`
	PaymentTxHash   string          `json:"paymentTxHash,omitempty"`
	CreatedAt       time.Time       `json:"createdAt"`
	CompletedAt     *time.Time      `json:"completedAt,omitempty"`
}

type OrderListResponse struct {
	Orders []OrderResponse `json:"orders"`
	Total  int64           `json:"total"`
	Page   int             `json:"page"`
	Limit  int             `json:"limit"`
}