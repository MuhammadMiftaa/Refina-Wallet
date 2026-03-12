package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"refina-wallet/internal/service/mocks"
	"refina-wallet/internal/types/dto"
	"refina-wallet/internal/types/model"
	"refina-wallet/internal/types/view"

	tpb "github.com/MuhammadMiftaa/Refina-Protobuf/transaction"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ---------- helpers ----------

type walletTestDeps struct {
	txManager   *mocks.MockTxManager
	walletsRepo *mocks.MockWalletsRepository
	typesRepo   *mocks.MockWalletTypesRepository
	outboxRepo  *mocks.MockOutboxRepository
	txClient    *mocks.MockTransactionClient
	rabbitMQ    *mocks.MockRabbitMQClient
	tx          *mocks.MockTransaction
}

func newWalletTestDeps() *walletTestDeps {
	return &walletTestDeps{
		txManager:   new(mocks.MockTxManager),
		walletsRepo: new(mocks.MockWalletsRepository),
		typesRepo:   new(mocks.MockWalletTypesRepository),
		outboxRepo:  new(mocks.MockOutboxRepository),
		txClient:    new(mocks.MockTransactionClient),
		rabbitMQ:    new(mocks.MockRabbitMQClient),
		tx:          new(mocks.MockTransaction),
	}
}

func (d *walletTestDeps) service() WalletsService {
	return NewWalletService(
		d.txManager,
		d.walletsRepo,
		d.typesRepo,
		d.outboxRepo,
		d.txClient,
		d.rabbitMQ,
	)
}

func (d *walletTestDeps) assertAll(t *testing.T) {
	t.Helper()
	d.txManager.AssertExpectations(t)
	d.walletsRepo.AssertExpectations(t)
	d.typesRepo.AssertExpectations(t)
	d.outboxRepo.AssertExpectations(t)
	d.txClient.AssertExpectations(t)
	d.rabbitMQ.AssertExpectations(t)
	d.tx.AssertExpectations(t)
}

var (
	walletID     = uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	userID       = uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")
	walletTypeID = uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc")
	fixedTime    = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
)

func sampleWalletType() model.WalletTypes {
	return model.WalletTypes{
		Base: model.Base{ID: walletTypeID, CreatedAt: fixedTime, UpdatedAt: fixedTime},
		Name: "BCA", Type: model.Bank, Description: "Bank BCA",
	}
}

func sampleWalletModel() model.Wallets {
	return model.Wallets{
		Base:         model.Base{ID: walletID, CreatedAt: fixedTime, UpdatedAt: fixedTime},
		UserID:       userID,
		WalletTypeID: walletTypeID,
		Name:         "My BCA",
		Number:       "1234567890",
		Balance:      100000,
		WalletType:   sampleWalletType(),
	}
}

func sampleWalletRequest() dto.WalletsRequest {
	return dto.WalletsRequest{
		UserID:       userID.String(),
		WalletTypeID: walletTypeID.String(),
		Name:         "My BCA",
		Number:       "1234567890",
		Balance:      100000,
	}
}

// =====================================================================
// GetAllWallets
// =====================================================================

func TestGetAllWallets_Success(t *testing.T) {
	d := newWalletTestDeps()
	svc := d.service()

	wallets := []model.Wallets{sampleWalletModel()}
	d.walletsRepo.On("GetAllWallets", mock.Anything, nil).Return(wallets, nil)

	result, err := svc.GetAllWallets(context.Background())

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, walletID.String(), result[0].ID)
	assert.Equal(t, "My BCA", result[0].Name)
	d.assertAll(t)
}

func TestGetAllWallets_EmptyList(t *testing.T) {
	d := newWalletTestDeps()
	svc := d.service()

	d.walletsRepo.On("GetAllWallets", mock.Anything, nil).Return([]model.Wallets{}, nil)

	result, err := svc.GetAllWallets(context.Background())

	assert.NoError(t, err)
	assert.Empty(t, result)
	d.assertAll(t)
}

func TestGetAllWallets_RepositoryError(t *testing.T) {
	d := newWalletTestDeps()
	svc := d.service()

	d.walletsRepo.On("GetAllWallets", mock.Anything, nil).
		Return([]model.Wallets{}, errors.New("db error"))

	result, err := svc.GetAllWallets(context.Background())

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "get all wallets")
	d.assertAll(t)
}

