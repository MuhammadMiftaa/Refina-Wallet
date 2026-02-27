package router

import (
	"net/http"

	"refina-wallet/config/db"
	"refina-wallet/config/env"
	"refina-wallet/interface/http/middleware"
	"refina-wallet/interface/http/routes"
	"refina-wallet/interface/queue"

	"github.com/gin-gonic/gin"
)

func SetupHTTPServer(dbInstance db.DatabaseClient, queueInstance queue.RabbitMQClient) *http.Server {
	router := gin.New()

	router.Use(
		middleware.CORSMiddleware(),
		middleware.RequestIDMiddleware(),
		middleware.GinMiddleware(),
	)

	router.GET("test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World",
		})
	})

	routes.WalletRoutes(router, dbInstance.GetDB(), queueInstance)
	routes.WalletTypesRoutes(router, dbInstance.GetDB())

	return &http.Server{
		Addr:    ":" + env.Cfg.Server.HTTPPort,
		Handler: router,
	}
}
