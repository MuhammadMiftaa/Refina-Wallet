package mocks

import (
	"context"

	tpb "github.com/MuhammadMiftaa/Refina-Protobuf/transaction"
	"github.com/stretchr/testify/mock"
)

type MockTransactionClient struct {
	mock.Mock
}

func (m *MockTransactionClient) InitialDeposit(ctx context.Context, walletID string, amount float64) (*tpb.TransactionDetail, error) {
	args := m.Called(ctx, walletID, amount)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*tpb.TransactionDetail), args.Error(1)
}

func (m *MockTransactionClient) CancelInitialDeposit(ctx context.Context, transactionID string) (*tpb.TransactionDetail, error) {
	args := m.Called(ctx, transactionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*tpb.TransactionDetail), args.Error(1)
}
