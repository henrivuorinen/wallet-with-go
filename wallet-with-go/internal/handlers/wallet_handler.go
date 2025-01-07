package handlers

import (
	_ "database/sql"
	_ "encoding/json"
	"github.com/google/uuid"
	_ "golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"wallet-with-go/internal/services"

	"github.com/gin-gonic/gin"
	_ "github.com/google/uuid"
	"wallet-with-go/internal/database"
)

type TransactionRequest struct {
	PlayerID      string  `json:"player_id" binding:"required"`
	TransactionID string  `json:"transaction_id" binding:"required"`
	Amount        float64 `json:"amount" binding:"required,gt=0"`
}

type TransactionResponse struct {
	PlayerID string  `json:"player_id"`
	Balance  float64 `json:"balance"`
}

type RegisterRequest struct {
	Username string  `json:"username" binding:"required"`
	Password string  `json:"password" binding:"required"`
	Name     string  `json:"name" binding:"required"`
	Balance  float64 `json:"balance" binding:"required,gt=0"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type WalletHandler struct {
	Service *services.WalletService
}

func (wh *WalletHandler) SetupRoutes(router *gin.Engine) {
	protected := router.Group("/")
	// Add CORS middleware or authentication here if needed
	{
		protected.POST("/purchase", wh.HandlePurchase)
		protected.POST("/win", wh.HandleWin)
		protected.POST("/register", wh.HandleRegister)
		protected.POST("/login", wh.HandleLogin)
	}
}

func (wh *WalletHandler) HandleRegister(c *gin.Context) {
	var input RegisterRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Check if username already exists
	exists, err := wh.Service.UsernameExists(input.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already taken"})
		return
	}

	// Generate player_id
	playerID := uuid.New().String()

	// Hash the password
	hashedPassword, err := database.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Save the player in the database
	err = wh.Service.SavePlayer(playerID, input.Username, hashedPassword, input.Name, input.Balance)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register player"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"player_id": playerID, "message": "Registration successful"})
}

func (h *WalletHandler) HandleLogin(c *gin.Context) {
	var loginRequest LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		log.Printf("Failed to bind login request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	log.Printf("Login attempt: username=%s", loginRequest.Username)

	// Use the WalletService's Login method
	playerID, balance, err := h.Service.Login(loginRequest.Username, loginRequest.Password)
	if err != nil {
		log.Printf("Login failed: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Login failed", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"player_id": playerID,
		"balance":   balance,
	})
}

func (wh *WalletHandler) HandlePurchase(c *gin.Context) {
	var req TransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	newBalance, err := wh.Service.ProcessPurchase(req.PlayerID, req.TransactionID, req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Transaction failed", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, TransactionResponse{PlayerID: req.PlayerID, Balance: newBalance})
}

func (wh *WalletHandler) HandleWin(c *gin.Context) {
	var req TransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	newBalance, err := wh.Service.ProcessWin(req.PlayerID, req.TransactionID, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction failed", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, TransactionResponse{PlayerID: req.PlayerID, Balance: newBalance})
}
