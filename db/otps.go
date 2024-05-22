package db

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"example.com/fintech-app/models"
	"github.com/dchest/uniuri"
)

func GenerateOtp(userId, walletId int64) (*models.OTPS, error) {
	// generated otp
	opt := uniuri.NewLen(6)

	newOTP := models.OTPS{
		Otp:       opt,
		UserID:    userId,
		WalletID:  walletId,
		Status:    "ACTIVE",
		CreatedAt: time.Now(),
	}

	query := `
	INSERT INTO otps (otp, userId, walletId, status, createdAt)
	VALUES (?, ?, ?, ?, ?)
	`

	stmt, err := DB.Prepare(query)

	if err != nil {
		return nil, fmt.Errorf("failed to generate otp item")
	}

	defer stmt.Close()

	_, err = stmt.Exec(newOTP.Otp, newOTP.UserID, newOTP.WalletID, newOTP.Status, newOTP.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to excute generate otp")
	}

	return &newOTP, nil
}

func ValidateOtp(otpStr string, userId, walletId int64) (bool, error) {
	query := `
	SELECT * FROM otps WHERE userId = ? AND walletId = ?
	`

	stmt, err := DB.Prepare(query)

	if err != nil {
		return false, fmt.Errorf("failed to validate otp")
	}

	defer stmt.Close()

	row := stmt.QueryRow(userId, walletId)

	var otpItem models.OTPS

	// Scan the row into the user object
	err = row.Scan(&otpItem.ID, &otpItem.Otp, &otpItem.Status, &otpItem.WalletID, &otpItem.UserID, &otpItem.CreatedAt)
	if err != nil {
		return false, fmt.Errorf("error fetching user: %v", err)
	}
	isValidOTP := otpItem.Otp == otpStr && otpItem.CreatedAt.Add(time.Minute*5).Before(time.Now())

	return isValidOTP, nil
}

func ValidateOTPAndActivateWallet(userID int64, otp string) error {
	var wg sync.WaitGroup

	// Prepare the SQL query to fetch the wallet status and OTP
	query := `
	SELECT otp, createdAt, walletId FROM otps AS o 
	INNER JOIN wallets AS w ON w.id = o.walletId
	WHERE w.status = 'REGISTERED' AND o.userId = ? AND o.createdAt <= DATE_SUB(NOW(), INTERVAL 5 MINUTE)
	`


	stmt, err := DB.Prepare(query)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("failed to prepare query: %v", err)
	}
	defer stmt.Close()

	row := stmt.QueryRow(userID)

	var newOTP models.OTPS
	err = row.Scan(&newOTP.Otp, &newOTP.CreatedAt, &newOTP.WalletID)
	if err != nil {
		fmt.Println(err)
		if err == sql.ErrNoRows {
			return fmt.Errorf("invalid OTP")
		}
		return fmt.Errorf("error fetching wallet and OTP: %v", err)
	}

	if newOTP.Otp != otp {
		return fmt.Errorf("invalid OTP")
	}


	// Update wallet status to "ACTIVE"
	wg.Add(1)
	go func() {
		defer wg.Done()
		updateQuery := `
			UPDATE wallets SET status = 'ACTIVE' WHERE id = ? AND userId =?
		`
		_, err := DB.Exec(updateQuery, newOTP.WalletID, userID)
		if err != nil {
			fmt.Println(err)
			err = fmt.Errorf("error updating wallet status: %v", err)
		}
	}()

	// Remove the OTP from the OTP table
	wg.Add(1)
	go func() {
		defer wg.Done()
		deleteQuery := `
			DELETE FROM otps WHERE walletId = ? AND userId =?
		`
		_, err := DB.Exec(deleteQuery, newOTP.WalletID, userID)
		if err != nil {
			fmt.Println(err)
			err = fmt.Errorf("error removing OTP: %v", err)
		}
	}()

	// Wait for both operations to finish
	wg.Wait()

	return nil
}

// check if the user Has an Active OTP to resend it otherwise create a new one.
func HasActiveOTP(userID int64) (string, error) {
    // Prepare the SQL query to fetch the wallet status and OTP
    query := `
        SELECT otp, createdAt, walletId FROM otps AS o 
        INNER JOIN wallets AS w ON w.id = o.walletId
        WHERE w.status = 'REGISTERED' AND o.userId = ? AND o.createdAt <= DATE_SUB(NOW(), INTERVAL 5 MINUTE)
    `

    stmt, err := DB.Prepare(query)
    if err != nil {
        return "", fmt.Errorf("failed to prepare query: %v", err)
    }
    defer stmt.Close()

    row := stmt.QueryRow(userID)

    var newOTP models.OTPS
    err = row.Scan(&newOTP.Otp, &newOTP.CreatedAt, &newOTP.WalletID)
    if err != nil {
        if err == sql.ErrNoRows {
            return "", nil
        }
        return "", fmt.Errorf("error fetching wallet and OTP: %v", err)
    }
	
    if newOTP.Otp != "" {
        return newOTP.Otp, nil
    }

    return "", nil
}

