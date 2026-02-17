package routes

import (
	"refina-wallet/interface/http/handler"
	"refina-wallet/interface/http/middleware"
	"refina-wallet/interface/queue"
	"refina-wallet/internal/repository"
	"refina-wallet/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func WalletRoutes(version *gin.Engine, db *gorm.DB, queueInstance queue.RabbitMQClient) {
	txManager := repository.NewTxManager(db)
	walletRepo := repository.NewWalletRepository(db)
	walletTypeRepo := repository.NewWalletTypesRepository(db)
	outboxRepo := repository.NewOutboxRepository(db)

	walletServ := service.NewWalletService(txManager, walletRepo, walletTypeRepo, outboxRepo, queueInstance)
	walletHandler := handler.NewWalletHandler(walletServ)

	wallets := version.Group("/wallets")
	wallets.Use(middleware.AuthMiddleware())

	wallets.GET("", walletHandler.GetAllWallets)
	wallets.GET(":id", walletHandler.GetWalletByID)
	wallets.GET("user", walletHandler.GetWalletsByUserID)
	wallets.GET("user-by-type", walletHandler.GetWalletsByUserIDGroupByType)
	wallets.POST("", walletHandler.CreateWallet)
	wallets.PUT(":id", walletHandler.UpdateWallet)
	wallets.DELETE(":id", walletHandler.DeleteWallet)
}