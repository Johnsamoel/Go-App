package wallets

import "example.com/fintech-app/models"

type WalletService interface {
	CreateNewWalletService(int64) (*models.Wallet, error)
	ActivateWalletService(int64, string) error
	RefundWalletService(int64, int64, float64) error
	WithdrawWalletService(int64, int64, float64) error
}

type walletService struct {
	repo WalletRepo
}

func NewWalletService(repo WalletRepo) WalletService {
	return &walletService{repo: repo}
}

func (s *walletService) CreateNewWalletService(userId int64) (*models.Wallet, error) {
	wallet, err := s.repo.CreateNewWallet(userId)

	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func (s *walletService) ActivateWalletService(userId int64, OTP string) error {
	err := s.repo.ActivateWallet(userId, OTP)

	if err != nil {
		return err
	}

	return nil
}

func (s *walletService) RefundWalletService(userId, walletId int64, amout float64) error {
	err := s.repo.Refund(userId, walletId, amout)

	if err != nil {
		return err
	}

	return nil
}

func (s *walletService) WithdrawWalletService(userId, walletId int64, amout float64) error {
	err := s.repo.Withdraw(userId, walletId, amout)

	if err != nil {
		return err
	}

	return nil
}
