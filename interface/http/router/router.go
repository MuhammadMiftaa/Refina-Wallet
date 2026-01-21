package router

import (
	"refina-wallet/config/db"
	"refina-wallet/interface/http/middleware"
	"refina-wallet/interface/http/routes"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.Use(middleware.CORSMiddleware(), middleware.GinMiddleware())

	router.GET("test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World",
		})
	})

	routes.WalletRoutes(router, db.DB)
	routes.WalletTypesRoutes(router, db.DB)

	return router
}
