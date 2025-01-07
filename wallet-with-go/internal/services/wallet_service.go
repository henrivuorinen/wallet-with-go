package services

import (
	"database/sql"
	"errors"
	"wallet-with-go/internal/database"
	"wallet-with-go/internal/models"
)

type WalletService struct {
	DB *sql.DB
}

// RegisterPlayer registers a new player in the database
func (ws *WalletService) RegisterPlayer(username, password, name string, initialBalance float64) error {
	// Hash the password
	hashedPassword, err := database.HashPassword(password)
	if err != nil {
		return err
	}

	// Insert user into the database
	_, err = ws.DB.Exec(
		"INSERT INTO players (username, password, name, balance) VALUES (?, ?, ?, ?)",
		username, hashedPassword, name, initialBalance,
	)
	return err
}

func (ws *WalletService) SavePlayer(playerID, username, hashedPassword, name string, initialBalance float64) error {
	_, err := ws.DB.Exec(
		`INSERT INTO players (player_id, username, password, name, balance) 
		 VALUES (?, ?, ?, ?, ?)`,
		playerID, username, hashedPassword, name, initialBalance,
	)
	if err != nil {
		return err
	}

	return nil
}

// UsernameExists checks if a username already exists in the database
func (ws *WalletService) UsernameExists(username string) (bool, error) {
	var count int
	err := ws.DB.QueryRow("SELECT COUNT(*) FROM players WHERE username = ?", username).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Login validates the player's credentials and retrieves their details
func (ws *WalletService) Login(username, password string) (string, float64, error) {
	var playerID string
	var hashedPassword string
	var balance float64

	err := ws.DB.QueryRow(
		"SELECT player_id, password, balance FROM players WHERE username = ?",
		username,
	).Scan(&playerID, &hashedPassword, &balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", 0, errors.New("invalid username or password")
		}
		return "", 0, err
	}

	if !database.CheckPasswordHash(password, hashedPassword) {
		return "", 0, errors.New("invalid username or password")
	}

	return playerID, balance, nil
}

// ProcessPurchase deducts the specified amount from the player's balance
func (ws *WalletService) ProcessPurchase(playerID, transactionID string, amount float64) (float64, error) {
	if amount <= 0 {
		return 0, errors.New("invalid amount")
	}

	tx, err := ws.DB.Begin()
	if err != nil {
		return 0, err
	}

	var currentBalance float64
	err = tx.QueryRow("SELECT balance FROM players WHERE player_id = ?", playerID).Scan(&currentBalance)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if currentBalance < amount {
		tx.Rollback()
		return 0, errors.New("insufficient funds")
	}

	_, err = tx.Exec("UPDATE players SET balance = balance - ? WHERE player_id = ?", amount, playerID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	_, err = tx.Exec("INSERT INTO transactions (id, player_id, type, amount) VALUES (?, ?, ?, ?)",
		transactionID, playerID, models.TransactionTypePurchase, amount)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return currentBalance - amount, nil
}

// ProcessWin adds the specified amount to the player's balance
func (ws *WalletService) ProcessWin(playerID, transactionID string, amount float64) (float64, error) {
	if amount <= 0 {
		return 0, errors.New("invalid amount")
	}

	tx, err := ws.DB.Begin()
	if err != nil {
		return 0, err
	}

	var currentBalance float64
	err = tx.QueryRow("SELECT balance FROM players WHERE player_id = ?", playerID).Scan(&currentBalance)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	_, err = tx.Exec("UPDATE players SET balance = balance + ? WHERE player_id = ?", amount, playerID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	_, err = tx.Exec("INSERT INTO transactions (id, player_id, type, amount) VALUES (?, ?, ?, ?)",
		transactionID, playerID, models.TransactionTypeWin, amount)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return currentBalance + amount, nil
}
