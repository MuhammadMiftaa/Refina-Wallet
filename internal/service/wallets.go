package service

import (
	"context"
	"encoding/json"
	"errors"

	"refina-wallet/interface/queue"
	"refina-wallet/internal/repository"
	"refina-wallet/internal/types/dto"
	"refina-wallet/internal/types/model"
	"refina-wallet/internal/types/view"
	"refina-wallet/internal/utils"
	"refina-wallet/internal/utils/data"
)

type WalletsService interface {
	GetAllWallets(ctx context.Context) ([]dto.WalletsResponse, error)
	GetWalletByID(ctx context.Context, id string) (dto.WalletsResponse, error)
	GetWalletsByUserID(ctx context.Context, token string) ([]dto.WalletsResponse, error)
	GetWalletsByUserIDGroupByType(ctx context.Context, token string) ([]view.ViewUserWalletsGroupByType, error)
	CreateWallet(ctx context.Context, token string, wallet dto.WalletsRequest) (dto.WalletsResponse, error)
	UpdateWallet(ctx context.Context, id string, wallet dto.WalletsRequest) (dto.WalletsResponse, error)
	DeleteWallet(ctx context.Context, id string) (dto.WalletsResponse, error)
}

type walletsService struct {
	txManager         repository.TxManager
	walletsRepository repository.WalletsRepository
	outboxRepository  repository.OutboxRepository
	queue             queue.RabbitMQClient
}

func NewWalletService(
	txManager repository.TxManager,
	walletsRepository repository.WalletsRepository,
	outboxRepository repository.OutboxRepository,
	queue queue.RabbitMQClient,
) WalletsService {
	return &walletsService{
		txManager:         txManager,
		walletsRepository: walletsRepository,
		outboxRepository:  outboxRepository,
		queue:             queue,
	}
}

func (wallet_serv *walletsService) GetAllWallets(ctx context.Context) ([]dto.WalletsResponse, error) {
	wallets, err := wallet_serv.walletsRepository.GetAllWallets(ctx, nil)
	if err != nil {
		return nil, errors.New("failed to get wallets")
	}

	var walletsResponse []dto.WalletsResponse
	for _, wallet := range wallets {
		walletResponse := utils.ConvertToResponseType(wallet).(dto.WalletsResponse)
		walletsResponse = append(walletsResponse, walletResponse)
	}

	return walletsResponse, nil
}

func (wallet_serv *walletsService) GetWalletByID(ctx context.Context, id string) (dto.WalletsResponse, error) {
	wallet, err := wallet_serv.walletsRepository.GetWalletByID(ctx, nil, id)
	if err != nil {
		return dto.WalletsResponse{}, errors.New("wallet not found")
	}

	walletResponse := utils.ConvertToResponseType(wallet).(dto.WalletsResponse)

	return walletResponse, nil
}

func (wallet_serv *walletsService) GetWalletsByUserID(ctx context.Context, token string) ([]dto.WalletsResponse, error) {
	userData, err := utils.VerifyToken(token[7:])
	if err != nil {
		return nil, errors.New("invalid token")
	}

	wallets, err := wallet_serv.walletsRepository.GetWalletsByUserID(ctx, nil, userData.ID)
	if err != nil {
		return nil, errors.New("failed to get wallets")
	}

	var walletsResponse []dto.WalletsResponse
	for _, wallet := range wallets {
		walletResponse := utils.ConvertToResponseType(wallet).(dto.WalletsResponse)
		walletsResponse = append(walletsResponse, walletResponse)
	}

	return walletsResponse, nil
}

func (wallet_serv *walletsService) GetWalletsByUserIDGroupByType(ctx context.Context, token string) ([]view.ViewUserWalletsGroupByType, error) {
	userData, err := utils.VerifyToken(token[7:])
	if err != nil {
		return nil, errors.New("invalid token")
	}

	wallets, err := wallet_serv.walletsRepository.GetWalletsByUserIDGroupByType(ctx, nil, userData.ID)
	if err != nil {
		return nil, err
	}

	return wallets, err
}

