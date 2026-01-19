package repository

import (
	"context"
	"errors"

	"refina-wallet/internal/types/model"

	"gorm.io/gorm"
)

type WalletTypesRepository interface {
	GetAllWalletTypes(ctx context.Context, tx Transaction) ([]model.WalletTypes, error)
	GetWalletTypeByID(ctx context.Context, tx Transaction, id string) (model.WalletTypes, error)
	CreateWalletType(ctx context.Context, tx Transaction, walletType model.WalletTypes) (model.WalletTypes, error)
	UpdateWalletType(ctx context.Context, tx Transaction, walletType model.WalletTypes) (model.WalletTypes, error)
	DeleteWalletType(ctx context.Context, tx Transaction, walletType model.WalletTypes) (model.WalletTypes, error)
}
type walletTypesRepository struct {
	db *gorm.DB
}

func NewWalletTypesRepository(db *gorm.DB) WalletTypesRepository {
	return &walletTypesRepository{db}
}

// Helper untuk mendapatkan DB instance (transaksi atau biasa)
func (wallet_type_repo *walletTypesRepository) getDB(ctx context.Context, tx Transaction) (*gorm.DB, error) {
	if tx != nil {
		gormTx, ok := tx.(*GormTx) // Type assertion ke GORM transaction
		if !ok {
			return nil, errors.New("invalid transaction type")
		}
		return gormTx.db.WithContext(ctx), nil
	}
	return wallet_type_repo.db.WithContext(ctx), nil
}

func (wallet_type_repo *walletTypesRepository) GetAllWalletTypes(ctx context.Context, tx Transaction) ([]model.WalletTypes, error) {
	db, err := wallet_type_repo.getDB(ctx, tx)
	if err != nil {
		return nil, err
	}

	var walletTypes []model.WalletTypes
	if err := db.Find(&walletTypes).Error; err != nil {
		return nil, err
	}
	return walletTypes, nil
}

func (wallet_type_repo *walletTypesRepository) GetWalletTypeByID(ctx context.Context, tx Transaction, id string) (model.WalletTypes, error) {
	db, err := wallet_type_repo.getDB(ctx, tx)
	if err != nil {
		return model.WalletTypes{}, err
	}

	var walletType model.WalletTypes
	if err := db.Where("id = ?", id).First(&walletType).Error; err != nil {
		return model.WalletTypes{}, err
	}
	return walletType, nil
}

func (wallet_type_repo *walletTypesRepository) CreateWalletType(ctx context.Context, tx Transaction, walletType model.WalletTypes) (model.WalletTypes, error) {
	db, err := wallet_type_repo.getDB(ctx, tx)
	if err != nil {
		return model.WalletTypes{}, err
	}

	if err := db.Create(&walletType).Error; err != nil {
		return model.WalletTypes{}, err
	}
	return walletType, nil
}

func (wallet_type_repo *walletTypesRepository) UpdateWalletType(ctx context.Context, tx Transaction, walletType model.WalletTypes) (model.WalletTypes, error) {
	db, err := wallet_type_repo.getDB(ctx, tx)
	if err != nil {
		return model.WalletTypes{}, err
	}

	if err := db.Save(&walletType).Error; err != nil {
		return model.WalletTypes{}, err
	}
	return walletType, nil
}

func (wallet_type_repo *walletTypesRepository) DeleteWalletType(ctx context.Context, tx Transaction, walletType model.WalletTypes) (model.WalletTypes, error) {
	db, err := wallet_type_repo.getDB(ctx, tx)
	if err != nil {
		return model.WalletTypes{}, err
	}

	if err := db.Delete(&walletType).Error; err != nil {
		return model.WalletTypes{}, err
	}
	return walletType, nil
}
