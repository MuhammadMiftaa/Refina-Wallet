package model

type WalletType string

const (
	Bank         WalletType = "bank"
	EWallet      WalletType = "e-wallet"
	Physical     WalletType = "physical"
	OthersWallet WalletType = "others"
)

type WalletTypes struct {
	Base
	Name        string     `gorm:"type:varchar(50);not null"`
	Type        WalletType `gorm:"type:varchar(50);not null"`
	Description string     `gorm:"type:text"`
}