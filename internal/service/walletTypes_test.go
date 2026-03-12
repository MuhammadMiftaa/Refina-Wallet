package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"refina-wallet/internal/service/mocks"
	"refina-wallet/internal/types/dto"
	"refina-wallet/internal/types/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ---------- helpers ----------

func newWalletTypesService(
	txManager *mocks.MockTxManager,
	walletTypesRepo *mocks.MockWalletTypesRepository,
) WalletTypesService {
	return NewWalletTypesService(txManager, walletTypesRepo)
}

func sampleWalletTypeModel() model.WalletTypes {
	return model.WalletTypes{
		Base: model.Base{
			ID:        uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		Name:        "BCA",
		Type:        model.Bank,
		Description: "Bank BCA",
	}
}

func sampleWalletTypeRequest() dto.WalletTypesRequest {
	return dto.WalletTypesRequest{
		Name:        "BCA",
		Type:        dto.Bank,
		Description: "Bank BCA",
	}
}

// =====================================================================
// GetAllWalletTypes
// =====================================================================

func TestGetAllWalletTypes_Success(t *testing.T) {
	txMgr := new(mocks.MockTxManager)
	repo := new(mocks.MockWalletTypesRepository)

	svc := newWalletTypesService(txMgr, repo)

	wt := sampleWalletTypeModel()
	repo.On("GetAllWalletTypes", mock.Anything, nil).Return([]model.WalletTypes{wt}, nil)

	result, err := svc.GetAllWalletTypes(context.Background())

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, wt.ID.String(), result[0].ID)
	assert.Equal(t, wt.Name, result[0].Name)
	assert.Equal(t, dto.WalletType(wt.Type), result[0].Type)
	assert.Equal(t, wt.Description, result[0].Description)
	repo.AssertExpectations(t)
}

func TestGetAllWalletTypes_EmptyList(t *testing.T) {
	txMgr := new(mocks.MockTxManager)
	repo := new(mocks.MockWalletTypesRepository)

	svc := newWalletTypesService(txMgr, repo)

	repo.On("GetAllWalletTypes", mock.Anything, nil).Return([]model.WalletTypes{}, nil)

	result, err := svc.GetAllWalletTypes(context.Background())

	assert.NoError(t, err)
	assert.Empty(t, result)
	repo.AssertExpectations(t)
}

func TestGetAllWalletTypes_RepositoryError(t *testing.T) {
	txMgr := new(mocks.MockTxManager)
	repo := new(mocks.MockWalletTypesRepository)

	svc := newWalletTypesService(txMgr, repo)

	repo.On("GetAllWalletTypes", mock.Anything, nil).
		Return([]model.WalletTypes{}, errors.New("db error"))

	result, err := svc.GetAllWalletTypes(context.Background())

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "get all wallet types")
	repo.AssertExpectations(t)
}

// =====================================================================
// GetWalletTypeByID
// =====================================================================

func TestGetWalletTypeByID_Success(t *testing.T) {
	txMgr := new(mocks.MockTxManager)
	repo := new(mocks.MockWalletTypesRepository)

	svc := newWalletTypesService(txMgr, repo)

	wt := sampleWalletTypeModel()
	id := wt.ID.String()
	repo.On("GetWalletTypeByID", mock.Anything, nil, id).Return(wt, nil)

	result, err := svc.GetWalletTypeByID(context.Background(), id)

	assert.NoError(t, err)
	assert.Equal(t, id, result.ID)
	assert.Equal(t, wt.Name, result.Name)
	repo.AssertExpectations(t)
}

func TestGetWalletTypeByID_NotFound(t *testing.T) {
	txMgr := new(mocks.MockTxManager)
	repo := new(mocks.MockWalletTypesRepository)

	svc := newWalletTypesService(txMgr, repo)

	id := uuid.New().String()
	repo.On("GetWalletTypeByID", mock.Anything, nil, id).
		Return(model.WalletTypes{}, errors.New("record not found"))

	result, err := svc.GetWalletTypeByID(context.Background(), id)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "wallet type not found")
	assert.Empty(t, result.ID)
	repo.AssertExpectations(t)
}

// =====================================================================
// CreateWalletType
// =====================================================================

func TestCreateWalletType_Success(t *testing.T) {
	txMgr := new(mocks.MockTxManager)
	repo := new(mocks.MockWalletTypesRepository)

	svc := newWalletTypesService(txMgr, repo)

	req := sampleWalletTypeRequest()
	created := sampleWalletTypeModel()

	repo.On("CreateWalletType", mock.Anything, nil, mock.MatchedBy(func(wt model.WalletTypes) bool {
		return wt.Name == req.Name && wt.Type == model.WalletType(req.Type) && wt.Description == req.Description
	})).Return(created, nil)

	result, err := svc.CreateWalletType(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, created.ID.String(), result.ID)
	assert.Equal(t, req.Name, result.Name)
	assert.Equal(t, req.Type, result.Type)
	repo.AssertExpectations(t)
}

