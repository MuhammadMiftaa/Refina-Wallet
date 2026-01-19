package handler

import (
	"net/http"

	"refina-wallet/internal/service"
	"refina-wallet/internal/types/dto"

	"github.com/gin-gonic/gin"
)

type walletHandler struct {
	walletService service.WalletsService
}

func NewWalletHandler(walletService service.WalletsService) *walletHandler {
	return &walletHandler{walletService}
}

func (wallet_handler *walletHandler) GetAllWallets(c *gin.Context) {
	ctx := c.Request.Context()

	wallets, err := wallet_handler.walletService.GetAllWallets(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": 500,
			"status":     false,
			"message":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": 200,
		"status":     true,
		"message":    "Get all wallets",
		"data":       wallets,
	})
}

func (wallet_handler *walletHandler) GetWalletByID(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")

	wallet, err := wallet_handler.walletService.GetWalletByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": 500,
			"status":     false,
			"message":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": 200,
		"status":     true,
		"message":    "Get wallet by ID",
		"data":       wallet,
	})
}

func (wallet_handler *walletHandler) GetWalletsByUserID(c *gin.Context) {
	ctx := c.Request.Context()
	token := c.GetHeader("Authorization")

	userWallets, err := wallet_handler.walletService.GetWalletsByUserID(ctx, token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": 400,
			"status":     false,
			"message":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": 200,
		"status":     true,
		"message":    "Get user wallets data",
		"data":       userWallets,
	})
}

func (wallet_handler *walletHandler) GetWalletsByUserIDGroupByType(c *gin.Context) {
	ctx := c.Request.Context()
	token := c.GetHeader("Authorization")

	userWallets, err := wallet_handler.walletService.GetWalletsByUserIDGroupByType(ctx, token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": 400,
			"status":     false,
			"message":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": 200,
		"status":     true,
		"message":    "Get user wallets data",
		"data":       userWallets,
	})
}

func (wallet_handler *walletHandler) CreateWallet(c *gin.Context) {
	ctx := c.Request.Context()
	token := c.GetHeader("Authorization")

	var walletRequest dto.WalletsRequest
	if err := c.ShouldBindJSON(&walletRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": 400,
			"status":     false,
			"message":    err.Error(),
		})
		return
	}

	wallet, err := wallet_handler.walletService.CreateWallet(ctx, token, walletRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": 500,
			"status":     false,
			"message":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"statusCode": 201,
		"status":     true,
		"message":    "Create wallet",
		"data":       wallet,
	})
}

func (wallet_handler *walletHandler) UpdateWallet(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")

	var walletRequest dto.WalletsRequest
	if err := c.ShouldBindJSON(&walletRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": 400,
			"status":     false,
			"message":    err.Error(),
		})
		return
	}

	wallet, err := wallet_handler.walletService.UpdateWallet(ctx, id, walletRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": 500,
			"status":     false,
			"message":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": 200,
		"status":     true,
		"message":    "Update wallet",
		"data":       wallet,
	})
}

func (wallet_handler *walletHandler) DeleteWallet(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")

	wallet, err := wallet_handler.walletService.DeleteWallet(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": 500,
			"status":     false,
			"message":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": 200,
		"status":     true,
		"message":    "Delete wallet",
		"data":       wallet,
	})
}
