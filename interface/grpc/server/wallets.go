package server

import (
	"context"
	"fmt"
	"time"

	"refina-wallet/config/log"
	"refina-wallet/internal/repository"
	"refina-wallet/internal/utils"
	"refina-wallet/internal/utils/data"

	wpb "github.com/MuhammadMiftaa/Refina-Protobuf/wallet"
)

type walletServer struct {
	wpb.UnimplementedWalletServiceServer
	walletsRepository repository.WalletsRepository
}

func (s *walletServer) GetWallets(req *wpb.GetWalletOptions, stream wpb.WalletService_GetUserWalletsServer) error {
	timeout, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	wallets, err := s.walletsRepository.GetAllWallets(timeout, nil)
	if err != nil {
		log.Error(data.LogGetAllWalletsFailed, map[string]any{
			"service": data.GRPCServerService,
			"error":   err.Error(),
		})
		return fmt.Errorf("get all wallets: %w", err)
	}

	for _, wallet := range wallets {
		if err := stream.Send(&wpb.Wallet{
			Id:             wallet.ID.String(),
			UserId:         wallet.UserID.String(),
			Name:           wallet.Name,
			Number:         wallet.Number,
			Balance:        wallet.Balance,
			WalletTypeId:   wallet.WalletTypeID.String(),
			WalletType:     string(wallet.WalletType.Type),
			WalletTypeName: wallet.WalletType.Name,
			CreatedAt:      wallet.CreatedAt.Format(time.RFC3339),
			UpdatedAt:      wallet.UpdatedAt.Format(time.RFC3339),
		}); err != nil {
			log.Error(data.LogGetAllWalletsStreamFailed, map[string]any{
				"service":   data.GRPCServerService,
				"wallet_id": wallet.ID.String(),
				"error":     err.Error(),
			})
			return fmt.Errorf("stream send wallet [id=%s]: %w", wallet.ID.String(), err)
		}
	}

	return nil
}

func (s *walletServer) GetUserWallets(req *wpb.UserID, stream wpb.WalletService_GetUserWalletsServer) error {
	timeout, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	userID := req.GetId()

	wallets, err := s.walletsRepository.GetWalletsByUserID(timeout, nil, userID)
	if err != nil {
		log.Error(data.LogGetUserWalletsFailed, map[string]any{
			"service": data.GRPCServerService,
			"user_id": userID,
			"error":   err.Error(),
		})
		return fmt.Errorf("get wallets by user [id=%s]: %w", userID, err)
	}

	for _, wallet := range wallets {
		if err := stream.Send(&wpb.Wallet{
			Id:             wallet.ID.String(),
			UserId:         wallet.UserID.String(),
			Name:           wallet.Name,
			Number:         wallet.Number,
			Balance:        wallet.Balance,
			WalletTypeId:   wallet.WalletTypeID.String(),
			WalletType:     string(wallet.WalletType.Type),
			WalletTypeName: wallet.WalletType.Name,
			CreatedAt:      wallet.CreatedAt.Format(time.RFC3339),
			UpdatedAt:      wallet.UpdatedAt.Format(time.RFC3339),
		}); err != nil {
			log.Error(data.LogGetUserWalletsStreamFailed, map[string]any{
				"service":   data.GRPCServerService,
				"user_id":   userID,
				"wallet_id": wallet.ID.String(),
				"error":     err.Error(),
			})
			return fmt.Errorf("stream send wallet [id=%s]: %w", wallet.ID.String(), err)
		}
	}

	return nil
}

func (s *walletServer) GetWalletByID(ctx context.Context, req *wpb.WalletID) (*wpb.Wallet, error) {
	walletID := req.GetId()

	wallet, err := s.walletsRepository.GetWalletByID(ctx, nil, walletID)
	if err != nil {
		log.Error(data.LogGetWalletByIDFailed, map[string]any{
			"service":   data.GRPCServerService,
			"wallet_id": walletID,
			"error":     err.Error(),
		})
		return nil, fmt.Errorf("get wallet [id=%s]: %w", walletID, err)
	}

	return &wpb.Wallet{
		Id:           wallet.ID.String(),
		UserId:       wallet.UserID.String(),
		Name:         wallet.Name,
		Number:       wallet.Number,
		Balance:      wallet.Balance,
		WalletTypeId: wallet.WalletTypeID.String(),
		CreatedAt:    wallet.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    wallet.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *walletServer) UpdateWallet(ctx context.Context, req *wpb.Wallet) (*wpb.Wallet, error) {
	walletID := req.GetId()

	wallet, err := s.walletsRepository.GetWalletByID(ctx, nil, walletID)
	if err != nil {
		log.Error(data.LogUpdateWalletNotFound, map[string]any{
			"service":   data.GRPCServerService,
			"wallet_id": walletID,
			"error":     err.Error(),
		})
		return nil, fmt.Errorf("wallet not found [id=%s]: %w", walletID, err)
	}

	wallet.Name = req.GetName()
	wallet.Number = req.GetNumber()
	wallet.Balance = req.GetBalance()

	walletTypeID, err := utils.ParseUUID(req.GetWalletTypeId())
	if err != nil {
		log.Error(data.LogUpdateWalletInvalidTypeID, map[string]any{
			"service":        data.GRPCServerService,
			"wallet_id":      walletID,
			"wallet_type_id": req.GetWalletTypeId(),
			"error":          err.Error(),
		})
		return nil, fmt.Errorf("invalid wallet type id [id=%s]: %w", req.GetWalletTypeId(), err)
	}
	wallet.WalletTypeID = walletTypeID

	updatedWallet, err := s.walletsRepository.UpdateWallet(ctx, nil, wallet)
	if err != nil {
		log.Error(data.LogUpdateWalletFailed, map[string]any{
			"service":   data.GRPCServerService,
			"wallet_id": walletID,
			"error":     err.Error(),
		})
		return nil, fmt.Errorf("update wallet [id=%s]: %w", walletID, err)
	}

	log.Info(data.LogWalletUpdated, map[string]any{
		"service":   data.GRPCServerService,
		"wallet_id": updatedWallet.ID.String(),
	})

	return &wpb.Wallet{
		Id:           updatedWallet.ID.String(),
		UserId:       updatedWallet.UserID.String(),
		Name:         updatedWallet.Name,
		Number:       updatedWallet.Number,
		Balance:      updatedWallet.Balance,
		WalletTypeId: updatedWallet.WalletTypeID.String(),
		CreatedAt:    updatedWallet.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    updatedWallet.UpdatedAt.Format(time.RFC3339),
	}, nil
}
