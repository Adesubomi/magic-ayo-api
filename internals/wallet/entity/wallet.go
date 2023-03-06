package entity

const (
	WalletAsset     = "sat"
	WalletAssetName = "satoshi"
	WalletPrecision = 8
)

type Wallet struct {
	ID      string `json:"id"`
	UserID  string `json:"-"`
	Asset   string `json:"asset"`
	Balance int64  `json:"balance"` // in sats
}
