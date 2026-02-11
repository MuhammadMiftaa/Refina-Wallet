package router

import (
	"net/http"

	"refina-wallet/config/db"
	"refina-wallet/config/env"
	"refina-wallet/interface/http/middleware"
	"refina-wallet/interface/http/routes"

	"github.com/gin-gonic/gin"
)

func SetupHTTPServer() *http.Server {
	router := gin.Default()

	router.Use(middleware.CORSMiddleware(), middleware.GinMiddleware())

	router.GET("test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World",
		})
	})

	routes.WalletRoutes(router, db.DB)
	routes.WalletTypesRoutes(router, db.DB)

	return &http.Server{
		Addr:    ":" + env.Cfg.Server.HTTPPort,
		Handler: router,
	}
}
