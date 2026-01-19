package dto

type WalletType string

const (
	Bank         WalletType = "bank"
	EWallet      WalletType = "e-wallet"
	Physical     WalletType = "physical"
	OthersWallet WalletType = "others"
)

type WalletTypesResponse struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Type        WalletType `json:"type"`
	Description string     `json:"description"`
}

type WalletTypesRequest struct {
	Name        string     `json:"name"`
	Type        WalletType `json:"type"`
	Description string     `json:"description"`
}
