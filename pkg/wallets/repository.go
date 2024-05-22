package wallets

type WalletRepo interface {
	ActivateWallet(int64 ,string) error
	Refund(int64, int64, int64) error
	Withdraw(int64, int64, int64) error
}
