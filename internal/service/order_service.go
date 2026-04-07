package service

import (
	"context"
	"fmt"
	"go-service/internal/dto"
	"go-service/internal/model"
	"go-service/internal/repository"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type OrderService struct {
	orderRepo    *repository.OrderRepository
	productRepo  *repository.ProductRepository
	productService *ProductService
}

func NewOrderService(orderRepo *repository.OrderRepository, productRepo *repository.ProductRepository, productService *ProductService) *OrderService {
	return &OrderService{
		orderRepo:      orderRepo,
		productRepo:    productRepo,
		productService: productService,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, req dto.CreateOrderRequest) (*model.Order, error) {
	// 获取商品信息
	product, err := s.productRepo.FindByID(req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("failed to find product: %w", err)
	}

	if product.Status != "AVAILABLE" {
		return nil, fmt.Errorf("product is not available")
	}

	if product.AvailableQty < req.Quantity {
		return nil, fmt.Errorf("insufficient stock")
	}

	// 检查过期
	if product.ExpiryTime != nil && product.ExpiryTime.Before(time.Now()) {
		return nil, fmt.Errorf("product has expired")
	}

	// 预留库存
	err = s.productService.ReserveProduct(ctx, req.ProductID, req.Quantity)
	if err != nil {
		return nil, fmt.Errorf("failed to reserve product: %w", err)
	}

	// 计算总价
	totalPrice := product.DiscountedPrice.Mul(decimal.NewFromInt(int64(req.Quantity)))

	order := &model.Order{
		OrderID:         uuid.New().String(),
		ProductID:       req.ProductID,
		CustomerAddress: req.CustomerAddress,
		MerchantAddress: product.MerchantAddress,
		Quantity:        req.Quantity,
		UnitPrice:       product.DiscountedPrice,
		TotalPrice:      totalPrice,
		Status:          "PLACED",
	}

	id, err := s.orderRepo.Create(order)
	if err != nil {
		// 如果创建订单失败，需要释放预留的库存
		// 这里应该有补偿逻辑，但暂时简化
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	order.ID = id
	return order, nil
}

func (s *OrderService) GetOrder(ctx context.Context, orderID string) (*model.Order, error) {
	return s.orderRepo.FindByOrderID(orderID)
}

func (s *OrderService) GetCustomerOrders(ctx context.Context, customerAddress string, page, limit int) ([]model.Order, int64, error) {
	offset := (page - 1) * limit
	return s.orderRepo.FindByCustomer(customerAddress, limit, offset)
}

func (s *OrderService) GetMerchantOrders(ctx context.Context, merchantAddress string, page, limit int) ([]model.Order, int64, error) {
	offset := (page - 1) * limit
	return s.orderRepo.FindByMerchant(merchantAddress, limit, offset)
}

func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderID string, req dto.UpdateOrderStatusRequest) error {
	var completedAt *string
	if req.Status == "COMPLETED" {
		now := time.Now().Format(time.RFC3339)
		completedAt = &now
	}

	return s.orderRepo.UpdateStatus(orderID, req.Status, req.PaymentTxHash, completedAt)
}

func (s *OrderService) ProcessExpiredOrders(ctx context.Context) error {
	// 处理超过24小时未确认的订单
	expiredOrders, err := s.orderRepo.FindByStatusAndAge("PLACED", 24)
	if err != nil {
		return fmt.Errorf("failed to get expired orders: %w", err)
	}

	for _, order := range expiredOrders {
		// 取消订单并释放库存
		err = s.UpdateOrderStatus(ctx, order.OrderID, dto.UpdateOrderStatusRequest{
			Status: "CANCELLED",
		})
		if err != nil {
			return fmt.Errorf("failed to cancel expired order %s: %w", order.OrderID, err)
		}

		// 这里应该有释放库存的逻辑，但暂时简化
		// 实际应该调用productService来释放库存
	}

	return nil
}