package dto

type WalletsResponse struct {
	ID                    string  `json:"id"`
	UserID                string  `json:"user_id"`
	WalletTypeID          string  `json:"wallet_type_id"`
	WalletType            string  `json:"wallet_type"`
	WalletTypeName        string  `json:"wallet_type_name"`
	WalletTypeDescription string  `json:"wallet_type_description"`
	Name                  string  `json:"name"`
	Number                string  `json:"number"`
	Balance               float64 `json:"balance"`
	CreatedAt             string  `json:"created_at"`
	UpdatedAt             string  `json:"updated_at"`
}

type WalletsRequest struct {
	UserID       string  `json:"user_id"`
	WalletTypeID string  `json:"wallet_type_id"`
	Name         string  `json:"name"`
	Number       string  `json:"number"`
	Balance      float64 `json:"balance"`
}
