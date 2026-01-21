package view

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
