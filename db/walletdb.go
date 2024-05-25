package db

import (
	"database/sql"
	"fmt"

	model "example.com/fintech-app/models"
)

type WalletRepository struct{}

var WalletRepo = WalletRepository{}


func (w WalletRepository) GetWalletByUserId(userId int64) (*model.Wallet, error) {
	// get the associated user wallet data..
	query := `SELECT id, balance, status, userId FROM wallets WHERE userId = ?`
	stmt, err := DB.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare query: %v", err)
	}
	defer stmt.Close()
	walletRow := stmt.QueryRow(userId)

	// Initialize a new wallet object to view the fetched wallet status data
	var userWallet model.Wallet

	// Scan the row into the wallet object
	err = walletRow.Scan(&userWallet.ID, &userWallet.Balance, &userWallet.Status, &userWallet.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			// No wallet found with the given user ID
			return nil, sql.ErrNoRows
		}
		// Other error occurred
		return nil, fmt.Errorf("something went wrong")
	}

	return &userWallet, nil
}

func (w WalletRepository) CreateNewWallet(userID int64) (*model.Wallet, error) {
	// Prepare the SQL query
	query := `
        INSERT INTO wallets (balance, userId, status)
        VALUES (?, ?, ?)
    `

	// Execute the SQL query with user data
	result, err := DB.Exec(query, 0, userID, "REGISTERED")
	if err != nil {
		return nil, fmt.Errorf("error creating new wallet: %v", err)
	}

	// Get the ID of the newly inserted wallet
	walletID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert ID: %v", err)
	}

	// Create a new Wallet object with the ID
	newWallet := &model.Wallet{
		ID:      walletID,
		UserID:  userID,
		Balance: 0,
		Status:  "REGISTERED",
	}

	return newWallet, nil
}

func (w WalletRepository) DeleteWallet(walletId int64) error {
	_, err := DB.Exec(`DELETE FROM users WHERE id = ?`, walletId)

	if err != nil {
		return fmt.Errorf("error deleting user wallet : %v", err)
	}

	return nil
}

func (w WalletRepository) RefundWallet(userID, walledId int64, amount float64) error {
	// Prepare the SQL query to update the wallet balance
	query := `
        UPDATE wallets 
        SET balance = balance + ?
        WHERE id= ? AND userId = ? AND status = 'ACTIVE'
    `

	// Start a transaction
	tx, err := DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // Re-throw panic after Rollback
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// Execute the SQL query within the transaction
	_, err = tx.Exec(query, amount, walledId, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("wallet not found or inactive")
		}
		return fmt.Errorf("error updating wallet balance: %v", err)
	}

	return nil
}

// Withdraw deducts the specified amount from the wallet balance if conditions are met.
func (w WalletRepository) WithdrawFromWallet(userID, walletId int64, amount float64) error {
	// Start a transaction
	tx, err := DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	// Check if the wallet is active and get its current balance
	var currentBalance float64
	err = tx.QueryRow("SELECT balance FROM wallets WHERE id= ? AND userId = ? AND status = 'ACTIVE'", walletId, userID).Scan(&currentBalance)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("wallet not found or not active")
		}
		return fmt.Errorf("failed to fetch wallet balance: %v", err)
	}

	// Check if the wallet has sufficient balance
	if currentBalance < amount {
		return fmt.Errorf("insufficient balance")
	}

	// Deduct the amount from the wallet balance
	_, err = tx.Exec("UPDATE wallets SET balance = balance - ? WHERE id = ? AND userId = ?", amount, walletId, userID)
	if err != nil {
		return fmt.Errorf("failed to deduct amount from wallet: %v", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}
