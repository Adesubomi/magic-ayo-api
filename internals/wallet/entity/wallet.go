package entity

const (
	WalletAsset     = "sat"
	WalletAssetName = "satoshi"
	WalletPrecision = 8
)

type Wallet struct {
	ID      string
	UserID  string
	Asset   string
	Balance int64 // in sats
}
