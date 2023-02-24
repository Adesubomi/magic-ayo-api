package entity

type Wallet struct {
	ID      string
	UserId  string
	Balance int64 // in sats
}

type LNUrl struct {
	UserId      string
	PaymentHash string
	Amount      string
}

type BitcoinWallet struct {
	ID       string
	Mnemonic string
	Address  string
}

type BitcoinAddress struct {
	ID              string
	UserId          string
	BitcoinWalletId string
	WalletId        string
	Address         string
	Type            string
}
