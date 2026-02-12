package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"refina-wallet/config/db"
	"refina-wallet/config/env"
	"refina-wallet/config/log"
	grpcserver "refina-wallet/interface/grpc/server"
	"refina-wallet/interface/http/router"
)

func init() {
	log.SetupLogger() // Initialize the logger configuration

	var err error
	var missing []string
	if missing, err = env.LoadByViper(); err != nil {
		log.Error("Failed to read JSON config file:" + err.Error())
		if missing, err = env.LoadNative(); err != nil {
			log.Log.Fatalf("Failed to load environment variables: %v", err)
		}
		log.SetupLogger()
		log.Info("Environment variables by .env file loaded successfully")
	} else {
		log.SetupLogger()
		log.Info("Environment variables by Viper loaded successfully")
	}

	if len(missing) > 0 {
		for _, envVar := range missing {
			log.Warn("Missing environment variable: " + envVar)
		}
	}

	log.Info("Setup Database Connection Start")
	db.SetupDatabase(env.Cfg.Database) // Initialize the database connection and run migrations
	log.Info("Setup Database Connection Success")

	log.Info("Starting Refina API...")
}

func main() {
	defer log.Info("Refina API stopped")

	// Set up the HTTP server
	httpServer := router.SetupHTTPServer()
	if httpServer == nil {
		go func() {
			if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Log.Fatalf("Failed to start HTTP server: %s\n", err)
			}
		}()
	}

	// Set up the gRPC server
	grpcServer, lis, err := grpcserver.SetupGRPCServer()
	if err != nil {
		log.Log.Fatalf("Failed to set up gRPC server: %v", err)
	}
	if grpcServer != nil && lis != nil {
		go func() {
			if err := grpcServer.Serve(*lis); err != nil {
				log.Log.Fatalf("Failed to serve gRPC: %v", err)
			}
		}()
		log.Info("Starting gRPC server successfully")
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Log.Info("Shutting down servers...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Log.Fatalf("Failed to shutdown HTTP server: %v", err)
	}

	grpcServer.GracefulStop()

	log.Log.Info("Servers gracefully stopped")
}
