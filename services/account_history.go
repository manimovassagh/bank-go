package services

import (
	"fmt"
	"github.com/manimovassagh/bank-go/types"
	"gorm.io/gorm"
)

func GetAccountHistoryByAccountNumber(db *gorm.DB, accountNumber string) (string, error) {
	var account types.Account
	if err := db.Where("account_number = ?", accountNumber).First(&account).Error; err != nil {
		return "", err
	}

	var transactions []types.Transaction
	if err := db.Where("from_account_id = ? OR to_account_id = ?", account.ID, account.ID).Find(&transactions).Error; err != nil {
		return "", err
	}

	history := fmt.Sprintf("Account Number: %s\n", account.AccountNumber)
	balance := account.Balance
	for _, transaction := range transactions {
		switch transaction.TransactionType {
		case "deposit":
			balance += transaction.Amount
			history += fmt.Sprintf("Deposit %.2f on %s, Balance: %.2f\n", transaction.Amount, transaction.CreatedAt.Format("2006-01-02"), balance)
		case "withdrawal":
			balance -= transaction.Amount
			history += fmt.Sprintf("Withdrawal %.2f on %s, Balance: %.2f\n", transaction.Amount, transaction.CreatedAt.Format("2006-01-02"), balance)
		case "transfer":
			if transaction.FromAccountID == account.ID {
				balance -= transaction.Amount
				history += fmt.Sprintf("Transfer out %.2f on %s, Balance: %.2f\n", transaction.Amount, transaction.CreatedAt.Format("2006-01-02"), balance)
			} else {
				balance += transaction.Amount
				history += fmt.Sprintf("Transfer in %.2f on %s, Balance: %.2f\n", transaction.Amount, transaction.CreatedAt.Format("2006-01-02"), balance)
			}
		}
	}

	return history, nil
}
