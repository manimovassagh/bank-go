package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Customer struct {
	ID          uint
	FirstName   string
	LastName    string
	PhoneNumber string
	Email       string
	Accounts    []Account `gorm:"foreignKey:CustomerID"`
}

type Account struct {
	ID            uint
	AccountNumber string
	Balance       float64
	CustomerID    uint
	Transactions  []Transaction `gorm:"foreignKey:FromAccountID"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Transaction struct {
	ID              uint
	FromAccountID   uint
	ToAccountID     uint
	TransactionType string
	Amount          float64
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type BankingService struct {
	DB *gorm.DB
}

func NewBankingService(db *gorm.DB) *BankingService {
	return &BankingService{
		DB: db,
	}
}

func (s *BankingService) Deposit(accountID uint, amount float64) error {
	var account Account
	if err := s.DB.First(&account, accountID).Error; err != nil {
		return err
	}

	account.Balance += amount

	transaction := Transaction{
		FromAccountID:   0,
		ToAccountID:     account.ID,
		TransactionType: "deposit",
		Amount:          amount,
	}

	if err := s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&account).Error; err != nil {
			return err
		}
		return tx.Create(&transaction).Error
	}); err != nil {
		return err
	}

	return nil
}

func (s *BankingService) Withdraw(accountID uint, amount float64) error {
	var account Account
	if err := s.DB.First(&account, accountID).Error; err != nil {
		return err
	}

	if account.Balance < amount {
		return errors.New("insufficient balance")
	}

	account.Balance -= amount

	transaction := Transaction{
		FromAccountID:   account.ID,
		ToAccountID:     0,
		TransactionType: "withdrawal",
		Amount:          amount,
	}

	if err := s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&account).Error; err != nil {
			return err
		}
		return tx.Create(&transaction).Error
	}); err != nil {
		return err
	}

	return nil
}

func (s *BankingService) Transfer(fromAccountID, toAccountID uint, amount float64) error {
	var fromAccount Account
	if err := s.DB.First(&fromAccount, fromAccountID).Error; err != nil {
		return err
	}

	var toAccount Account
	if err := s.DB.First(&toAccount, toAccountID).Error; err != nil {
		return err
	}

	if fromAccount.Balance < amount {
		return errors.New("insufficient balance")
	}

	fromAccount.Balance -= amount
	toAccount.Balance += amount

	transaction := Transaction{
		FromAccountID:   fromAccount.ID,
		ToAccountID:     toAccount.ID,
		TransactionType: "transfer",
		Amount:          amount,
	}

	if err := s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&fromAccount).Error; err != nil {
			return err
		}
		if err := tx.Save(&toAccount).Error; err != nil {
			return err
		}
		return tx.Create(&transaction).Error
	}); err != nil {
		return err
	}

	return nil
}

type TransactionEntry struct {
	TransactionType string  `json:"transaction_type"`
	Amount          float64 `json:"amount"`
	Date            string  `json:"date"`
	Balance         float64 `json:"balance"`
}

type AccountHistoryResponse struct {
	AccountNumber string             `json:"account_number"`
	History       []TransactionEntry `json:"history"`
}

func (s *BankingService) GetAccountHistoryByAccountNumber(accountNumber string) (*AccountHistoryResponse, error) {
	var account Account
	if err := s.DB.Where("account_number = ?", accountNumber).First(&account).Error; err != nil {
		return nil, err
	}

	var transactions []Transaction
	if err := s.DB.Where("from_account_id = ? OR to_account_id = ?", account.ID, account.ID).Find(&transactions).Error; err != nil {
		return nil, err
	}

	history := []TransactionEntry{}
	balance := account.Balance
	for _, transaction := range transactions {
		entry := TransactionEntry{
			TransactionType: transaction.TransactionType,
			Amount:          transaction.Amount,
			Date:            transaction.CreatedAt.Format("2006-01-02"),
			Balance:         balance,
		}
		switch transaction.TransactionType {
		case "deposit":
			balance += transaction.Amount
		case "withdrawal":
			balance -= transaction.Amount
		case "transfer":
			if transaction.FromAccountID == account.ID {
				balance -= transaction.Amount
			} else {
				balance += transaction.Amount
			}
		}
		history = append(history, entry)
	}

	response := &AccountHistoryResponse{
		AccountNumber: account.AccountNumber,
		History:       history,
	}

	return response, nil
}

func SeedData(db *gorm.DB) {
	// Seed customers
	customer1 := Customer{
		FirstName:   "John",
		LastName:    "Doe",
		PhoneNumber: "1234567890",
		Email:       "john.doe@example.com",
	}
	db.Create(&customer1)

	// Seed accounts
	account1 := Account{
		AccountNumber: "123456789",
		Balance:       1000.00,
		CustomerID:    customer1.ID,
	}
	db.Create(&account1)

	// Seed transactions
	transaction1 := Transaction{
		FromAccountID:   0,
		ToAccountID:     account1.ID,
		TransactionType: "deposit",
		Amount:          300.00,
	}
	db.Create(&transaction1)

	transaction2 := Transaction{
		FromAccountID:   0,
		ToAccountID:     account1.ID,
		TransactionType: "deposit",
		Amount:          200.00,
	}
	db.Create(&transaction2)

	transaction3 := Transaction{
		FromAccountID:   account1.ID,
		ToAccountID:     0,
		TransactionType: "withdrawal",
		Amount:          50.00,
	}
	db.Create(&transaction3)

	transaction4 := Transaction{
		FromAccountID:   account1.ID,
		ToAccountID:     0,
		TransactionType: "transfer",
		Amount:          100.00,
	}
	db.Create(&transaction4)
}

func main() {
	// Connect to the PostgreSQL database
	dsn := "host=localhost user=postgres password=postgres dbname=bank port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate all tables
	db.AutoMigrate(&Customer{}, &Account{}, &Transaction{})

	// Create banking service
	service := NewBankingService(db)

	// Seed sample data
	SeedData(db)

	// Initialize Echo
	e := echo.New()

	// Define route for getting account history
	e.GET("/accounts/:account_number/history", func(c echo.Context) error {
		accountNumber := c.Param("account_number")
		history, err := service.GetAccountHistoryByAccountNumber(accountNumber)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, history)
	})

	// Start the server
	e.Logger.Fatal(e.Start(":8080"))
}
