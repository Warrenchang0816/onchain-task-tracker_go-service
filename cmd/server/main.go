package main

import (
	"log"

	"go-service/internal/config"
	appdb "go-service/internal/db"
	"go-service/internal/handler"
	"go-service/internal/repository"
	appRouter "go-service/internal/router"
	"go-service/internal/service"
)

func main() {
	postgresDB, err := appdb.NewPostgresDB()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer postgresDB.Close()

	log.Println("✅ DB connected")

	if err := appdb.RunMigrations(postgresDB, "./migrations"); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	log.Println("✅ migrations completed")

	blockchainConfig := config.LoadBlockchainConfig()

	taskRepo := repository.NewTaskRepository(postgresDB)
	logRepo := repository.NewBlockchainLogRepository(postgresDB)
	nftOrderRepo := repository.NewNFTOrderRepository(postgresDB)

	taskPermissionSvc := service.NewTaskPermissionService(
		blockchainConfig.GodModeWalletAddress,
	)

	taskRewardVaultSvc, err := service.NewTaskRewardVaultService()
	if err != nil {
		log.Fatalf("failed to init task reward vault service: %v", err)
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
	nftOrderHandler := handler.NewNFTOrderHandler(nftOrderRepo)

	r := appRouter.SetupRouter(taskHandler, logHandler, authHandler, nftOrderHandler)

	log.Println("🚀 server starting on :8080")

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}