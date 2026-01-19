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

	v1 := router.Group("/v1")

	routes.WalletRoutes(v1, db.DB)
	routes.WalletTypesRoutes(v1, db.DB)

	return router
}
