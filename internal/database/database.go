package database

import (
	"database/sql"
	_ "database/sql"
	"golang.org/x/crypto/bcrypt"
	_ "log"

	_ "github.com/mattn/go-sqlite3"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func InitDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./wallet.db")
	if err != nil {
		return nil, err
	}

	//create tables
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS players(
		    player_id SERIAL PRIMARY KEY,
		    username TEXT NOT NULL,
		    password TEXT NOT NULL,
		    name TEXT NOT NULL,
		    balance DECIMAL(10,2) NOT NULL DEFAULT 0
		);
		CREATE TABLE IF NOT EXISTS transaction_log(
		    id TEXT PRIMARY KEY,
		    player_id TEXT NOT NULL,
		    timestamp DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		    type TEXT NOT NULL,
		    amount DECIMAL(10,2) NOT NULL DEFAULT 0,
		    FOREIGN KEY(player_id) REFERENCES players(id)
		);
	`)
	if err != nil {
		return nil, err
	}
	return db, nil
}