func (wallet_serv *walletsService) CreateWallet(ctx context.Context, token string, wallet dto.WalletsRequest) (dto.WalletsResponse, error) {
	userData, err := utils.VerifyToken(token[7:])
	if err != nil {
		return dto.WalletsResponse{}, errors.New("invalid token")
	}

	UserID, err := utils.ParseUUID(userData.ID)
	if err != nil {
		return dto.WalletsResponse{}, errors.New("invalid user id")
	}

	WalletTypeID, err := utils.ParseUUID(wallet.WalletTypeID)
	if err != nil {
		return dto.WalletsResponse{}, errors.New("invalid wallet type id")
	}

	tx, err := wallet_serv.txManager.Begin(ctx)
	if err != nil {
		return dto.WalletsResponse{}, errors.New("failed to begin transaction")
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// Create wallet
	newWallet, err := wallet_serv.walletsRepository.CreateWallet(ctx, tx, model.Wallets{
		UserID:       UserID,
		WalletTypeID: WalletTypeID,
		Name:         wallet.Name,
		Number:       wallet.Number,
		Balance:      wallet.Balance,
	})
	if err != nil {
		tx.Rollback()
		return dto.WalletsResponse{}, err
	}

	walletResponse := utils.ConvertToResponseType(newWallet).(dto.WalletsResponse)

	payload, err := json.Marshal(walletResponse)
	if err != nil {
		return dto.WalletsResponse{}, errors.New("failed to marshal wallet response")
	}

	outboxMsg := &model.OutboxMessage{
		AggregateID: walletResponse.ID,
		EventType:   data.OUTBOX_EVENT_WALLET_CREATED,
		Payload:     payload,
		Published:   false,
		MaxRetries:  data.OUTBOX_PUBLISH_MAX_RETRIES,
	}

	if err := wallet_serv.outboxRepository.Create(ctx, tx, outboxMsg); err != nil {
		tx.Rollback()
		return dto.WalletsResponse{}, err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return dto.WalletsResponse{}, errors.New("failed to commit transaction")
	}

	return walletResponse, nil
}

func (wallet_serv *walletsService) UpdateWallet(ctx context.Context, id string, wallet dto.WalletsRequest) (dto.WalletsResponse, error) {
	existingWallet, err := wallet_serv.walletsRepository.GetWalletByID(ctx, nil, id)
	if err != nil {
		return dto.WalletsResponse{}, errors.New("wallet not found")
	}

	existingWallet.Name = wallet.Name
	existingWallet.Number = wallet.Number
	existingWallet.Balance = wallet.Balance

	tx, err := wallet_serv.txManager.Begin(ctx)
	if err != nil {
		return dto.WalletsResponse{}, errors.New("failed to begin transaction")
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	walletUpdated, err := wallet_serv.walletsRepository.UpdateWallet(ctx, tx, existingWallet)
	if err != nil {
		return dto.WalletsResponse{}, err
	}

	walletResponse := utils.ConvertToResponseType(walletUpdated).(dto.WalletsResponse)

	payload, err := json.Marshal(walletResponse)
	if err != nil {
		return dto.WalletsResponse{}, errors.New("failed to marshal wallet response")
	}

	outboxMsg := &model.OutboxMessage{
		AggregateID: walletResponse.ID,
		EventType:   data.OUTBOX_EVENT_WALLET_UPDATED,
		Payload:     payload,
		Published:   false,
		MaxRetries:  data.OUTBOX_PUBLISH_MAX_RETRIES,
	}

	if err := wallet_serv.outboxRepository.Create(ctx, tx, outboxMsg); err != nil {
		tx.Rollback()
		return dto.WalletsResponse{}, err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return dto.WalletsResponse{}, errors.New("failed to commit transaction")
	}

	return walletResponse, nil
}

func (wallet_serv *walletsService) DeleteWallet(ctx context.Context, id string) (dto.WalletsResponse, error) {
	existingWallet, err := wallet_serv.walletsRepository.GetWalletByID(ctx, nil, id)
	if err != nil {
		return dto.WalletsResponse{}, errors.New("wallet not found")
	}

	tx, err := wallet_serv.txManager.Begin(ctx)
	if err != nil {
		return dto.WalletsResponse{}, errors.New("failed to begin transaction")
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	deletedWallet, err := wallet_serv.walletsRepository.DeleteWallet(ctx, tx, existingWallet)
	if err != nil {
		return dto.WalletsResponse{}, err
	}

	walletResponse := utils.ConvertToResponseType(deletedWallet).(dto.WalletsResponse)

	payload, err := json.Marshal(walletResponse)
	if err != nil {
		return dto.WalletsResponse{}, errors.New("failed to marshal wallet response")
	}

	outboxMsg := &model.OutboxMessage{
		AggregateID: walletResponse.ID,
		EventType:   data.OUTBOX_EVENT_WALLET_DELETED,
		Payload:     payload,
		Published:   false,
		MaxRetries:  data.OUTBOX_PUBLISH_MAX_RETRIES,
	}

	if err := wallet_serv.outboxRepository.Create(ctx, tx, outboxMsg); err != nil {
		tx.Rollback()
		return dto.WalletsResponse{}, err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return dto.WalletsResponse{}, errors.New("failed to commit transaction")
	}

	return walletResponse, nil
}
