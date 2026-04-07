package main

import (
	"context"
	"fmt"
	"go-service/internal/config"
	"go-service/internal/db"
	"go-service/internal/handler"
	"go-service/internal/repository"
	"go-service/internal/router"
	"go-service/internal/service"
	"log"
)

// noopVaultService 是 TaskRewardVaultService 的空實作。
// 在沒有設定鏈上 RPC 的本地開發環境中使用，讓伺服器能正常啟動。
type noopVaultService struct{}

func (noopVaultService) AssignWorker(_ context.Context, _ string, _ string) (string, error) {
	return "", fmt.Errorf("blockchain RPC not configured")
}

func (noopVaultService) ApproveTask(_ context.Context, _ string) (string, error) {
	return "", fmt.Errorf("blockchain RPC not configured")
}

func main() {
	postgresDB, err := db.NewPostgresDB()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	log.Printf("Database connected successfully")
	defer postgresDB.Close()

	// ── 任務媒介（原有系統）──────────────────────────────────────
	taskRepo := repository.NewTaskRepository(postgresDB)
	logRepo := repository.NewBlockchainLogRepository(postgresDB)
	blockchainConfig := config.LoadBlockchainConfig()

	taskPermissionSvc := service.NewTaskPermissionService(blockchainConfig.GodModeWalletAddress)

	taskRewardVaultSvc, err := service.NewTaskRewardVaultService()
	if err != nil {
		// 沒有設定鏈上 RPC 時，任務媒介的鏈上功能停用，但其他 API 仍可正常運作
		log.Printf("[WARN] task reward vault service disabled: %v", err)
		taskRewardVaultSvc = noopVaultService{}
	}

	taskService := service.NewTaskService(
		taskRepo,
		logRepo,
		taskPermissionSvc,
		blockchainConfig.PlatformFeeBps,
		taskRewardVaultSvc,
		blockchainConfig,
	)
	taskHandler := handler.NewTaskHandler(taskService, taskPermissionSvc)
	logHandler := handler.NewBlockchainLogHandler(logRepo)
	authHandler := handler.NewAuthHandler(postgresDB)

	// ── 餐廳即期食物販售系統（新增）──────────────────────────────
	productRepo := repository.NewProductRepository(postgresDB)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	orderRepo := repository.NewOrderRepository(postgresDB)
	orderService := service.NewOrderService(orderRepo, productRepo, productService)
	orderHandler := handler.NewOrderHandler(orderService)

	// ── 路由 ─────────────────────────────────────────────────────
	r := router.SetupRouter(
		taskHandler,
		logHandler,
		authHandler,
		productHandler,
		orderHandler,
	)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
