package wallets

import (
	"database/sql"

	dataSource "example.com/fintech-app/db"
	"example.com/fintech-app/models"
)

type WalletRepo interface {
	CreateNewWallet(int64) (*models.Wallet, error)
	ActivateWallet(int64, string) error
	Refund(int64, int64, float64) error
	Withdraw(int64, int64, float64) error
}

type walletRepo struct {
	db      *sql.DB
	WalletMethods dataSource.WalletRepository
	OTPMehtods dataSource.OTPRepository
}

// NewUserRepo initializes and returns an instance of userRepo implementing UserRepo
func NewWalletRepo(db *sql.DB) WalletRepo {
	return &walletRepo{db: db}
}

func (w *walletRepo) CreateNewWallet(userID int64) (*models.Wallet, error) {
	wallet , err := w.WalletMethods.CreateNewWallet(userID)

	return wallet , err
}

func (w *walletRepo) ActivateWallet(userID int64, OTP string) error {
	err := w.OTPMehtods.ValidateOTPAndActivateWallet(userID, OTP)
	return err
}

func (w *walletRepo) Refund(userId, walletId int64, amout float64) error {
	err := w.WalletMethods.RefundWallet(userId, walletId, amout)
	return err
}

func (w *walletRepo) Withdraw(userId, walletId int64, amout float64) error {
	err := w.WalletMethods.WithdrawFromWallet(userId, walletId, amout)
	return err
}
