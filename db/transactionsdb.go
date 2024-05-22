package db

import (
	"fmt"

	"example.com/fintech-app/models"
)

func CreateTransaction(transaction *models.Transactions) error {
	query := `INSERT INTO transactions (status, type, amount, walletId, createdAt)
	VALUES (?, ?, ?, ?, ?)
	`
	stmt, err := DB.Prepare(query)

	if err != nil {
		return fmt.Errorf("failed to generate transaction")
	}

	defer stmt.Close()

	_, err = stmt.Exec(transaction.Status, transaction.Type, transaction.Amount, transaction.WalletID, transaction.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to excute generate transaction")
	}

	return nil

}
