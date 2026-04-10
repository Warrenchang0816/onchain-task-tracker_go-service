package router

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go-service/internal/handler"
)

func SetupRouter(
	taskHandler *handler.TaskHandler,
	logHandler *handler.BlockchainLogHandler,
	authHandler *handler.AuthHandler,
	nftOrderHandler *handler.NFTOrderHandler,
) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://127.0.0.1:3000",
		},
		AllowMethods: []string{
			"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS",
		},
		AllowHeaders: []string{
			"Origin", "Content-Type", "Accept", "Authorization",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Static("/uploads", "./uploads")

	// Task routes
	r.GET("/api/tasks", taskHandler.GetTasks)
	r.GET("/api/tasks/:id", taskHandler.GetTask)
	r.POST("/api/tasks", taskHandler.CreateTask)
	r.PUT("/api/tasks/:id", taskHandler.UpdateTask)
	r.PATCH("/api/tasks/:id/status", taskHandler.UpdateTaskStatus)
	r.POST("/api/tasks/:id/accept", taskHandler.AcceptTask)
	r.POST("/api/tasks/:id/cancel", taskHandler.CancelTask)
	r.POST("/api/tasks/:id/submit", taskHandler.SubmitTask)
	r.POST("/api/tasks/:id/approve", taskHandler.ApproveTask)
	r.POST("/api/tasks/:id/claim", taskHandler.ClaimReward)
	r.POST("/api/tasks/:id/mark-funded", taskHandler.MarkTaskFunded)
	r.POST("/api/tasks/:id/mark-claimed-onchain", taskHandler.MarkTaskClaimedOnchain)

	// Blockchain log routes
	r.GET("/api/blockchain-logs", logHandler.GetLogs)

	// Auth routes
	r.POST("/api/auth/wallet/siwe/message", authHandler.SIWEMessageHandler)
	r.POST("/api/auth/wallet/siwe/verify", authHandler.SIWEVerifyHandler)
	r.GET("/api/auth/me", authHandler.AuthMeHandler)
	r.POST("/api/auth/logout", authHandler.AuthLogoutHandler)

	// Upload + NFT order routes
	r.POST("/api/upload", handler.UploadFile)
	r.POST("/api/nft-orders", nftOrderHandler.CreateNFTOrder)
	r.POST("/api/nft-orders/:id/purchase", nftOrderHandler.MarkNFTOrderPurchased)
	r.GET("/api/nft-orders", nftOrderHandler.GetNFTOrders)

	return r
}