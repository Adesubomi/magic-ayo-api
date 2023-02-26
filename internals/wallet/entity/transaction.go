package entity

import "time"

type TransactionCategory string

const (
	LightningFunding TransactionCategory = "lightning-funding"
	BitcoinFunding   TransactionCategory = "bitcoin-funding"
	GamePlay         TransactionCategory = "game-play"
)

type Transaction struct {
	ID           string    `json:"id" gorm:"id,primaryKey,unique,not null"`
	UserID       string    `json:"-" gorm:"user_id,not null"`
	WalletID     string    `json:"walletID" gorm:"wallet_id"`
	Category     string    `json:"category" gorm:"category"`
	Credit       int64     `json:"credit" gorm:"credit"`
	Debit        int64     `json:"debit" gorm:"debit"`
	BalanceAfter int64     `json:"balanceAfter" gorm:"balance_after,not null"`
	Status       string    `json:"status" gorm:"remark"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"createdAt"`
}
