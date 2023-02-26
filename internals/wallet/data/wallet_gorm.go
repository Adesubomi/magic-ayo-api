package data

import (
	"errors"
	"github.com/Adesubomi/magic-ayo-api/internals/wallet/entity"
	logPkg "github.com/Adesubomi/magic-ayo-api/pkg/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r Repo) UserWallet(userID string) (*entity.Wallet, error) {
	wallet := &entity.Wallet{}
	walletID := userID + "-" + entity.WalletAsset

	result := r.DbClient.
		Where(&entity.Wallet{ID: walletID}).
		Attrs(&entity.Wallet{
			UserID:  userID,
			Asset:   entity.WalletAsset,
			Balance: 0.00}).
		FirstOrCreate(wallet)

	if result.Error != nil {
		return wallet, result.Error
	}

	return wallet, nil
}

func (r Repo) Credit(
	reference string,
	userID string,
	trxCategory entity.TransactionCategory,
	amount int64,
	status string) (*entity.Transaction, error) {

	wallet := &entity.Wallet{}
	creditWalletID := userID + "-" + entity.WalletAsset

	transaction := &entity.Transaction{
		ID:       reference,
		UserID:   userID,
		WalletID: creditWalletID,
		Category: string(trxCategory),
		Debit:    0.00,
		Credit:   amount,
		Status:   status,
	}

	err := r.DbClient.Transaction(func(tx *gorm.DB) error {
		var t_ []*entity.Transaction
		result := tx.Model(&entity.Transaction{}).
			Where(&entity.Transaction{ID: reference}).
			Find(&t_)

		if result.RowsAffected > 0 {
			return logPkg.DuplicateRecordError
		}

		result = tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where(entity.Wallet{ID: creditWalletID}).
			Find(wallet)

		if wallet.ID == "" {
			return logPkg.RecordNotFoundError
		}

		balanceAfter := wallet.Balance + amount
		result = tx.
			Where(entity.Wallet{
				ID:      creditWalletID,
				Balance: wallet.Balance,
			}).
			Updates(entity.Wallet{
				Balance: balanceAfter,
			})
		if result.Error != nil {
			return result.Error
		}

		transaction.BalanceAfter = balanceAfter
		result = tx.Create(transaction)
		if result.Error != nil {
			return result.Error
		}
		return nil
	})

	return transaction, err
}

func (r Repo) Debit(
	reference string,
	userID string,
	trxCategory entity.TransactionCategory,
	amount int64,
	status string) (*entity.Transaction, error) {

	wallet := &entity.Wallet{}
	debitWalletID := userID + "-" + entity.WalletAsset

	transaction := &entity.Transaction{
		ID:       reference,
		UserID:   userID,
		WalletID: debitWalletID,
		Category: string(trxCategory),
		Debit:    amount,
		Credit:   0,
		Status:   status,
	}

	err := r.DbClient.Transaction(func(tx *gorm.DB) error {
		var tN []*entity.Transaction
		result := tx.Model(&entity.Transaction{}).
			Where(&entity.Transaction{ID: reference}).
			Find(&tN)

		if result.RowsAffected > 0 {
			return logPkg.DuplicateRecordError
		}

		result = tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where(entity.Wallet{ID: debitWalletID}).
			Find(wallet)
		if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return logPkg.RecordNotFoundError
		} else if result.Error != nil {
			return result.Error
		}

		if wallet.ID == "" {
			return logPkg.RecordNotFoundError
		} else if wallet.Balance <= 0 || wallet.Balance < amount {
			return logPkg.InsufficientBalanceError
		}

		balanceAfter := wallet.Balance - amount
		result = tx.
			Where(entity.Wallet{
				ID:      debitWalletID,
				Balance: wallet.Balance,
			}).
			Updates(entity.Wallet{
				Balance: balanceAfter,
			})
		if result.Error != nil {
			return result.Error
		}

		transaction.BalanceAfter = balanceAfter
		result = tx.Create(transaction)
		if result.Error != nil {
			return result.Error
		}

		return nil
	})

	return transaction, err
}
