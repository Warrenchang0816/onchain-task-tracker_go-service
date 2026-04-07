package service

import (
	"context"
	"fmt"
	"go-service/internal/dto"
	"go-service/internal/model"
	"go-service/internal/repository"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ProductService struct {
	productRepo *repository.ProductRepository
}

func NewProductService(productRepo *repository.ProductRepository) *ProductService {
	return &ProductService{productRepo: productRepo}
}

func (s *ProductService) CreateProduct(ctx context.Context, req dto.CreateProductRequest) (*model.Product, error) {
	// 计算折扣价
	discountedPrice := req.OriginalPrice
	if req.DiscountRate > 0 {
		discount := decimal.NewFromInt(int64(req.DiscountRate)).Div(decimal.NewFromInt(10000))
		discountedPrice = req.OriginalPrice.Mul(decimal.NewFromInt(1).Sub(discount))
	}

	product := &model.Product{
		ProductID:       uuid.New().String(),
		MerchantAddress: req.MerchantAddress,
		Name:            req.Name,
		Description:     req.Description,
		Category:        req.Category,
		ImageURL:        req.ImageURL,
		OriginalPrice:   req.OriginalPrice,
		DiscountedPrice: discountedPrice,
		DiscountRate:    req.DiscountRate,
		Quantity:        req.Quantity,
		AvailableQty:    req.Quantity,
		Status:          "AVAILABLE",
		ExpiryTime:      req.ExpiryTime,
		SalesTarget:     req.SalesTarget,
		SoldQuantity:    0,
		PaymentStatus:   "NOT_FUNDED",
	}

	id, err := s.productRepo.Create(product)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	product.ID = id
	return product, nil
}

func (s *ProductService) GetProduct(ctx context.Context, id int64) (*model.Product, error) {
	return s.productRepo.FindByID(id)
}

func (s *ProductService) ListAvailableProducts(ctx context.Context, filters dto.ProductFilters, page, limit int) ([]model.Product, int64, error) {
	filterMap := make(map[string]interface{})

	if filters.Category != "" {
		filterMap["category"] = filters.Category
	}
	if !filters.MaxPrice.IsZero() {
		filterMap["maxPrice"] = filters.MaxPrice
	}
	if filters.MinDiscount > 0 {
		filterMap["minDiscount"] = filters.MinDiscount
	}
	if filters.Merchant != "" {
		filterMap["merchant"] = filters.Merchant
	}
	if filters.SortBy != "" {
		filterMap["sortBy"] = filters.SortBy
	}
	if filters.SortOrder != "" {
		filterMap["sortOrder"] = filters.SortOrder
	}

	offset := (page - 1) * limit
	return s.productRepo.FindAvailable(filterMap, limit, offset)
}

func (s *ProductService) UpdateProduct(ctx context.Context, id int64, req dto.UpdateProductRequest) error {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("failed to find product: %w", err)
	}

	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Category != "" {
		product.Category = req.Category
	}
	if req.ImageURL != "" {
		product.ImageURL = req.ImageURL
	}
	if !req.OriginalPrice.IsZero() {
		product.OriginalPrice = req.OriginalPrice
		// 重新计算折扣价
		if product.DiscountRate > 0 {
			discount := decimal.NewFromInt(int64(product.DiscountRate)).Div(decimal.NewFromInt(10000))
			product.DiscountedPrice = req.OriginalPrice.Mul(decimal.NewFromInt(1).Sub(discount))
		} else {
			product.DiscountedPrice = req.OriginalPrice
		}
	}
	if req.DiscountRate > 0 {
		product.DiscountRate = req.DiscountRate
		// 重新计算折扣价
		discount := decimal.NewFromInt(int64(req.DiscountRate)).Div(decimal.NewFromInt(10000))
		product.DiscountedPrice = product.OriginalPrice.Mul(decimal.NewFromInt(1).Sub(discount))
	}
	if req.Quantity > 0 {
		product.Quantity = req.Quantity
		product.AvailableQty = req.Quantity - product.SoldQuantity
	}
	if req.ExpiryTime != nil {
		product.ExpiryTime = req.ExpiryTime
	}
	if req.Status != "" {
		product.Status = req.Status
	}

	return s.productRepo.Update(product)
}

func (s *ProductService) ReserveProduct(ctx context.Context, productID int64, quantity int) error {
	return s.productRepo.ReserveQuantity(productID, quantity)
}

func (s *ProductService) GetExpiredProducts(ctx context.Context) ([]model.Product, error) {
	return s.productRepo.FindExpired()
}

func (s *ProductService) UpdateExpiredProducts(ctx context.Context, productIDs []int64) error {
	return s.productRepo.UpdateExpiredStatus(productIDs)
}

func (s *ProductService) ProcessExpiredProducts(ctx context.Context) error {
	expiredProducts, err := s.GetExpiredProducts(ctx)
	if err != nil {
		return fmt.Errorf("failed to get expired products: %w", err)
	}

	if len(expiredProducts) == 0 {
		return nil
	}

	ids := make([]int64, len(expiredProducts))
	for i, p := range expiredProducts {
		ids[i] = p.ID
	}

	return s.UpdateExpiredProducts(ctx, ids)
}