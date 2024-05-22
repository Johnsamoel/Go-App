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

func ActivateWalletHandler(Context *gin.Context) {
	// Read request body
	reqBody, err := ioutil.ReadAll(Context.Request.Body)
	if err != nil {
		Context.JSON(http.StatusBadRequest, gin.H{"error": "failed to read request body"})
		return
	}

	// Parse request body into User struct
	var newOTP models.OTPS
	err = json.Unmarshal(reqBody, &newOTP)
	if err != nil {
		Context.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse request body"})
		return
	}

	// get userid from req obj
	userId, exists := Context.Get("userID")
	if !exists {
		Context.JSON(http.StatusBadRequest, gin.H{"error": "You're Unauthorized"})
		return
	}

	// Type assertion to convert interface{} to int64
	userIdInt64, ok := userId.(int64)
	if !ok {
		Context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = wallets.ActivateWalletService(userIdInt64, newOTP.Otp)

	if err != nil {
		Context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	Context.JSON(http.StatusAccepted, gin.H{"message": "wallet was activated sucessfully."})
}

func RefundWalletHandler(Context *gin.Context) {
	// Read request body
	reqBody, err := ioutil.ReadAll(Context.Request.Body)
	if err != nil {
		fmt.Println(err)
		Context.JSON(http.StatusBadRequest, gin.H{"error": "failed to read request body"})
		return
	}

	// Parse request body into User struct
	newTransaction := models.Transactions{
		CreatedAt: time.Now(),
		Status:    "SUCCESS",
		Type:      "REFUND",
	}

	err = json.Unmarshal(reqBody, &newTransaction)
	if err != nil {
		fmt.Println(err)
		Context.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse request body"})
		return
	}

	// get userid from req obj
	userId, exists := Context.Get("userID")
	if !exists {
		Context.JSON(http.StatusBadRequest, gin.H{"error": "You're Unauthorized"})
		return
	}

	// Type assertion to convert interface{} to int64
	userIdInt64, ok := userId.(int64)
	if !ok {
		Context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// get userid from req obj
	walletId, exists := Context.Get("walletID")
	if !exists {
		Context.JSON(http.StatusBadRequest, gin.H{"error": "You're Unauthorized"})
		return
	}

	// Type assertion to convert interface{} to int64
	walletIdInt64, ok := walletId.(int64)
	if !ok {
		Context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Call the RefundWalletService function with walletID and RefundAmount
	err = wallets.RefundWalletService(userIdInt64, walletIdInt64, newTransaction.Amount)
	if err != nil {
		Context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newTransaction.WalletID = walletIdInt64

	err = db.CreateTransaction(&newTransaction)

	if err != nil {
		Context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	Context.JSON(http.StatusOK, gin.H{"message": "Wallet refunded successfully"})
}

func WithdrawWalletHandler(Context *gin.Context) {
	// Read request body
	reqBody, err := ioutil.ReadAll(Context.Request.Body)
	if err != nil {
		Context.JSON(http.StatusBadRequest, gin.H{"error": "failed to read request body"})
		return
	}

	// Parse request body into User struct
	newTransaction := models.Transactions{
		CreatedAt: time.Now(),
		Status:    "SUCCESS",
		Type:      "WITHDRAW",
	}
	err = json.Unmarshal(reqBody, &newTransaction)
	if err != nil {
		Context.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse request body"})
		return
	}

	// get userid from req obj
	userId, exists := Context.Get("userID")
	if !exists {
		Context.JSON(http.StatusBadRequest, gin.H{"error": "You're Unauthorized"})
		return
	}

	// Type assertion to convert interface{} to int64
	userIdInt64, ok := userId.(int64)
	if !ok {
		Context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// get userid from req obj
	walletId, exists := Context.Get("walletID")
	if !exists {
		Context.JSON(http.StatusBadRequest, gin.H{"error": "You're Unauthorized"})
		return
	}

	// Type assertion to convert interface{} to int64
	walletIdInt64, ok := walletId.(int64)
	if !ok {
		Context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Call the withdraw function with walletID and RefundAmount
	err = wallets.WithdrawWalletService(userIdInt64, walletIdInt64, newTransaction.Amount)
	if err != nil {
		Context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newTransaction.WalletID = walletIdInt64

	err = db.CreateTransaction(&newTransaction)

	if err != nil {
		Context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	Context.JSON(http.StatusOK, gin.H{"message": "Wallet balace updated successfully"})
}
