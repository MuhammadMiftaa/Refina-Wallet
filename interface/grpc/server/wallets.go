package server

import (
	"context"
	"time"

	"refina-wallet/internal/repository"
	"refina-wallet/internal/utils"

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
		return err
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
			return err
		}
	}

	return nil
}

func (s *walletServer) GetUserWallets(req *wpb.UserID, stream wpb.WalletService_GetUserWalletsServer) error {
	timeout, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	wallets, err := s.walletsRepository.GetWalletsByUserID(timeout, nil, req.GetId())
	if err != nil {
		return err
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
			return err
		}
	}

	return nil
}

func (s *walletServer) GetWalletByID(ctx context.Context, req *wpb.WalletID) (*wpb.Wallet, error) {
	wallet, err := s.walletsRepository.GetWalletByID(ctx, nil, req.GetId())
	if err != nil {
		return nil, err
	}

	walletResponse := &wpb.Wallet{
		Id:           wallet.ID.String(),
		UserId:       wallet.UserID.String(),
		Name:         wallet.Name,
		Number:       wallet.Number,
		Balance:      wallet.Balance,
		WalletTypeId: wallet.WalletTypeID.String(),
		CreatedAt:    wallet.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    wallet.UpdatedAt.Format(time.RFC3339),
	}

	return walletResponse, nil
}

func (s *walletServer) UpdateWallet(ctx context.Context, req *wpb.Wallet) (*wpb.Wallet, error) {
	wallet, err := s.walletsRepository.GetWalletByID(ctx, nil, req.GetId())
	if err != nil {
		return nil, err
	}

	wallet.Name = req.GetName()
	wallet.Number = req.GetNumber()
	wallet.Balance = req.GetBalance()
	walletTypeID, err := utils.ParseUUID(req.GetWalletTypeId())
	if err != nil {
		return nil, err
	}
	wallet.WalletTypeID = walletTypeID

	updatedWallet, err := s.walletsRepository.UpdateWallet(ctx, nil, wallet)
	if err != nil {
		return nil, err
	}

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
