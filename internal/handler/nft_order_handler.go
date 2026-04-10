package handler

import (
	"net/http"
	"strconv"

	"go-service/internal/repository"

	"github.com/gin-gonic/gin"
)

type NFTOrderHandler struct {
	repo *repository.NFTOrderRepository
}

func NewNFTOrderHandler(repo *repository.NFTOrderRepository) *NFTOrderHandler {
	return &NFTOrderHandler{
		repo: repo,
	}
}

type CreateNFTOrderRequest struct {
	Title           string `json:"title"`
	Description     string `json:"description"`
	Image           string `json:"image"`
	Price           string `json:"price"`
	RecipientWallet string `json:"recipientWallet"`
	CreatorWallet   string `json:"creatorWallet"`
}

func (h *NFTOrderHandler) CreateNFTOrder(c *gin.Context) {
	var req CreateNFTOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request body",
		})
		return
	}

	id, err := h.repo.Create(repository.NFTOrder{
		Title:           req.Title,
		Description:     req.Description,
		Image:           req.Image,
		Price:           req.Price,
		RecipientWallet: req.RecipientWallet,
		CreatorWallet:   req.CreatorWallet,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to create nft order",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data": gin.H{
			"id":              id,
			"title":           req.Title,
			"description":     req.Description,
			"image":           req.Image,
			"price":           req.Price,
			"recipientWallet": req.RecipientWallet,
			"creatorWallet":   req.CreatorWallet,
		},
		"message": "NFT order created",
	})
}

func (h *NFTOrderHandler) GetNFTOrders(c *gin.Context) {
	orders, err := h.repo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to get nft orders",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
	"success": true,
	"data":    orders,
	"count":   len(orders),
})
}

type PurchaseNFTOrderRequest struct {
	PurchaseTxHash string `json:"purchaseTxHash"`
}

func (h *NFTOrderHandler) MarkNFTOrderPurchased(c *gin.Context) {
	var req PurchaseNFTOrderRequest

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid order id",
		})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request body",
		})
		return
	}

	if req.PurchaseTxHash == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "purchaseTxHash is required",
		})
		return
	}

	if err := h.repo.UpdatePurchase(id, req.PurchaseTxHash); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to update nft order purchase",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "NFT order marked as purchased",
	})
}


