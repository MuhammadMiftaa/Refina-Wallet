package dto

type CategoryType string

const (
	Income       CategoryType = "income"
	Expense      CategoryType = "expense"
	FundTransfer CategoryType = "fund_transfer"
)