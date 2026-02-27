package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"refina-wallet/config/db"
	"refina-wallet/config/env"
	logger "refina-wallet/config/log"
	"refina-wallet/interface/grpc/client"
	grpcserver "refina-wallet/interface/grpc/server"
	"refina-wallet/interface/http/router"
	"refina-wallet/interface/queue"
	"refina-wallet/internal/repository"
	"refina-wallet/internal/service"
	"refina-wallet/internal/utils"
)

func init() {
	var err error
	var missing []string
	if missing, err = env.LoadByViper(); err != nil {
		log.Printf("Failed to read JSON config file: %v", err)
		if missing, err = env.LoadNative(); err != nil {
			log.Fatalf("Failed to load environment variables: %v", err)
		}
		log.Printf("Environment variables by .env file loaded successfully")
	} else {
		log.Printf("Environment variables by Viper loaded successfully")
	}

	logger.SetupLogger()

	if len(missing) > 0 {
		for _, envVar := range missing {
			logger.Warn("Missing environment variable", map[string]any{"service": "env", "document_id": envVar})
		}
	}
}

func main() {
	// Setup Database Connection
	startTime := time.Now()
	dbInstance := db.GetInstance(env.Cfg.Database)
	logger.Info("Setup Database Connection successfully", map[string]any{"service": "database", "duration": utils.Ms(time.Since(startTime))})

	// Setup RabbitMQ Connection
	startTime = time.Now()
	queueInstance := queue.GetInstance(env.Cfg.RabbitMQ)
	logger.Info("Setup RabbitMQ Connection successfully", map[string]any{"service": "rabbitmq", "duration": utils.Ms(time.Since(startTime))})

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup Outbox Publisher
	startTime = time.Now()
	outboxRepo := repository.NewOutboxRepository(dbInstance.GetDB())
	outboxPublisher := service.NewOutboxPublisher(outboxRepo, queueInstance)

	// Start outbox publisher worker
	go outboxPublisher.Start(ctx)

	// Start cleanup job (optional)
	go outboxPublisher.StartCleanupJob(ctx)
	logger.Info("Outbox Publisher started successfully", map[string]any{"service": "outbox", "duration": utils.Ms(time.Since(startTime))})

	// Set up the gRPC client
	startTime = time.Now()
	grpcManager := client.GetManager()
	err := grpcManager.SetupGRPCClient()
	if err != nil {
		logger.Fatal("Failed to set up gRPC client", map[string]any{"service": "grpc_client", "error": err})
	}
	logger.Info("Setup gRPC client successfully", map[string]any{"service": "grpc_client", "duration": utils.Ms(time.Since(startTime))})

	// Set up the HTTP server
	startTime = time.Now()
	httpServer := router.SetupHTTPServer(dbInstance, queueInstance)
	if httpServer != nil {
		go func() {
			if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logger.Fatal("Failed to start HTTP server", map[string]any{"service": "http_server", "error": err})
			}
		}()
		logger.Info("Starting HTTP server successfully.", map[string]any{"service": "http_server", "port": env.Cfg.Server.HTTPPort, "duration": utils.Ms(time.Since(startTime))})
	}

	// Set up the gRPC server
	startTime = time.Now()
	grpcServer, lis, err := grpcserver.SetupGRPCServer(dbInstance)
	if err != nil {
		logger.Fatal("Failed to set up gRPC server", map[string]any{"service": "grpc_server", "error": err})
	}
	if grpcServer != nil && lis != nil {
		go func() {
			if err := grpcServer.Serve(*lis); err != nil {
				logger.Fatal("Failed to serve gRPC", map[string]any{"service": "grpc_server", "error": err})
			}
		}()
		logger.Info("Starting gRPC server successfully.", map[string]any{"service": "grpc_server", "port": env.Cfg.Server.GRPCPort, "duration": utils.Ms(time.Since(startTime))})
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutdown signal received, stopping services...", map[string]any{"service": "main"})

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	startTime = time.Now()
	shutdownErrors := map[string]any{"service": "main"}

	// Shutdown HTTP server
	if httpServer != nil {
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			logger.Error("Failed to shutdown HTTP server", map[string]any{"service": "http_server", "error": err})
			shutdownErrors["http_error"] = true
		}
	}

	// Cancel context to stop outbox publisher
	cancel()
	time.Sleep(2 * time.Second) // Give some time for outbox publisher to stop

	// Shutdown gRPC server
	if grpcServer != nil {
		grpcServer.GracefulStop()
		if err := grpcManager.Shutdown(shutdownCtx); err != nil {
			logger.Error("Failed to shutdown gRPC clients", map[string]any{"service": "grpc_client", "error": err})
			shutdownErrors["grpc_error"] = true
		}
	}

	// Close RabbitMQ connection
	if err := queueInstance.Close(); err != nil {
		logger.Error("Failed to close RabbitMQ connection", map[string]any{"service": "rabbitmq", "error": err})
		shutdownErrors["rabbitmq_error"] = true
	}

	// Close database connection
	if err := dbInstance.Close(); err != nil {
		logger.Error("Failed to close database connection", map[string]any{"service": "database", "error": err})
		shutdownErrors["database_error"] = true
	}

	if len(shutdownErrors) > 0 {
		logger.Info("Servers stopped with errors", shutdownErrors)
	} else {
		logger.Info("Servers gracefully stopped", map[string]any{"service": "main", "duration": utils.Ms(time.Since(startTime))})
	}
}
