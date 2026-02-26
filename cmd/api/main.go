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
			logger.Log.Warnf("Missing environment variable: %s", envVar)
		}
	}
}

func main() {
	// Setup Database Connection
	startTime := time.Now()
	dbInstance := db.GetInstance(env.Cfg.Database)
	logger.Info("Setup Database Connection successfully", map[string]any{"duration": utils.Ms(time.Since(startTime))})

	// Setup RabbitMQ Connection
	startTime = time.Now()
	queueInstance := queue.GetInstance(env.Cfg.RabbitMQ)
	logger.Info("Setup RabbitMQ Connection successfully", map[string]any{"duration": utils.Ms(time.Since(startTime))})

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
	logger.Info("Outbox Publisher started successfully", map[string]any{"duration": utils.Ms(time.Since(startTime))})

	// Set up the gRPC client
	startTime = time.Now()
	grpcManager := client.GetManager()
	err := grpcManager.SetupGRPCClient()
	if err != nil {
		logger.Log.Fatalf("Failed to set up gRPC client: %v", err)
	}
	logger.Info("Setup gRPC client successfully", map[string]any{"duration": utils.Ms(time.Since(startTime))})

	// Set up the HTTP server
	startTime = time.Now()
	httpServer := router.SetupHTTPServer(dbInstance, queueInstance)
	if httpServer != nil {
		go func() {
			if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logger.Log.Fatalf("Failed to start HTTP server: %s\n", err)
			}
		}()
		logger.Info("Starting HTTP server successfully.", map[string]any{"port": env.Cfg.Server.HTTPPort, "duration": utils.Ms(time.Since(startTime))})
	}

	// Set up the gRPC server
	startTime = time.Now()
	grpcServer, lis, err := grpcserver.SetupGRPCServer(dbInstance)
	if err != nil {
		logger.Log.Fatalf("Failed to set up gRPC server: %v", err)
	}
	if grpcServer != nil && lis != nil {
		go func() {
			if err := grpcServer.Serve(*lis); err != nil {
				logger.Log.Fatalf("Failed to serve gRPC: %v", err)
			}
		}()
		logger.Info("Starting gRPC server successfully.", map[string]any{"port": env.Cfg.Server.GRPCPort, "duration": utils.Ms(time.Since(startTime))})
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutdown signal received, stopping services...", nil)

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	startTime = time.Now()
	shutdownErrors := make(map[string]any)

	// Shutdown HTTP server
	if httpServer != nil {
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			logger.Log.Errorf("Failed to shutdown HTTP server: %v", err)
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
			logger.Log.Errorf("Failed to shutdown gRPC clients: %v", err)
			shutdownErrors["grpc_error"] = true
		}
	}

	// Close RabbitMQ connection
	if err := queueInstance.Close(); err != nil {
		logger.Log.Errorf("Failed to close RabbitMQ connection: %v", err)
		shutdownErrors["rabbitmq_error"] = true
	}

	// Close database connection
	if err := dbInstance.Close(); err != nil {
		logger.Log.Errorf("Failed to close database connection: %v", err)
		shutdownErrors["database_error"] = true
	}

	if len(shutdownErrors) > 0 {
		logger.Info("Servers stopped with errors", shutdownErrors)
	} else {
		logger.Info("Servers gracefully stopped", map[string]any{"duration": utils.Ms(time.Since(startTime))})
	}
}
