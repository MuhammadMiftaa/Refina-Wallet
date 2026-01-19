package service

import (
	"context"

	"refina-wallet/internal/repository"
	"refina-wallet/internal/types/dto"
	"refina-wallet/internal/types/model"
	"refina-wallet/internal/utils"
)

type WalletTypesService interface {
	GetAllWalletTypes(ctx context.Context) ([]dto.WalletTypesResponse, error)
	GetWalletTypeByID(ctx context.Context, id string) (dto.WalletTypesResponse, error)
	CreateWalletType(ctx context.Context, walletType dto.WalletTypesRequest) (dto.WalletTypesResponse, error)
	UpdateWalletType(ctx context.Context, id string, walletType dto.WalletTypesRequest) (dto.WalletTypesResponse, error)
	DeleteWalletType(ctx context.Context, id string) (dto.WalletTypesResponse, error)
}

type walletTypesService struct {
	txManager       repository.TxManager
	walletTypesRepo repository.WalletTypesRepository
}

func NewWalletTypesService(txManager repository.TxManager, walletTypesRepo repository.WalletTypesRepository) WalletTypesService {
	return &walletTypesService{
		txManager:       txManager,
		walletTypesRepo: walletTypesRepo,
	}
}

func (walletTypeServ *walletTypesService) GetAllWalletTypes(ctx context.Context) ([]dto.WalletTypesResponse, error) {
	walletTypes, err := walletTypeServ.walletTypesRepo.GetAllWalletTypes(ctx, nil)
	if err != nil {
		return nil, err
	}

	var walletTypesResponse []dto.WalletTypesResponse
	for _, walletType := range walletTypes {
		walletTypeResponse := utils.ConvertToResponseType(walletType).(dto.WalletTypesResponse)
		walletTypesResponse = append(walletTypesResponse, walletTypeResponse)
	}

	return walletTypesResponse, nil
}

func (walletTypeServ *walletTypesService) GetWalletTypeByID(ctx context.Context, id string) (dto.WalletTypesResponse, error) {
	walletType, err := walletTypeServ.walletTypesRepo.GetWalletTypeByID(ctx, nil, id)
	if err != nil {
		return dto.WalletTypesResponse{}, err
	}

	walletTypeResponse := utils.ConvertToResponseType(walletType).(dto.WalletTypesResponse)

	return walletTypeResponse, nil
}

func (walletTypeServ *walletTypesService) CreateWalletType(ctx context.Context, walletType dto.WalletTypesRequest) (dto.WalletTypesResponse, error) {
	walletTypeModel := model.WalletTypes{
		Name:        walletType.Name,
		Type:        model.WalletType(walletType.Type),
		Description: walletType.Description,
	}

	walletTypeModel, err := walletTypeServ.walletTypesRepo.CreateWalletType(ctx, nil, walletTypeModel)
	if err != nil {
		return dto.WalletTypesResponse{}, err
	}

	walletTypeResponse := utils.ConvertToResponseType(walletTypeModel).(dto.WalletTypesResponse)

	return walletTypeResponse, nil
}

func (walletTypeServ *walletTypesService) UpdateWalletType(ctx context.Context, id string, walletType dto.WalletTypesRequest) (dto.WalletTypesResponse, error) {
	walletTypeModel := model.WalletTypes{
		Name:        walletType.Name,
		Type:        model.WalletType(walletType.Type),
		Description: walletType.Description,
	}

	walletTypeModel, err := walletTypeServ.walletTypesRepo.UpdateWalletType(ctx, nil, walletTypeModel)
	if err != nil {
		return dto.WalletTypesResponse{}, err
	}

	walletTypeResponse := utils.ConvertToResponseType(walletTypeModel).(dto.WalletTypesResponse)

	return walletTypeResponse, nil
}

func (walletTypeServ *walletTypesService) DeleteWalletType(ctx context.Context, id string) (dto.WalletTypesResponse, error) {
	walletTypeModel, err := walletTypeServ.walletTypesRepo.GetWalletTypeByID(ctx, nil, id)
	if err != nil {
		return dto.WalletTypesResponse{}, err
	}

	walletTypeModel, err = walletTypeServ.walletTypesRepo.DeleteWalletType(ctx, nil, walletTypeModel)
	if err != nil {
		return dto.WalletTypesResponse{}, err
	}

	walletTypeResponse := utils.ConvertToResponseType(walletTypeModel).(dto.WalletTypesResponse)

	return walletTypeResponse, nil
}
