package view

type ViewUserWallets struct {
	ID             string  `json:"id"`
	UserID         string  `json:"user_id"`
	WalletNumber   string  `json:"wallet_number"`
	WalletBalance  float64 `json:"wallet_balance"`
	WalletName     string  `json:"wallet_name"`
	WalletTypeName string  `json:"wallet_type_name"`
	WalletType     string  `json:"wallet_type"`
}

type ViewUserWalletsGroupByTypeDetailWallet struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Number  string  `json:"number"`
	Balance float64 `json:"balance"`
}

type ViewUserWalletsGroupByType struct {
	UserID  string                                   `json:"user_id"`
	Type    string                                   `json:"type"`
	Wallets []ViewUserWalletsGroupByTypeDetailWallet `gorm:"type:jsonb" json:"wallets"`
}
