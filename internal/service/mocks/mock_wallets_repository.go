package mocks

import (
	"context"

	"refina-wallet/internal/repository"
	"refina-wallet/internal/types/model"
	"refina-wallet/internal/types/view"

	"github.com/stretchr/testify/mock"
)

type MockWalletsRepository struct {
	mock.Mock
}

func (m *MockWalletsRepository) GetAllWallets(ctx context.Context, tx repository.Transaction) ([]model.Wallets, error) {
	args := m.Called(ctx, tx)
	return args.Get(0).([]model.Wallets), args.Error(1)
}

func (m *MockWalletsRepository) GetWalletByID(ctx context.Context, tx repository.Transaction, id string) (model.Wallets, error) {
	args := m.Called(ctx, tx, id)
	return args.Get(0).(model.Wallets), args.Error(1)
}

func (m *MockWalletsRepository) GetWalletsByUserID(ctx context.Context, tx repository.Transaction, id string) ([]model.Wallets, error) {
	args := m.Called(ctx, tx, id)
	return args.Get(0).([]model.Wallets), args.Error(1)
}

func (m *MockWalletsRepository) GetWalletsByUserIDGroupByType(ctx context.Context, tx repository.Transaction, id string) ([]view.ViewUserWalletsGroupByType, error) {
	args := m.Called(ctx, tx, id)
	return args.Get(0).([]view.ViewUserWalletsGroupByType), args.Error(1)
}

func (m *MockWalletsRepository) CreateWallet(ctx context.Context, tx repository.Transaction, wallet model.Wallets) (model.Wallets, error) {
	args := m.Called(ctx, tx, wallet)
	return args.Get(0).(model.Wallets), args.Error(1)
}

func (m *MockWalletsRepository) UpdateWallet(ctx context.Context, tx repository.Transaction, wallet model.Wallets) (model.Wallets, error) {
	args := m.Called(ctx, tx, wallet)
	return args.Get(0).(model.Wallets), args.Error(1)
}

func (m *MockWalletsRepository) DeleteWallet(ctx context.Context, tx repository.Transaction, wallet model.Wallets) (model.Wallets, error) {
	args := m.Called(ctx, tx, wallet)
	return args.Get(0).(model.Wallets), args.Error(1)
}
