package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"wallet-with-go/internal/handlers"
	"wallet-with-go/internal/services"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3" // SQLite driver for in-memory testing
	"github.com/stretchr/testify/assert"
)

// Initialize an in-memory SQLite database for testing
func initTestDB() *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatalf("Failed to initialize test database: %v", err)
	}

	// Create necessary schema for the tests
	_, err = db.Exec(`
		CREATE TABLE players (
			player_id TEXT PRIMARY KEY,
			username TEXT UNIQUE,
			password TEXT,
			name TEXT,
			balance REAL
		);

		CREATE TABLE transactions (
			id TEXT PRIMARY KEY,
			player_id TEXT,
			type TEXT,
			amount REAL,
			FOREIGN KEY(player_id) REFERENCES players(player_id)
		);
	`)
	if err != nil {
		log.Fatalf("Failed to create test database schema: %v", err)
	}

	return db
}

func setupTestRouter(db *sql.DB) *gin.Engine {
	// Initialize the WalletService with the test DB connection
	walletService := &services.WalletService{DB: db}

	// Create the WalletHandler
	walletHandler := &handlers.WalletHandler{Service: walletService}

	// Set up Gin router
	router := gin.Default()
	walletHandler.SetupRoutes(router)

	return router
}

func TestRegisterAPI(t *testing.T) {
	db := initTestDB()
	defer db.Close()

	router := setupTestRouter(db)

	// Prepare request payload
	requestBody, _ := json.Marshal(map[string]interface{}{
		"username": "testuser",
		"password": "password123",
		"name":     "Test User",
		"balance":  100.0,
	})

	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, resp.Code)
	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Contains(t, response, "player_id")
	assert.Equal(t, "Registration successful", response["message"])
}

func TestLoginAPI(t *testing.T) {
	db := initTestDB()
	defer db.Close()

	// Insert a test user
	_, err := db.Exec(`
		INSERT INTO players (player_id, username, password, name, balance) 
		VALUES ('player123', 'testuser', '$2a$10$8l8VJiKy8E.ciy.dSPKUOeFdIzs8.zP6cEpyJqc/qNNmW4SPvQ.aW', 'Test User', 100.0)
	`)
	if err != nil {
		t.Fatalf("Failed to insert test user: %v", err)
	}

	router := setupTestRouter(db)

	// Prepare request payload
	requestBody, _ := json.Marshal(map[string]string{
		"username": "testuser",
		"password": "password123",
	})

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, resp.Code)
	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Contains(t, response, "player_id")
	assert.Contains(t, response, "balance")
}

func TestPurchaseAPI(t *testing.T) {
	db := initTestDB()
	defer db.Close()

	// Insert a test user
	_, err := db.Exec(`
		INSERT INTO players (player_id, username, password, name, balance) 
		VALUES ('player123', 'testuser', 'hashedpassword', 'Test User', 100.0)
	`)
	if err != nil {
		t.Fatalf("Failed to insert test user: %v", err)
	}

	router := setupTestRouter(db)

	// Prepare request payload
	requestBody, _ := json.Marshal(map[string]interface{}{
		"player_id":      "player123",
		"transaction_id": "txn123",
		"amount":         10.0,
	})

	req, _ := http.NewRequest("POST", "/purchase", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, resp.Code)
	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Contains(t, response, "player_id")
	assert.Contains(t, response, "balance")
}

func TestWinAPI(t *testing.T) {
	db := initTestDB()
	defer db.Close()

	// Insert a test user
	_, err := db.Exec(`
		INSERT INTO players (player_id, username, password, name, balance) 
		VALUES ('player123', 'testuser', 'hashedpassword', 'Test User', 100.0)
	`)
	if err != nil {
		t.Fatalf("Failed to insert test user: %v", err)
	}

	router := setupTestRouter(db)

	// Prepare request payload
	requestBody, _ := json.Marshal(map[string]interface{}{
		"player_id":      "player123",
		"transaction_id": "txn124",
		"amount":         20.0,
	})

	req, _ := http.NewRequest("POST", "/win", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, resp.Code)
	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Contains(t, response, "player_id")
	assert.Contains(t, response, "balance")
}
