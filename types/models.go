package types

import "time"

type Customer struct {
	ID          uint
	FirstName   string
	LastName    string
	PhoneNumber string
	Email       string
}

type Account struct {
	ID            uint
	AccountNumber string
	Balance       float64
	CustomerID    uint
	CreatedAt     time.Time // Include CreatedAt field
	UpdatedAt     time.Time // Include UpdatedAt field
}

type Transaction struct {
	ID              uint
	FromAccountID   uint
	ToAccountID     uint
	TransactionType string
	Amount          float64
	CreatedAt       time.Time // Include CreatedAt field
	UpdatedAt       time.Time // Include UpdatedAt field
}
