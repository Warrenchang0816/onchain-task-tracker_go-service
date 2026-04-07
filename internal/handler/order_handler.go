package handler

import (
	"go-service/internal/dto"
	"go-service/internal/model"
	"go-service/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req dto.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.orderService.CreateOrder(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, toOrderResponse(*order))
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	orderID := c.Param("id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order id is required"})
		return
	}

	order, err := h.orderService.GetOrder(c.Request.Context(), orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	c.JSON(http.StatusOK, toOrderResponse(*order))
}

func (h *OrderHandler) GetCustomerOrders(c *gin.Context) {
	customerAddress := c.Param("address")
	if customerAddress == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "customer address is required"})
		return
	}

	page := 1
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	limit := 20
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	orders, total, err := h.orderService.GetCustomerOrders(c.Request.Context(), customerAddress, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]dto.OrderResponse, 0, len(orders))
	for _, order := range orders {
		response = append(response, toOrderResponse(order))
	}

	c.JSON(http.StatusOK, dto.OrderListResponse{
		Orders: response,
		Total:  total,
		Page:   page,
		Limit:  limit,
	})
}

func (h *OrderHandler) GetMerchantOrders(c *gin.Context) {
	merchantAddress := c.Param("address")
	if merchantAddress == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "merchant address is required"})
		return
	}

	page := 1
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	limit := 20
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	orders, total, err := h.orderService.GetMerchantOrders(c.Request.Context(), merchantAddress, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]dto.OrderResponse, 0, len(orders))
	for _, order := range orders {
		response = append(response, toOrderResponse(order))
	}

	c.JSON(http.StatusOK, dto.OrderListResponse{
		Orders: response,
		Total:  total,
		Page:   page,
		Limit:  limit,
	})
}

func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	orderID := c.Param("id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order id is required"})
		return
	}

	var req dto.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.orderService.UpdateOrderStatus(c.Request.Context(), orderID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order status updated successfully"})
}

func toOrderResponse(order model.Order) dto.OrderResponse {
	return dto.OrderResponse{
		ID:              order.ID,
		OrderID:         order.OrderID,
		ProductID:       order.ProductID,
		CustomerAddress: order.CustomerAddress,
		MerchantAddress: order.MerchantAddress,
		Quantity:        order.Quantity,
		UnitPrice:       order.UnitPrice,
		TotalPrice:      order.TotalPrice,
		Status:          order.Status,
		PaymentTxHash:   order.PaymentTxHash,
		CreatedAt:       order.CreatedAt,
		CompletedAt:     order.CompletedAt,
	}
}