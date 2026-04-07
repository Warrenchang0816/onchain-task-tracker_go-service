package handler

import (
	"go-service/internal/dto"
	"go-service/internal/model"
	"go-service/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productService *service.ProductService
}

func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.productService.CreateProduct(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, toProductResponse(*product))
}

func (h *ProductHandler) GetProducts(c *gin.Context) {
	var filters dto.ProductFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	products, total, err := h.productService.ListAvailableProducts(c.Request.Context(), filters, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]dto.ProductResponse, 0, len(products))
	for _, product := range products {
		response = append(response, toProductResponse(product))
	}

	c.JSON(http.StatusOK, dto.ProductListResponse{
		Products: response,
		Total:    total,
		Page:     page,
		Limit:    limit,
	})
}

func (h *ProductHandler) GetProductDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	product, err := h.productService.GetProduct(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	c.JSON(http.StatusOK, toProductResponse(*product))
}

func (h *ProductHandler) GetMerchantProducts(c *gin.Context) {
	merchantAddress := c.Param("address")
	if merchantAddress == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "merchant address is required"})
		return
	}

	filters := dto.ProductFilters{
		Merchant: merchantAddress,
	}

	page := 1
	limit := 50

	products, total, err := h.productService.ListAvailableProducts(c.Request.Context(), filters, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]dto.ProductResponse, 0, len(products))
	for _, product := range products {
		response = append(response, toProductResponse(product))
	}

	c.JSON(http.StatusOK, dto.ProductListResponse{
		Products: response,
		Total:    total,
		Page:     page,
		Limit:    limit,
	})
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.productService.UpdateProduct(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product updated successfully"})
}

func (h *ProductHandler) UpdateStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateReq := dto.UpdateProductRequest{
		Status: req.Status,
	}

	err = h.productService.UpdateProduct(c.Request.Context(), id, updateReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product status updated successfully"})
}

func (h *ProductHandler) UpdateQuantity(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	var req struct {
		Quantity int `json:"quantity" binding:"required,min=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateReq := dto.UpdateProductRequest{
		Quantity: req.Quantity,
	}

	err = h.productService.UpdateProduct(c.Request.Context(), id, updateReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product quantity updated successfully"})
}

func toProductResponse(product model.Product) dto.ProductResponse {
	discountPercentage := float64(product.DiscountRate) / 100.0

	return dto.ProductResponse{
		ID:                  product.ID,
		ProductID:           product.ProductID,
		MerchantAddress:     product.MerchantAddress,
		CustomerAddress:     product.CustomerAddress,
		Name:                product.Name,
		Description:         product.Description,
		Category:            product.Category,
		ImageURL:            product.ImageURL,
		OriginalPrice:       product.OriginalPrice,
		DiscountedPrice:     product.DiscountedPrice,
		DiscountRate:        product.DiscountRate,
		DiscountPercentage:  discountPercentage,
		Quantity:            product.Quantity,
		AvailableQty:        product.AvailableQty,
		Status:              product.Status,
		ExpiryTime:          product.ExpiryTime,
		SalesTarget:         product.SalesTarget,
		SoldQuantity:        product.SoldQuantity,
		PaymentStatus:       product.PaymentStatus,
		CreatedAt:           product.CreatedAt,
		UpdatedAt:           product.UpdatedAt,
		ListedAt:            product.ListedAt,
		SoldOutAt:           product.SoldOutAt,
	}
}