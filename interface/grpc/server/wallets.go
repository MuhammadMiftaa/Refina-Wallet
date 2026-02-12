package server

import (
	"context"

	"refina-wallet/internal/repository"
	"refina-wallet/internal/utils"

	wpb "github.com/MuhammadMiftaa/Golang-Refina-Protobuf/wallet"
)

type walletServer struct {
	wpb.UnimplementedWalletServiceServer
	walletsRepository repository.WalletsRepository
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
		CreatedAt:    wallet.CreatedAt.String(),
		UpdatedAt:    wallet.UpdatedAt.String(),
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
		CreatedAt:    updatedWallet.CreatedAt.String(),
		UpdatedAt:    updatedWallet.UpdatedAt.String(),
	}, nil
}
