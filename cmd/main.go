package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"convenienceStore/internal/handler"
	"convenienceStore/internal/service"
	"convenienceStore/pkg/config"
	"convenienceStore/pkg/logger"
	"convenienceStore/pkg/payment"
	"convenienceStore/routes"
)

func main() {
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	appLogger := logger.FromConfig(cfg.Logging)

	paymentClient := payment.NewWeChatClient(cfg.Payment, appLogger)

	services := service.NewServices(service.Dependencies{
		Config:  cfg,
		Logger:  appLogger,
		Payment: paymentClient,
	})

	handlers := handler.NewHandlers(services)

	engine := gin.Default()
	routes.RegisterRoutes(engine, routes.HandlerSet{
		User:     handlers.User,
		Product:  handlers.Product,
		Cart:     handlers.Cart,
		Order:    handlers.Order,
		Payment:  handlers.Payment,
		Delivery: handlers.Delivery,
	})

	if err := engine.Run(cfg.Server.Address()); err != nil {
		appLogger.Fatalf("failed to start server: %v", err)
	}
}
