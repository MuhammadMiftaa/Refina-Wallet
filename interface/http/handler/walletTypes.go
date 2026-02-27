package handler

import (
	"net/http"

	"refina-wallet/config/log"
	"refina-wallet/internal/service"
	"refina-wallet/internal/types/dto"
	"refina-wallet/internal/utils/data"

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
	requestID, _ := c.Get(data.REQUEST_ID_LOCAL_KEY)

	walletTypes, err := walletTypeHandler.walletTypeServ.GetAllWalletTypes(ctx)
	if err != nil {
		log.Error("get_all_wallet_types_failed", map[string]any{
			"service":    data.WalletTypeService,
			"request_id": requestID,
			"error":      err.Error(),
		})
		statusCode, message := mapServiceError(err)
		c.JSON(statusCode, gin.H{
			"statusCode": statusCode,
			"status":     false,
			"message":    message,
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
	requestID, _ := c.Get(data.REQUEST_ID_LOCAL_KEY)

	id := c.Param("id")

	walletType, err := walletTypeHandler.walletTypeServ.GetWalletTypeByID(ctx, id)
	if err != nil {
		log.Error("get_wallet_type_by_id_failed", map[string]any{
			"service":        data.WalletTypeService,
			"request_id":     requestID,
			"wallet_type_id": id,
			"error":          err.Error(),
		})
		statusCode, message := mapServiceError(err)
		c.JSON(statusCode, gin.H{
			"statusCode": statusCode,
			"status":     false,
			"message":    message,
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
	requestID, _ := c.Get(data.REQUEST_ID_LOCAL_KEY)

	var walletTypeRequest dto.WalletTypesRequest
	if err := c.ShouldBindJSON(&walletTypeRequest); err != nil {
		log.Warn("create_wallet_type_bad_request", map[string]any{
			"service":    data.WalletTypeService,
			"request_id": requestID,
			"error":      err.Error(),
		})
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": 400,
			"status":     false,
			"message":    "invalid request body",
		})
		return
	}

	walletType, err := walletTypeHandler.walletTypeServ.CreateWalletType(ctx, walletTypeRequest)
	if err != nil {
		log.Error("create_wallet_type_failed", map[string]any{
			"service":    data.WalletTypeService,
			"request_id": requestID,
			"name":       walletTypeRequest.Name,
			"type":       walletTypeRequest.Type,
			"error":      err.Error(),
		})
		statusCode, message := mapServiceError(err)
		c.JSON(statusCode, gin.H{
			"statusCode": statusCode,
			"status":     false,
			"message":    message,
		})
		return
	}

	log.Info("wallet_type_created", map[string]any{
		"service":        data.WalletTypeService,
		"request_id":     requestID,
		"wallet_type_id": walletType.ID,
		"name":           walletType.Name,
		"type":           walletType.Type,
	})

	c.JSON(http.StatusCreated, gin.H{
		"statusCode": 201,
		"status":     true,
		"message":    "Create wallet type",
		"data":       walletType,
	})
}

func (walletTypeHandler *walletTypeHandler) UpdateWalletType(c *gin.Context) {
	ctx := c.Request.Context()
	requestID, _ := c.Get(data.REQUEST_ID_LOCAL_KEY)

	id := c.Param("id")

	var walletTypeRequest dto.WalletTypesRequest
	if err := c.ShouldBindJSON(&walletTypeRequest); err != nil {
		log.Warn("update_wallet_type_bad_request", map[string]any{
			"service":        data.WalletTypeService,
			"request_id":     requestID,
			"wallet_type_id": id,
			"error":          err.Error(),
		})
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": 400,
			"status":     false,
			"message":    "invalid request body",
		})
		return
	}

	walletType, err := walletTypeHandler.walletTypeServ.UpdateWalletType(ctx, id, walletTypeRequest)
	if err != nil {
		log.Error("update_wallet_type_failed", map[string]any{
			"service":        data.WalletTypeService,
			"request_id":     requestID,
			"wallet_type_id": id,
			"error":          err.Error(),
		})
		statusCode, message := mapServiceError(err)
		c.JSON(statusCode, gin.H{
			"statusCode": statusCode,
			"status":     false,
			"message":    message,
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
	requestID, _ := c.Get(data.REQUEST_ID_LOCAL_KEY)

	id := c.Param("id")

	walletType, err := walletTypeHandler.walletTypeServ.DeleteWalletType(ctx, id)
	if err != nil {
		log.Error("delete_wallet_type_failed", map[string]any{
			"service":        data.WalletTypeService,
			"request_id":     requestID,
			"wallet_type_id": id,
			"error":          err.Error(),
		})
		statusCode, message := mapServiceError(err)
		c.JSON(statusCode, gin.H{
			"statusCode": statusCode,
			"status":     false,
			"message":    message,
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