// =====================================================================
// GetWalletByID
// =====================================================================

func TestGetWalletByID_Success(t *testing.T) {
	d := newWalletTestDeps()
	svc := d.service()

	w := sampleWalletModel()
	id := w.ID.String()
	d.walletsRepo.On("GetWalletByID", mock.Anything, nil, id).Return(w, nil)

	result, err := svc.GetWalletByID(context.Background(), id)

	assert.NoError(t, err)
	assert.Equal(t, id, result.ID)
	assert.Equal(t, w.Name, result.Name)
	assert.Equal(t, w.Balance, result.Balance)
	d.assertAll(t)
}

func TestGetWalletByID_NotFound(t *testing.T) {
	d := newWalletTestDeps()
	svc := d.service()

	id := uuid.New().String()
	d.walletsRepo.On("GetWalletByID", mock.Anything, nil, id).
		Return(model.Wallets{}, errors.New("record not found"))

	result, err := svc.GetWalletByID(context.Background(), id)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "get wallet")
	assert.Empty(t, result.ID)
	d.assertAll(t)
}

// =====================================================================
// GetWalletsByUserID
// =====================================================================

func TestGetWalletsByUserID_Success(t *testing.T) {
	d := newWalletTestDeps()
	svc := d.service()

	uid := userID.String()
	wallets := []model.Wallets{sampleWalletModel()}
	d.walletsRepo.On("GetWalletsByUserID", mock.Anything, nil, uid).Return(wallets, nil)

	result, err := svc.GetWalletsByUserID(context.Background(), uid)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, uid, result[0].UserID)
	d.assertAll(t)
}

func TestGetWalletsByUserID_RepositoryError(t *testing.T) {
	d := newWalletTestDeps()
	svc := d.service()

	uid := userID.String()
	d.walletsRepo.On("GetWalletsByUserID", mock.Anything, nil, uid).
		Return([]model.Wallets{}, errors.New("db error"))

	result, err := svc.GetWalletsByUserID(context.Background(), uid)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "get wallets by user")
	d.assertAll(t)
}

// =====================================================================
// GetWalletsByUserIDGroupByType
// =====================================================================

func TestGetWalletsByUserIDGroupByType_Success(t *testing.T) {
	d := newWalletTestDeps()
	svc := d.service()

	uid := userID.String()
	grouped := []view.ViewUserWalletsGroupByType{
		{
			UserID: uid,
			Type:   "bank",
			Wallets: []view.ViewUserWalletsGroupByTypeDetailWallet{
				{ID: walletID.String(), Name: "My BCA", Number: "1234567890", Balance: 100000},
			},
		},
	}
	d.walletsRepo.On("GetWalletsByUserIDGroupByType", mock.Anything, nil, uid).Return(grouped, nil)

	result, err := svc.GetWalletsByUserIDGroupByType(context.Background(), uid)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "bank", result[0].Type)
	assert.Len(t, result[0].Wallets, 1)
	d.assertAll(t)
}

func TestGetWalletsByUserIDGroupByType_RepositoryError(t *testing.T) {
	d := newWalletTestDeps()
	svc := d.service()

	uid := userID.String()
	d.walletsRepo.On("GetWalletsByUserIDGroupByType", mock.Anything, nil, uid).
		Return([]view.ViewUserWalletsGroupByType{}, errors.New("db error"))

	result, err := svc.GetWalletsByUserIDGroupByType(context.Background(), uid)

	assert.Error(t, err)
	assert.Nil(t, result)
	d.assertAll(t)
}

// =====================================================================
// CreateWallet
// =====================================================================

func TestCreateWallet_Success(t *testing.T) {
	d := newWalletTestDeps()
	svc := d.service()

	req := sampleWalletRequest()
	wt := sampleWalletType()
	w := sampleWalletModel()

	d.typesRepo.On("GetWalletTypeByID", mock.Anything, nil, req.WalletTypeID).Return(wt, nil)
	d.txManager.On("Begin", mock.Anything).Return(d.tx, nil)
	d.walletsRepo.On("CreateWallet", mock.Anything, d.tx, mock.Anything).Return(w, nil)
	d.txClient.On("InitialDeposit", mock.Anything, mock.AnythingOfType("string"), req.Balance).
		Return(&tpb.TransactionDetail{Id: "tx-123"}, nil)
	d.outboxRepo.On("Create", mock.Anything, d.tx, mock.Anything).Return(nil)
	d.tx.On("Commit").Return(nil)
	d.tx.On("Rollback").Return(nil)

	result, err := svc.CreateWallet(context.Background(), userID.String(), req)

	assert.NoError(t, err)
	assert.Equal(t, w.Name, result.Name)
	assert.Equal(t, w.Balance, result.Balance)
	d.assertAll(t)
}

