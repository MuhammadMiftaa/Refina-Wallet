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
	Wallet_repo := repository.NewWalletRepository(db)
	Wallet_serv := service.NewWalletService(txManager, Wallet_repo, queueInstance)
	Wallet_handler := handler.NewWalletHandler(Wallet_serv)

	wallets := version.Group("/wallets")
	wallets.Use(middleware.AuthMiddleware())

	wallets.GET("", Wallet_handler.GetAllWallets)
	wallets.GET(":id", Wallet_handler.GetWalletByID)
	wallets.GET("user", Wallet_handler.GetWalletsByUserID)
	wallets.GET("user-by-type", Wallet_handler.GetWalletsByUserIDGroupByType)
	wallets.POST("", Wallet_handler.CreateWallet)
	wallets.PUT(":id", Wallet_handler.UpdateWallet)
	wallets.DELETE(":id", Wallet_handler.DeleteWallet)
}