func TestCreateWalletType_RepositoryError(t *testing.T) {
	txMgr := new(mocks.MockTxManager)
	repo := new(mocks.MockWalletTypesRepository)

	svc := newWalletTypesService(txMgr, repo)

	req := sampleWalletTypeRequest()

	repo.On("CreateWalletType", mock.Anything, nil, mock.Anything).
		Return(model.WalletTypes{}, errors.New("db error"))

	result, err := svc.CreateWalletType(context.Background(), req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "create wallet type")
	assert.Empty(t, result.ID)
	repo.AssertExpectations(t)
}

// =====================================================================
// UpdateWalletType
// =====================================================================

func TestUpdateWalletType_Success(t *testing.T) {
	txMgr := new(mocks.MockTxManager)
	repo := new(mocks.MockWalletTypesRepository)

	svc := newWalletTypesService(txMgr, repo)

	id := uuid.New().String()
	req := dto.WalletTypesRequest{
		Name:        "Updated Name",
		Type:        dto.EWallet,
		Description: "Updated Description",
	}

	updated := model.WalletTypes{
		Base: model.Base{
			ID:        uuid.MustParse(id),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:        req.Name,
		Type:        model.WalletType(req.Type),
		Description: req.Description,
	}

	repo.On("UpdateWalletType", mock.Anything, nil, mock.MatchedBy(func(wt model.WalletTypes) bool {
		return wt.Name == req.Name && wt.Type == model.WalletType(req.Type)
	})).Return(updated, nil)

	result, err := svc.UpdateWalletType(context.Background(), id, req)

	assert.NoError(t, err)
	assert.Equal(t, req.Name, result.Name)
	assert.Equal(t, req.Type, result.Type)
	repo.AssertExpectations(t)
}

func TestUpdateWalletType_RepositoryError(t *testing.T) {
	txMgr := new(mocks.MockTxManager)
	repo := new(mocks.MockWalletTypesRepository)

	svc := newWalletTypesService(txMgr, repo)

	id := uuid.New().String()
	req := sampleWalletTypeRequest()

	repo.On("UpdateWalletType", mock.Anything, nil, mock.Anything).
		Return(model.WalletTypes{}, errors.New("db error"))

	result, err := svc.UpdateWalletType(context.Background(), id, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "update wallet type")
	assert.Empty(t, result.ID)
	repo.AssertExpectations(t)
}

// =====================================================================
// DeleteWalletType
// =====================================================================

func TestDeleteWalletType_Success(t *testing.T) {
	txMgr := new(mocks.MockTxManager)
	repo := new(mocks.MockWalletTypesRepository)

	svc := newWalletTypesService(txMgr, repo)

	wt := sampleWalletTypeModel()
	id := wt.ID.String()

	repo.On("GetWalletTypeByID", mock.Anything, nil, id).Return(wt, nil)
	repo.On("DeleteWalletType", mock.Anything, nil, wt).Return(wt, nil)

	result, err := svc.DeleteWalletType(context.Background(), id)

	assert.NoError(t, err)
	assert.Equal(t, id, result.ID)
	assert.Equal(t, wt.Name, result.Name)
	repo.AssertExpectations(t)
}

func TestDeleteWalletType_NotFound(t *testing.T) {
	txMgr := new(mocks.MockTxManager)
	repo := new(mocks.MockWalletTypesRepository)

	svc := newWalletTypesService(txMgr, repo)

	id := uuid.New().String()

	repo.On("GetWalletTypeByID", mock.Anything, nil, id).
		Return(model.WalletTypes{}, errors.New("record not found"))

	result, err := svc.DeleteWalletType(context.Background(), id)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "wallet type not found")
	assert.Empty(t, result.ID)
	repo.AssertExpectations(t)
}

func TestDeleteWalletType_DeleteError(t *testing.T) {
	txMgr := new(mocks.MockTxManager)
	repo := new(mocks.MockWalletTypesRepository)

	svc := newWalletTypesService(txMgr, repo)

	wt := sampleWalletTypeModel()
	id := wt.ID.String()

	repo.On("GetWalletTypeByID", mock.Anything, nil, id).Return(wt, nil)
	repo.On("DeleteWalletType", mock.Anything, nil, wt).
		Return(model.WalletTypes{}, errors.New("delete failed"))

	result, err := svc.DeleteWalletType(context.Background(), id)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "delete wallet type")
	assert.Empty(t, result.ID)
	repo.AssertExpectations(t)
}