func TestCreateWallet_InvalidUserID(t *testing.T) {
	d := newWalletTestDeps()
	svc := d.service()

	req := sampleWalletRequest()

	result, err := svc.CreateWallet(context.Background(), "invalid-uuid", req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid user id")
	assert.Empty(t, result.ID)
	d.assertAll(t)
}

func TestCreateWallet_InvalidWalletTypeID(t *testing.T) {
	d := newWalletTestDeps()
	svc := d.service()

	req := sampleWalletRequest()
	req.WalletTypeID = "invalid-uuid"

	result, err := svc.CreateWallet(context.Background(), userID.String(), req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid wallet type id")
	assert.Empty(t, result.ID)
	d.assertAll(t)
}

func TestCreateWallet_WalletTypeNotFound(t *testing.T) {
	d := newWalletTestDeps()
	svc := d.service()

	req := sampleWalletRequest()

	d.typesRepo.On("GetWalletTypeByID", mock.Anything, nil, req.WalletTypeID).
		Return(model.WalletTypes{}, errors.New("record not found"))

	result, err := svc.CreateWallet(context.Background(), userID.String(), req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "wallet type not found")
	assert.Empty(t, result.ID)
	d.assertAll(t)
}

func TestCreateWallet_BeginTxError(t *testing.T) {
	d := newWalletTestDeps()
	svc := d.service()

	req := sampleWalletRequest()
	wt := sampleWalletType()

	d.typesRepo.On("GetWalletTypeByID", mock.Anything, nil, req.WalletTypeID).Return(wt, nil)
	d.txManager.On("Begin", mock.Anything).Return(nil, errors.New("tx error"))

	result, err := svc.CreateWallet(context.Background(), userID.String(), req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "begin transaction")
	assert.Empty(t, result.ID)
	d.assertAll(t)
}

func TestCreateWallet_InsertError(t *testing.T) {
	d := newWalletTestDeps()
	svc := d.service()

	req := sampleWalletRequest()
	wt := sampleWalletType()

	d.typesRepo.On("GetWalletTypeByID", mock.Anything, nil, req.WalletTypeID).Return(wt, nil)
	d.txManager.On("Begin", mock.Anything).Return(d.tx, nil)
	d.walletsRepo.On("CreateWallet", mock.Anything, d.tx, mock.Anything).
		Return(model.Wallets{}, errors.New("insert failed"))
	d.tx.On("Rollback").Return(nil)

	result, err := svc.CreateWallet(context.Background(), userID.String(), req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "insert to db")
	assert.Empty(t, result.ID)
	d.assertAll(t)
}

func TestCreateWallet_GRPCDepositError(t *testing.T) {
	d := newWalletTestDeps()
	svc := d.service()

	req := sampleWalletRequest()
	wt := sampleWalletType()
	w := sampleWalletModel()

	d.typesRepo.On("GetWalletTypeByID", mock.Anything, nil, req.WalletTypeID).Return(wt, nil)
	d.txManager.On("Begin", mock.Anything).Return(d.tx, nil)
	d.walletsRepo.On("CreateWallet", mock.Anything, d.tx, mock.Anything).Return(w, nil)
	d.txClient.On("InitialDeposit", mock.Anything, mock.AnythingOfType("string"), req.Balance).
		Return(nil, errors.New("grpc error"))
	d.tx.On("Rollback").Return(nil)

	result, err := svc.CreateWallet(context.Background(), userID.String(), req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "initial deposit via grpc")
	assert.Empty(t, result.ID)
	d.assertAll(t)
}

func TestCreateWallet_OutboxCreateError(t *testing.T) {
	d := newWalletTestDeps()
	svc := d.service()

	req := sampleWalletRequest()
	wt := sampleWalletType()
	w := sampleWalletModel()

	d.typesRepo.On("GetWalletTypeByID", mock.Anything, nil, req.WalletTypeID).Return(wt, nil)
	d.txManager.On("Begin", mock.Anything).Return(d.tx, nil)
	d.walletsRepo.On("CreateWallet", mock.Anything, d.tx, mock.Anything).Return(w, nil)
	d.txClient.On("InitialDeposit", mock.Anything, mock.AnythingOfType("string"), req.Balance).
		Return(&tpb.TransactionDetail{Id: "tx-123"}, nil)
	d.outboxRepo.On("Create", mock.Anything, d.tx, mock.Anything).Return(errors.New("outbox error"))
	d.tx.On("Rollback").Return(nil)
	// NOTE: CancelInitialDeposit is NOT called because the outbox create uses `:=`
	// which shadows the outer `err`, so the defer sees err == nil.

	result, err := svc.CreateWallet(context.Background(), userID.String(), req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "save outbox message")
	assert.Empty(t, result.ID)
	d.txClient.AssertNotCalled(t, "CancelInitialDeposit")
	d.assertAll(t)
}

func TestCreateWallet_CommitError(t *testing.T) {
	d := newWalletTestDeps()
	svc := d.service()

	req := sampleWalletRequest()
	wt := sampleWalletType()
	w := sampleWalletModel()

	d.typesRepo.On("GetWalletTypeByID", mock.Anything, nil, req.WalletTypeID).Return(wt, nil)
	d.txManager.On("Begin", mock.Anything).Return(d.tx, nil)
	d.walletsRepo.On("CreateWallet", mock.Anything, d.tx, mock.Anything).Return(w, nil)
	d.txClient.On("InitialDeposit", mock.Anything, mock.AnythingOfType("string"), req.Balance).
		Return(&tpb.TransactionDetail{Id: "tx-123"}, nil)
	d.outboxRepo.On("Create", mock.Anything, d.tx, mock.Anything).Return(nil)
	d.tx.On("Commit").Return(errors.New("commit error"))
	d.tx.On("Rollback").Return(nil)
	// NOTE: CancelInitialDeposit is NOT called because commit uses `:=`
	// which shadows the outer `err`, so the defer sees err == nil.

	result, err := svc.CreateWallet(context.Background(), userID.String(), req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "commit transaction")
	assert.Empty(t, result.ID)
	d.txClient.AssertNotCalled(t, "CancelInitialDeposit")
	d.assertAll(t)
}

// =====================================================================
// CreateWalletGRPC
// =====================================================================

func TestCreateWalletGRPC_Success(t *testing.T) {
	d := newWalletTestDeps()
	svc := d.service()

	req := sampleWalletRequest()
	wt := sampleWalletType()
	w := sampleWalletModel()

	d.typesRepo.On("GetWalletTypeByID", mock.Anything, nil, req.WalletTypeID).Return(wt, nil)
	d.txManager.On("Begin", mock.Anything).Return(d.tx, nil)
	d.walletsRepo.On("CreateWallet", mock.Anything, d.tx, mock.Anything).Return(w, nil)
	d.txClient.On("InitialDeposit", mock.Anything, mock.AnythingOfType("string"), req.Balance).
		Return(&tpb.TransactionDetail{Id: "tx-456"}, nil)
	d.outboxRepo.On("Create", mock.Anything, d.tx, mock.Anything).Return(nil)
	d.tx.On("Commit").Return(nil)
	d.tx.On("Rollback").Return(nil)

	result, err := svc.CreateWalletGRPC(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, w.Name, result.Name)
	assert.Equal(t, w.Balance, result.Balance)
	d.assertAll(t)
}

func TestCreateWalletGRPC_SuccessZeroBalance(t *testing.T) {
	d := newWalletTestDeps()
	svc := d.service()

	req := sampleWalletRequest()
	req.Balance = 0

	wt := sampleWalletType()
	w := sampleWalletModel()
	w.Balance = 0

	d.typesRepo.On("GetWalletTypeByID", mock.Anything, nil, req.WalletTypeID).Return(wt, nil)
	d.txManager.On("Begin", mock.Anything).Return(d.tx, nil)
	d.walletsRepo.On("CreateWallet", mock.Anything, d.tx, mock.Anything).Return(w, nil)
	// InitialDeposit should NOT be called when balance is 0
	d.outboxRepo.On("Create", mock.Anything, d.tx, mock.Anything).Return(nil)
	d.tx.On("Commit").Return(nil)
	d.tx.On("Rollback").Return(nil)

	result, err := svc.CreateWalletGRPC(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, float64(0), result.Balance)
	d.txClient.AssertNotCalled(t, "InitialDeposit")
	d.assertAll(t)
}

func TestCreateWalletGRPC_InvalidUserID(t *testing.T) {
	d := newWalletTestDeps()
	svc := d.service()

	req := sampleWalletRequest()
	req.UserID = "not-a-uuid"

	result, err := svc.CreateWalletGRPC(context.Background(), req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid user id")
	assert.Empty(t, result.ID)
	d.assertAll(t)
}

func TestCreateWalletGRPC_InvalidWalletTypeID(t *testing.T) {
	d := newWalletTestDeps()
	svc := d.service()

	req := sampleWalletRequest()
	req.WalletTypeID = "not-a-uuid"

	result, err := svc.CreateWalletGRPC(context.Background(), req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid wallet type id")
	assert.Empty(t, result.ID)
	d.assertAll(t)
}

func TestCreateWalletGRPC_WalletTypeNotFound(t *testing.T) {
	d := newWalletTestDeps()
	svc := d.service()

	req := sampleWalletRequest()

	d.typesRepo.On("GetWalletTypeByID", mock.Anything, nil, req.WalletTypeID).
		Return(model.WalletTypes{}, errors.New("not found"))

	result, err := svc.CreateWalletGRPC(context.Background(), req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "wallet type not found")
	assert.Empty(t, result.ID)
	d.assertAll(t)
}

func TestCreateWalletGRPC_BeginTxError(t *testing.T) {
	d := newWalletTestDeps()
	svc := d.service()

	req := sampleWalletRequest()
	wt := sampleWalletType()

	d.typesRepo.On("GetWalletTypeByID", mock.Anything, nil, req.WalletTypeID).Return(wt, nil)
	d.txManager.On("Begin", mock.Anything).Return(nil, errors.New("tx error"))

	result, err := svc.CreateWalletGRPC(context.Background(), req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "begin transaction")
	assert.Empty(t, result.ID)
	d.assertAll(t)
}

func TestCreateWalletGRPC_InsertError(t *testing.T) {
	d := newWalletTestDeps()
	svc := d.service()

	req := sampleWalletRequest()
	wt := sampleWalletType()

	d.typesRepo.On("GetWalletTypeByID", mock.Anything, nil, req.WalletTypeID).Return(wt, nil)
	d.txManager.On("Begin", mock.Anything).Return(d.tx, nil)
	d.walletsRepo.On("CreateWallet", mock.Anything, d.tx, mock.Anything).
		Return(model.Wallets{}, errors.New("insert failed"))
	d.tx.On("Rollback").Return(nil)

	result, err := svc.CreateWalletGRPC(context.Background(), req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "insert to db")
	assert.Empty(t, result.ID)
	d.assertAll(t)
}

func TestCreateWalletGRPC_GRPCDepositError(t *testing.T) {
	d := newWalletTestDeps()
	svc := d.service()

	req := sampleWalletRequest()
	wt := sampleWalletType()
	w := sampleWalletModel()

	d.typesRepo.On("GetWalletTypeByID", mock.Anything, nil, req.WalletTypeID).Return(wt, nil)
	d.txManager.On("Begin", mock.Anything).Return(d.tx, nil)
	d.walletsRepo.On("CreateWallet", mock.Anything, d.tx, mock.Anything).Return(w, nil)
	d.txClient.On("InitialDeposit", mock.Anything, mock.AnythingOfType("string"), req.Balance).
		Return(nil, errors.New("grpc error"))
	d.tx.On("Rollback").Return(nil)

	result, err := svc.CreateWalletGRPC(context.Background(), req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "initial deposit via grpc")
	assert.Empty(t, result.ID)
	d.assertAll(t)
}

func TestCreateWalletGRPC_OutboxCreateError(t *testing.T) {
	d := newWalletTestDeps()
	svc := d.service()

	req := sampleWalletRequest()
	wt := sampleWalletType()
	w := sampleWalletModel()

	d.typesRepo.On("GetWalletTypeByID", mock.Anything, nil, req.WalletTypeID).Return(wt, nil)
	d.txManager.On("Begin", mock.Anything).Return(d.tx, nil)
	d.walletsRepo.On("CreateWallet", mock.Anything, d.tx, mock.Anything).Return(w, nil)
	d.txClient.On("InitialDeposit", mock.Anything, mock.AnythingOfType("string"), req.Balance).
		Return(&tpb.TransactionDetail{Id: "tx-789"}, nil)
	d.outboxRepo.On("Create", mock.Anything, d.tx, mock.Anything).Return(errors.New("outbox error"))
	d.tx.On("Rollback").Return(nil)
	// CancelInitialDeposit NOT called — `:=` shadows outer `err`

	result, err := svc.CreateWalletGRPC(context.Background(), req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "save outbox message")
	assert.Empty(t, result.ID)
	d.txClient.AssertNotCalled(t, "CancelInitialDeposit")
	d.assertAll(t)
}

func TestCreateWalletGRPC_CommitError(t *testing.T) {
	d := newWalletTestDeps()
	svc := d.service()

	req := sampleWalletRequest()
	wt := sampleWalletType()
	w := sampleWalletModel()

	d.typesRepo.On("GetWalletTypeByID", mock.Anything, nil, req.WalletTypeID).Return(wt, nil)
	d.txManager.On("Begin", mock.Anything).Return(d.tx, nil)
	d.walletsRepo.On("CreateWallet", mock.Anything, d.tx, mock.Anything).Return(w, nil)
	d.txClient.On("InitialDeposit", mock.Anything, mock.AnythingOfType("string"), req.Balance).
		Return(&tpb.TransactionDetail{Id: "tx-789"}, nil)
	d.outboxRepo.On("Create", mock.Anything, d.tx, mock.Anything).Return(nil)
	d.tx.On("Commit").Return(errors.New("commit failed"))
	d.tx.On("Rollback").Return(nil)
	// CancelInitialDeposit NOT called — `:=` shadows outer `err`

	result, err := svc.CreateWalletGRPC(context.Background(), req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "commit transaction")
	assert.Empty(t, result.ID)
	d.txClient.AssertNotCalled(t, "CancelInitialDeposit")
	d.assertAll(t)
}

// =====================================================================
// UpdateWallet
// =====================================================================

func TestUpdateWallet_Success(t *testing.T) {
	d := newWalletTestDeps()
	svc := d.service()

	existing := sampleWalletModel()
	id := existing.ID.String()
	req := dto.WalletsRequest{
		WalletTypeID: walletTypeID.String(),
		Name:         "Updated BCA",
		Number:       "9999999999",
		Balance:      200000,
	}

	updated := existing
	updated.Name = req.Name
	updated.Number = req.Number
	updated.Balance = req.Balance

	d.walletsRepo.On("GetWalletByID", mock.Anything, nil, id).Return(existing, nil)
	d.txManager.On("Begin", mock.Anything).Return(d.tx, nil)
	d.walletsRepo.On("UpdateWallet", mock.Anything, d.tx, mock.Anything).Return(updated, nil)
	d.outboxRepo.On("Create", mock.Anything, d.tx, mock.Anything).Return(nil)
	d.tx.On("Commit").Return(nil)
	d.tx.On("Rollback").Return(nil)

	result, err := svc.UpdateWallet(context.Background(), id, req)

	assert.NoError(t, err)
	assert.Equal(t, req.Name, result.Name)
	assert.Equal(t, req.Number, result.Number)
	assert.Equal(t, req.Balance, result.Balance)
	d.assertAll(t)
}

func TestUpdateWallet_NotFound(t *testing.T) {
	d := newWalletTestDeps()
	svc := d.service()

	id := uuid.New().String()
	req := sampleWalletRequest()

	d.walletsRepo.On("GetWalletByID", mock.Anything, nil, id).
		Return(model.Wallets{}, errors.New("record not found"))

	result, err := svc.UpdateWallet(context.Background(), id, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "wallet not found")
	assert.Empty(t, result.ID)
	d.assertAll(t)
}

func TestUpdateWallet_InvalidWalletTypeID(t *testing.T) {
	d := newWalletTestDeps()
	svc := d.service()

	existing := sampleWalletModel()
	id := existing.ID.String()
	req := dto.WalletsRequest{
		WalletTypeID: "not-valid",
		Name:         "Updated",
		Number:       "123",
		Balance:      100,
	}

	d.walletsRepo.On("GetWalletByID", mock.Anything, nil, id).Return(existing, nil)

	result, err := svc.UpdateWallet(context.Background(), id, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid wallet type id")
	assert.Empty(t, result.ID)
	d.assertAll(t)
}

func TestUpdateWallet_BeginTxError(t *testing.T) {
	d := newWalletTestDeps()
	svc := d.service()

	existing := sampleWalletModel()
	id := existing.ID.String()
	req := dto.WalletsRequest{
		WalletTypeID: walletTypeID.String(),
		Name:         "Updated",
		Number:       "123",
		Balance:      100,
	}

	d.walletsRepo.On("GetWalletByID", mock.Anything, nil, id).Return(existing, nil)
	d.txManager.On("Begin", mock.Anything).Return(nil, errors.New("tx error"))

	result, err := svc.UpdateWallet(context.Background(), id, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "begin transaction")
	assert.Empty(t, result.ID)
	d.assertAll(t)
}

func TestUpdateWallet_UpdateError(t *testing.T) {
	d := newWalletTestDeps()
	svc := d.service()

	existing := sampleWalletModel()
	id := existing.ID.String()
	req := dto.WalletsRequest{
		WalletTypeID: walletTypeID.String(),
		Name:         "Updated",
		Number:       "123",
		Balance:      100,
	}

	d.walletsRepo.On("GetWalletByID", mock.Anything, nil, id).Return(existing, nil)
	d.txManager.On("Begin", mock.Anything).Return(d.tx, nil)
	d.walletsRepo.On("UpdateWallet", mock.Anything, d.tx, mock.Anything).
		Return(model.Wallets{}, errors.New("update failed"))
	d.tx.On("Rollback").Return(nil)

	result, err := svc.UpdateWallet(context.Background(), id, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "update in db")
	assert.Empty(t, result.ID)
	d.assertAll(t)
}

func TestUpdateWallet_OutboxError(t *testing.T) {
	d := newWalletTestDeps()
	svc := d.service()

	existing := sampleWalletModel()
	id := existing.ID.String()
	req := dto.WalletsRequest{
		WalletTypeID: walletTypeID.String(),
		Name:         "Updated",
		Number:       "123",
		Balance:      100,
	}

	updated := existing
	updated.Name = req.Name

	d.walletsRepo.On("GetWalletByID", mock.Anything, nil, id).Return(existing, nil)
	d.txManager.On("Begin", mock.Anything).Return(d.tx, nil)
	d.walletsRepo.On("UpdateWallet", mock.Anything, d.tx, mock.Anything).Return(updated, nil)
	d.outboxRepo.On("Create", mock.Anything, d.tx, mock.Anything).Return(errors.New("outbox error"))
	d.tx.On("Rollback").Return(nil)

	result, err := svc.UpdateWallet(context.Background(), id, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "save outbox message")
	assert.Empty(t, result.ID)
	d.assertAll(t)
}

func TestUpdateWallet_CommitError(t *testing.T) {
	d := newWalletTestDeps()
	svc := d.service()

	existing := sampleWalletModel()
	id := existing.ID.String()
	req := dto.WalletsRequest{
		WalletTypeID: walletTypeID.String(),
		Name:         "Updated",
		Number:       "123",
		Balance:      100,
	}

	updated := existing
	updated.Name = req.Name

	d.walletsRepo.On("GetWalletByID", mock.Anything, nil, id).Return(existing, nil)
	d.txManager.On("Begin", mock.Anything).Return(d.tx, nil)
	d.walletsRepo.On("UpdateWallet", mock.Anything, d.tx, mock.Anything).Return(updated, nil)
	d.outboxRepo.On("Create", mock.Anything, d.tx, mock.Anything).Return(nil)
	d.tx.On("Commit").Return(errors.New("commit error"))
	d.tx.On("Rollback").Return(nil)

	result, err := svc.UpdateWallet(context.Background(), id, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "commit transaction")
	assert.Empty(t, result.ID)
	d.assertAll(t)
}

// =====================================================================
// DeleteWallet
// =====================================================================

func TestDeleteWallet_Success(t *testing.T) {
	d := newWalletTestDeps()
	svc := d.service()

	existing := sampleWalletModel()
	existing.Balance = 0 // must be zero
	id := existing.ID.String()

	d.walletsRepo.On("GetWalletByID", mock.Anything, nil, id).Return(existing, nil)
	d.txManager.On("Begin", mock.Anything).Return(d.tx, nil)
	d.walletsRepo.On("DeleteWallet", mock.Anything, d.tx, existing).Return(existing, nil)
	d.outboxRepo.On("Create", mock.Anything, d.tx, mock.Anything).Return(nil)
	d.tx.On("Commit").Return(nil)
	d.tx.On("Rollback").Return(nil)

	result, err := svc.DeleteWallet(context.Background(), id)

	assert.NoError(t, err)
	assert.Equal(t, id, result.ID)
	d.assertAll(t)
}

func TestDeleteWallet_NotFound(t *testing.T) {
	d := newWalletTestDeps()
	svc := d.service()

	id := uuid.New().String()

	d.walletsRepo.On("GetWalletByID", mock.Anything, nil, id).
		Return(model.Wallets{}, errors.New("record not found"))

	result, err := svc.DeleteWallet(context.Background(), id)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "wallet not found")
	assert.Empty(t, result.ID)
	d.assertAll(t)
}

func TestDeleteWallet_BalanceNotZero(t *testing.T) {
	d := newWalletTestDeps()
	svc := d.service()

	existing := sampleWalletModel()
	existing.Balance = 50000 // not zero
	id := existing.ID.String()

	d.walletsRepo.On("GetWalletByID", mock.Anything, nil, id).Return(existing, nil)

	result, err := svc.DeleteWallet(context.Background(), id)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "wallet balance must be zero")
	assert.Empty(t, result.ID)
	d.assertAll(t)
}

func TestDeleteWallet_BeginTxError(t *testing.T) {
	d := newWalletTestDeps()
	svc := d.service()

	existing := sampleWalletModel()
	existing.Balance = 0
	id := existing.ID.String()

	d.walletsRepo.On("GetWalletByID", mock.Anything, nil, id).Return(existing, nil)
	d.txManager.On("Begin", mock.Anything).Return(nil, errors.New("tx error"))

	result, err := svc.DeleteWallet(context.Background(), id)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "begin transaction")
	assert.Empty(t, result.ID)
	d.assertAll(t)
}

func TestDeleteWallet_DeleteError(t *testing.T) {
	d := newWalletTestDeps()
	svc := d.service()

	existing := sampleWalletModel()
	existing.Balance = 0
	id := existing.ID.String()

	d.walletsRepo.On("GetWalletByID", mock.Anything, nil, id).Return(existing, nil)
	d.txManager.On("Begin", mock.Anything).Return(d.tx, nil)
	d.walletsRepo.On("DeleteWallet", mock.Anything, d.tx, existing).
		Return(model.Wallets{}, errors.New("delete failed"))
	d.tx.On("Rollback").Return(nil)

	result, err := svc.DeleteWallet(context.Background(), id)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "delete from db")
	assert.Empty(t, result.ID)
	d.assertAll(t)
}

func TestDeleteWallet_OutboxError(t *testing.T) {
	d := newWalletTestDeps()
	svc := d.service()

	existing := sampleWalletModel()
	existing.Balance = 0
	id := existing.ID.String()

	d.walletsRepo.On("GetWalletByID", mock.Anything, nil, id).Return(existing, nil)
	d.txManager.On("Begin", mock.Anything).Return(d.tx, nil)
	d.walletsRepo.On("DeleteWallet", mock.Anything, d.tx, existing).Return(existing, nil)
	d.outboxRepo.On("Create", mock.Anything, d.tx, mock.Anything).Return(errors.New("outbox error"))
	d.tx.On("Rollback").Return(nil)

	result, err := svc.DeleteWallet(context.Background(), id)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "save outbox message")
	assert.Empty(t, result.ID)
	d.assertAll(t)
}

func TestDeleteWallet_CommitError(t *testing.T) {
	d := newWalletTestDeps()
	svc := d.service()

	existing := sampleWalletModel()
	existing.Balance = 0
	id := existing.ID.String()

	d.walletsRepo.On("GetWalletByID", mock.Anything, nil, id).Return(existing, nil)
	d.txManager.On("Begin", mock.Anything).Return(d.tx, nil)
	d.walletsRepo.On("DeleteWallet", mock.Anything, d.tx, existing).Return(existing, nil)
	d.outboxRepo.On("Create", mock.Anything, d.tx, mock.Anything).Return(nil)
	d.tx.On("Commit").Return(errors.New("commit error"))
	d.tx.On("Rollback").Return(nil)

	result, err := svc.DeleteWallet(context.Background(), id)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "commit transaction")
	assert.Empty(t, result.ID)
	d.assertAll(t)
}
