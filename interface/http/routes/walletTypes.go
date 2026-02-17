package routes

import (
	"refina-wallet/interface/http/handler"
	"refina-wallet/interface/http/middleware"
	"refina-wallet/internal/repository"
	"refina-wallet/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func WalletTypesRoutes(version *gin.Engine, db *gorm.DB) {
	txManager := repository.NewTxManager(db)
	WalletTypesRepo := repository.NewWalletTypesRepository(db)
	WalletTypesServ := service.NewWalletTypesService(txManager, WalletTypesRepo)
	WalletTypesHandler := handler.NewWalletTypesHandler(WalletTypesServ)

	version.Use(middleware.AuthMiddleware())
	version.GET("wallet-types", WalletTypesHandler.GetAllWalletTypes)
	version.GET("wallet-types/:id", WalletTypesHandler.GetWalletTypeByID)
	version.POST("wallet-types", WalletTypesHandler.CreateWalletType)
	version.PUT("wallet-types/:id", WalletTypesHandler.UpdateWalletType)
	version.DELETE("wallet-types/:id", WalletTypesHandler.DeleteWalletType)
}
