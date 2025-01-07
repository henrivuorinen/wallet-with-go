package models

type TransactionType string

const (
	TransactionTypePurchase TransactionType = "purchase"
	TransactionTypeWin      TransactionType = "win"
)

type Transaction struct {
	ID        string          `json:"id"`
	PlayerID  string          `json:"playerId"`
	Timestamp string          `json:"timestamp"`
	Type      TransactionType `json:"type"`
	Amount    float64         `json:"amount"`
}
