package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"example.com/fintech-app/db"
	"example.com/fintech-app/models"
	"example.com/fintech-app/pkg/wallets"
	"github.com/gin-gonic/gin"
)

// ActivateWalletHandler handles wallet activation
func ActivateWalletHandler(service wallets.WalletService) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqBody, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read request body"})
			return
		}

		var newOTP models.OTPS
		err = json.Unmarshal(reqBody, &newOTP)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse request body"})
			return
		}

		userId, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "You're Unauthorized"})
			return
		}

		userIdInt64, ok := userId.(int64)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		err = service.ActivateWalletService(userIdInt64, newOTP.Otp)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusAccepted, gin.H{"message": "wallet was activated successfully."})
	}
}

// RefundWalletHandler handles wallet refunds
func RefundWalletHandler(service wallets.WalletService) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqBody, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read request body"})
			return
		}

		newTransaction := models.Transactions{
			CreatedAt: time.Now(),
			Status:    "SUCCESS",
			Type:      "REFUND",
		}

		err = json.Unmarshal(reqBody, &newTransaction)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse request body"})
			return
		}

		userId, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "You're Unauthorized"})
			return
		}

		userIdInt64, ok := userId.(int64)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		walletId, exists := c.Get("walletID")
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "You're Unauthorized"})
			return
		}

		walletIdInt64, ok := walletId.(int64)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wallet ID"})
			return
		}

		err = service.RefundWalletService(userIdInt64, walletIdInt64, newTransaction.Amount)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		newTransaction.WalletID = walletIdInt64

		err = db.CreateTransaction(&newTransaction)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Wallet refunded successfully"})
	}
}

// WithdrawWalletHandler handles wallet withdrawals
func WithdrawWalletHandler(service wallets.WalletService) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqBody, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read request body"})
			return
		}

		newTransaction := models.Transactions{
			CreatedAt: time.Now(),
			Status:    "SUCCESS",
			Type:      "WITHDRAW",
		}
		err = json.Unmarshal(reqBody, &newTransaction)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse request body"})
			return
		}

		userId, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "You're Unauthorized"})
			return
		}

		userIdInt64, ok := userId.(int64)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		walletId, exists := c.Get("walletID")
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "You're Unauthorized"})
			return
		}

		walletIdInt64, ok := walletId.(int64)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wallet ID"})
			return
		}

		err = service.WithdrawWalletService(userIdInt64, walletIdInt64, newTransaction.Amount)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		newTransaction.WalletID = walletIdInt64

		err = db.CreateTransaction(&newTransaction)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Wallet balance updated successfully"})
	}
}
