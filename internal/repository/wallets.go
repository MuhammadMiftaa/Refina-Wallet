package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"refina-wallet/internal/types/model"
	"refina-wallet/internal/types/view"

	"gorm.io/gorm"
)

type WalletsRepository interface {
	GetAllWallets(ctx context.Context, tx Transaction) ([]model.Wallets, error)
	GetWalletByID(ctx context.Context, tx Transaction, id string) (model.Wallets, error)
	GetWalletsByUserID(ctx context.Context, tx Transaction, id string) ([]model.Wallets, error)
	GetWalletsByUserIDGroupByType(ctx context.Context, tx Transaction, id string) ([]view.ViewUserWalletsGroupByType, error)
	CreateWallet(ctx context.Context, tx Transaction, wallet model.Wallets) (model.Wallets, error)
	UpdateWallet(ctx context.Context, tx Transaction, wallet model.Wallets) (model.Wallets, error)
	DeleteWallet(ctx context.Context, tx Transaction, wallet model.Wallets) (model.Wallets, error)
}

type walletsRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) WalletsRepository {
	return &walletsRepository{db}
}

// Helper untuk mendapatkan DB instance (transaksi atau biasa)
func (wallet_repo *walletsRepository) getDB(ctx context.Context, tx Transaction) (*gorm.DB, error) {
	if tx != nil {
		gormTx, ok := tx.(*GormTx) // Type assertion ke GORM transaction
		if !ok {
			return nil, errors.New("invalid transaction type")
		}
		return gormTx.db.WithContext(ctx), nil
	}
	return wallet_repo.db.WithContext(ctx), nil
}

// Implementasi method dengan transaksi opsional
func (wallet_repo *walletsRepository) GetAllWallets(ctx context.Context, tx Transaction) ([]model.Wallets, error) {
	db, err := wallet_repo.getDB(ctx, tx)
	if err != nil {
		return nil, err
	}

	var wallets []model.Wallets
	if err := db.Find(&wallets).Error; err != nil {
		return nil, err
	}
	return wallets, nil
}

func (wallet_repo *walletsRepository) GetWalletByID(ctx context.Context, tx Transaction, id string) (model.Wallets, error) {
	db, err := wallet_repo.getDB(ctx, tx)
	if err != nil {
		return model.Wallets{}, err
	}

	var wallet model.Wallets
	if err := db.Where("id = ?", id).First(&wallet).Error; err != nil {
		return model.Wallets{}, err
	}
	return wallet, nil
}

func (wallet_repo *walletsRepository) GetWalletsByUserID(ctx context.Context, tx Transaction, id string) ([]model.Wallets, error) {
	db, err := wallet_repo.getDB(ctx, tx)
	if err != nil {
		return nil, err
	}

	var userWallets []model.Wallets
	err = db.Where("user_id = ?", id).Find(&userWallets).Error
	if err != nil {
		return nil, errors.New("user wallets not found")
	}

	if len(userWallets) == 0 {
		return nil, nil
	}

	return userWallets, nil
}

func (wallet_repo *walletsRepository) GetWalletsByUserIDGroupByType(ctx context.Context, tx Transaction, id string) ([]view.ViewUserWalletsGroupByType, error) {
	db, err := wallet_repo.getDB(ctx, tx)
	if err != nil {
		return nil, err
	}

	var rawResults []struct {
		UserID  string
		Type    string
		Wallets []byte
	}
	err = db.Raw(`SELECT * FROM view_user_wallets_group_by_type WHERE user_id = $1`, id).Scan(&rawResults).Error
	if err != nil {
		return nil, errors.New("user wallets group by type not found")
	}

	var results []view.ViewUserWalletsGroupByType

	for _, row := range rawResults {
		var wallets []view.ViewUserWalletsGroupByTypeDetailWallet

		err := json.Unmarshal(row.Wallets, &wallets)
		if err != nil {
			return nil, fmt.Errorf("gagal decode JSON wallets (type: %s): %w", row.Type, err)
		}

		results = append(results, view.ViewUserWalletsGroupByType{
			UserID:  row.UserID,
			Type:    row.Type,
			Wallets: wallets,
		})
	}

	return results, nil
}

func (wallet_repo *walletsRepository) CreateWallet(ctx context.Context, tx Transaction, wallet model.Wallets) (model.Wallets, error) {
	db, err := wallet_repo.getDB(ctx, tx)
	if err != nil {
		return model.Wallets{}, err
	}

	if err := db.Create(&wallet).Error; err != nil {
		return model.Wallets{}, err
	}

	return wallet, nil
}

func (wallet_repo *walletsRepository) UpdateWallet(ctx context.Context, tx Transaction, wallet model.Wallets) (model.Wallets, error) {
	db, err := wallet_repo.getDB(ctx, tx)
	if err != nil {
		return model.Wallets{}, err
	}

	if err := db.Omit("User", "WalletType").Save(&wallet).Error; err != nil {
		return model.Wallets{}, err
	}

	return wallet, nil
}

func (wallet_repo *walletsRepository) DeleteWallet(ctx context.Context, tx Transaction, wallet model.Wallets) (model.Wallets, error) {
	db, err := wallet_repo.getDB(ctx, tx)
	if err != nil {
		return model.Wallets{}, err
	}

	if err := db.Delete(&wallet).Error; err != nil {
		return model.Wallets{}, err
	}

	return wallet, nil
}
