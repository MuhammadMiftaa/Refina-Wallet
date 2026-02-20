package client

import (
	"context"
	"time"

	"refina-wallet/internal/utils/data"

	tpb "github.com/MuhammadMiftaa/Refina-Protobuf/transaction"
)

type TransactionClient interface {
	InitialDeposit(ctx context.Context, walletID string, amount float64) (*tpb.Transaction, error)
	CancelInitialDeposit(ctx context.Context, transactionID string) (*tpb.Transaction, error)
}

type transactionClientImpl struct {
	client tpb.TransactionServiceClient
}

func NewTransactionClient(grpcClient tpb.TransactionServiceClient) TransactionClient {
	return &transactionClientImpl{
		client: grpcClient,
	}
}

func (t *transactionClientImpl) InitialDeposit(ctx context.Context, walletID string, amount float64) (*tpb.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return t.client.CreateTransaction(ctx, &tpb.NewTransaction{
		WalletId:        walletID,
		Amount:          amount,
		CategoryId:      data.INITIAL_DEPOSIT_CATEGORY_ID,
		TransactionDate: time.Now().Format(time.RFC3339),
		Description:     data.INITIAL_DEPOSIT_DESC,
	})
}

func (t *transactionClientImpl) CancelInitialDeposit(ctx context.Context, transactionID string) (*tpb.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return t.client.DeleteTransaction(ctx, &tpb.TransactionID{Id: transactionID})
}
