package dto

import (
	"time"

	"github.com/shopspring/decimal"
)

type CreateProductRequest struct {
	Name            string          `json:"name" binding:"required"`
	Description     string          `json:"description"`
	Category        string          `json:"category"`
	ImageURL        string          `json:"imageUrl"`
	OriginalPrice   decimal.Decimal `json:"originalPrice" binding:"required"`
	DiscountRate    int             `json:"discountRate" binding:"required,min=0,max=10000"`
	Quantity        int             `json:"quantity" binding:"required,min=1"`
	ExpiryTime      *time.Time      `json:"expiryTime"`
	SalesTarget     int             `json:"salesTarget,omitempty"`
	MerchantAddress string          `json:"merchantAddress" binding:"required"`
}

type UpdateProductRequest struct {
	Name         string          `json:"name,omitempty"`
	Description  string          `json:"description,omitempty"`
	Category     string          `json:"category,omitempty"`
	ImageURL     string          `json:"imageUrl,omitempty"`
	OriginalPrice decimal.Decimal `json:"originalPrice,omitempty"`
	DiscountRate int             `json:"discountRate,omitempty"`
	Quantity     int             `json:"quantity,omitempty"`
	ExpiryTime   *time.Time      `json:"expiryTime,omitempty"`
	Status       string          `json:"status,omitempty"`
}

type ProductResponse struct {
	ID                  int64           `json:"id"`
	ProductID           string          `json:"productId"`
	MerchantAddress     string          `json:"merchantAddress"`
	CustomerAddress     *string         `json:"customerAddress,omitempty"`
	Name                string          `json:"name"`
	Description         string          `json:"description"`
	Category            string          `json:"category"`
	ImageURL            string          `json:"imageUrl"`
	OriginalPrice       decimal.Decimal `json:"originalPrice"`
	DiscountedPrice     decimal.Decimal `json:"discountedPrice"`
	DiscountRate        int             `json:"discountRate"`
	DiscountPercentage  float64         `json:"discountPercentage"`
	Quantity            int             `json:"quantity"`
	AvailableQty        int             `json:"availableQty"`
	Status              string          `json:"status"`
	ExpiryTime          *time.Time      `json:"expiryTime,omitempty"`
	SalesTarget         int             `json:"salesTarget,omitempty"`
	SoldQuantity        int             `json:"soldQuantity"`
	PaymentStatus       string          `json:"paymentStatus"`
	CreatedAt           time.Time       `json:"createdAt"`
	UpdatedAt           time.Time       `json:"updatedAt"`
	ListedAt            time.Time       `json:"listedAt"`
	SoldOutAt           *time.Time      `json:"soldOutAt,omitempty"`
}

type ProductListResponse struct {
	Products []ProductResponse `json:"products"`
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	Limit    int               `json:"limit"`
}

type ProductFilters struct {
	Category    string          `form:"category"`
	MaxPrice    decimal.Decimal `form:"maxPrice"`
	MinDiscount int             `form:"minDiscount"`
	Status      string          `form:"status"`
	Merchant    string          `form:"merchant"`
	SortBy      string          `form:"sortBy"` // price, expiry, discount
	SortOrder   string          `form:"sortOrder"` // asc, desc
}