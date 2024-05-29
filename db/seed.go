package db

import (
	"github.com/manimovassagh/bank-go/types"
	"gorm.io/gorm"
)

func SeedData(db *gorm.DB) {
	// Seed customers
	customer1 := types.Customer{
		FirstName:   "John",
		LastName:    "Doe",
		PhoneNumber: "1234567890",
		Email:       "john.doe@example.com",
	}
	db.Create(&customer1)

	// Seed accounts
	account1 := types.Account{
		AccountNumber: "123456789",
		Balance:       1000.00,
		CustomerID:    customer1.ID,
	}
	db.Create(&account1)

	// Seed transactions
	transaction1 := types.Transaction{
		FromAccountID:   account1.ID,
		ToAccountID:     0, // For deposits, the ToAccountID is 0
		TransactionType: "deposit",
		Amount:          300.00,
	}
	db.Create(&transaction1)

	transaction2 := types.Transaction{
		FromAccountID:   account1.ID,
		ToAccountID:     0, // For deposits, the ToAccountID is 0
		TransactionType: "deposit",
		Amount:          200.00,
	}
	db.Create(&transaction2)

	transaction3 := types.Transaction{
		FromAccountID:   account1.ID,
		ToAccountID:     0, // For deposits, the ToAccountID is 0
		TransactionType: "withdrawal",
		Amount:          50.00,
	}
	db.Create(&transaction3)

	transaction4 := types.Transaction{
		FromAccountID:   account1.ID,
		ToAccountID:     0, // For deposits, the ToAccountID is 0
		TransactionType: "transfer",
		Amount:          100.00,
	}
	db.Create(&transaction4)
}
