package mocks

import (
	"context"

	"refina-wallet/internal/repository"
	"refina-wallet/internal/types/model"

	"github.com/stretchr/testify/mock"
)

type MockWalletTypesRepository struct {
	mock.Mock
}

func (m *MockWalletTypesRepository) GetAllWalletTypes(ctx context.Context, tx repository.Transaction) ([]model.WalletTypes, error) {
	args := m.Called(ctx, tx)
	return args.Get(0).([]model.WalletTypes), args.Error(1)
}

func (m *MockWalletTypesRepository) GetWalletTypeByID(ctx context.Context, tx repository.Transaction, id string) (model.WalletTypes, error) {
	args := m.Called(ctx, tx, id)
	return args.Get(0).(model.WalletTypes), args.Error(1)
}

func (m *MockWalletTypesRepository) CreateWalletType(ctx context.Context, tx repository.Transaction, walletType model.WalletTypes) (model.WalletTypes, error) {
	args := m.Called(ctx, tx, walletType)
	return args.Get(0).(model.WalletTypes), args.Error(1)
}

func (m *MockWalletTypesRepository) UpdateWalletType(ctx context.Context, tx repository.Transaction, walletType model.WalletTypes) (model.WalletTypes, error) {
	args := m.Called(ctx, tx, walletType)
	return args.Get(0).(model.WalletTypes), args.Error(1)
}

func (m *MockWalletTypesRepository) DeleteWalletType(ctx context.Context, tx repository.Transaction, walletType model.WalletTypes) (model.WalletTypes, error) {
	args := m.Called(ctx, tx, walletType)
	return args.Get(0).(model.WalletTypes), args.Error(1)
}
