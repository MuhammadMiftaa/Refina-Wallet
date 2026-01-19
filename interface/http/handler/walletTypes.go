package handler

import (
	"net/http"

	"refina-wallet/internal/service"
	"refina-wallet/internal/types/dto"

	"github.com/gin-gonic/gin"
)

type walletTypeHandler struct {
	walletTypeServ service.WalletTypesService
}

func NewWalletTypesHandler(walletTypeServ service.WalletTypesService) *walletTypeHandler {
	return &walletTypeHandler{walletTypeServ}
}

func (walletTypeHandler *walletTypeHandler) GetAllWalletTypes(c *gin.Context) {
	ctx := c.Request.Context()

	walletTypes, err := walletTypeHandler.walletTypeServ.GetAllWalletTypes(ctx)
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
		"message":    "Get all wallet types",
		"data":       walletTypes,
	})
}

func (walletTypeHandler *walletTypeHandler) GetWalletTypeByID(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")

	walletType, err := walletTypeHandler.walletTypeServ.GetWalletTypeByID(ctx, id)
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
		"message":    "Get wallet type by ID",
		"data":       walletType,
	})
}

func (walletTypeHandler *walletTypeHandler) CreateWalletType(c *gin.Context) {
	ctx := c.Request.Context()

	var walletTypeRequest dto.WalletTypesRequest
	if err := c.ShouldBindJSON(&walletTypeRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": 400,
			"status":     false,
			"message":    err.Error(),
		})
		return
	}

	walletType, err := walletTypeHandler.walletTypeServ.CreateWalletType(ctx, walletTypeRequest)
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
		"message":    "Create wallet type",
		"data":       walletType,
	})
}

func (walletTypeHandler *walletTypeHandler) UpdateWalletType(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")

	var walletTypeRequest dto.WalletTypesRequest
	if err := c.ShouldBindJSON(&walletTypeRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": 400,
			"status":     false,
			"message":    err.Error(),
		})
		return
	}

	walletType, err := walletTypeHandler.walletTypeServ.UpdateWalletType(ctx, id, walletTypeRequest)
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
		"message":    "Update wallet type",
		"data":       walletType,
	})
}

func (walletTypeHandler *walletTypeHandler) DeleteWalletType(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")

	walletType, err := walletTypeHandler.walletTypeServ.DeleteWalletType(ctx, id)
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
		"message":    "Delete wallet type",
		"data":       walletType,
	})
}
