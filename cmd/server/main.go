package main

import (
	"go-service/internal/config"
	"go-service/internal/db"
	"go-service/internal/handler"
	"go-service/internal/repository"
	"go-service/internal/router"
	"go-service/internal/service"
	"log"
)

func main() {
	cfg := config.Load()

	postgresDB, err := db.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	taskRepo := repository.NewTaskRepository(postgresDB)
	taskService := service.NewTaskService(taskRepo)
	taskHandler := handler.NewTaskHandler(taskService)

	r := router.SetupRouter(taskHandler)

	log.Printf("server is running on :%s", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
