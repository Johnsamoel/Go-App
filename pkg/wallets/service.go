package wallets

import (
	"example.com/fintech-app/db"
)

type WalletService interface {
	ActivateWalletService(int64, string) error
	RefundWalletService(int64, int64, int64) error
	WithdrawWalletService(int64, int64, int64) error
}

func ActivateWalletService(userId int64, OTP string) error {
	err := db.ValidateOTPAndActivateWallet(userId, OTP)

	if err != nil {
		return err
	}

	return nil
}


func RefundWalletService(userId, walletId int64, amout float64) error {
	err := db.RefundWallet(userId, walletId ,amout)

	if err != nil {
		return err
	}

	return nil
}


func WithdrawWalletService(userId, walletId int64, amout float64) error {
	err := db.WithdrawFromWallet(userId, walletId ,amout)

	if err != nil {
		return err
	}

	return nil
}
