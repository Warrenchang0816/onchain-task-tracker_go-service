package router

import (
	"time"

	"go-service/internal/auth"
	"go-service/internal/handler"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	taskHandler *handler.TaskHandler,
	logHandler *handler.BlockchainLogHandler,
	authHandler *handler.AuthHandler,
	productHandler *handler.ProductHandler,
	orderHandler *handler.OrderHandler,
) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := r.Group("/api")
	{
		// 健康檢查
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		// ── 身份驗證 ─────────────────────────────────────────────
		api.POST("/auth/wallet/siwe/message", authHandler.SIWEMessageHandler)
		api.POST("/auth/wallet/siwe/verify", authHandler.SIWEVerifyHandler)
		api.GET("/auth/me", authHandler.AuthMeHandler)
		api.POST("/auth/logout", authHandler.AuthLogoutHandler)

		// ── 任務媒介（原有系統）──────────────────────────────────
		publicTask := api.Group("")
		publicTask.Use(auth.OptionalAuthMiddleware(authHandler.GetSessionRepository()))
		{
			publicTask.GET("/tasks", taskHandler.GetTasks)
			publicTask.GET("/tasks/:id", taskHandler.GetTask)
		}

		protected := api.Group("")
		protected.Use(auth.AuthMiddleware(authHandler.GetSessionRepository()))
		{
			protected.POST("/tasks", taskHandler.CreateTask)
			protected.PUT("/tasks/:id", taskHandler.UpdateTask)
			protected.PUT("/tasks/:id/status", taskHandler.UpdateTaskStatus)

			protected.PUT("/tasks/:id/accept", taskHandler.AcceptTask)
			protected.PUT("/tasks/:id/cancel", taskHandler.CancelTask)
			protected.POST("/tasks/:id/submissions", taskHandler.SubmitTask)

			protected.PUT("/tasks/:id/approve", taskHandler.ApproveTask)
			protected.POST("/tasks/:id/claim", taskHandler.ClaimReward)

			protected.POST("/tasks/:id/onchain/funded", taskHandler.MarkTaskFunded)
			protected.POST("/tasks/:id/onchain/claimed", taskHandler.MarkTaskClaimedOnchain)

			protected.PUT("/tasks/:id/fund", taskHandler.FundTask)
		}

		// 鏈上事件記錄
		api.GET("/blockchain-logs", logHandler.GetLogs)

		// ── 餐廳即期食物販售系統（新增）──────────────────────────
		// 商品管理 API
		api.POST("/products", productHandler.CreateProduct)
		api.GET("/products", productHandler.GetProducts)
		api.GET("/products/:id", productHandler.GetProductDetail)
		api.GET("/products/merchant/:address", productHandler.GetMerchantProducts)
		api.PUT("/products/:id", productHandler.UpdateProduct)
		api.PUT("/products/:id/status", productHandler.UpdateStatus)
		api.PUT("/products/:id/quantity", productHandler.UpdateQuantity)

		// 訂單管理 API
		api.POST("/orders", orderHandler.CreateOrder)
		api.GET("/orders/:id", orderHandler.GetOrder)
		api.GET("/orders/customer/:address", orderHandler.GetCustomerOrders)
		api.GET("/orders/merchant/:address", orderHandler.GetMerchantOrders)
		api.PUT("/orders/:id/status", orderHandler.UpdateOrderStatus)
	}

	return r
}
