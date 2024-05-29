package main

import (
    "fmt"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "time"
)

type Customer struct {
    ID          uint
    FirstName   string
    LastName    string
    PhoneNumber string
    Email       string
}

type Account struct {
    ID           uint
    AccountNumber string
    Balance      float64
    CustomerID   uint
    CreatedAt    time.Time // Include CreatedAt field
    UpdatedAt    time.Time // Include UpdatedAt field
}

type Transaction struct {
    ID            uint
    FromAccountID uint
    ToAccountID   uint
    TransactionType string
    Amount        float64
    CreatedAt    time.Time // Include CreatedAt field
    UpdatedAt    time.Time // Include UpdatedAt field
}

func MigrateTables(db *gorm.DB) {
    // Auto-migrate the tables managed by GORM
    db.AutoMigrate(&Customer{}, &Account{}, &Transaction{})
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
        FromAccountID: account1.ID,
        ToAccountID:   0, // For deposits, the ToAccountID is 0
        TransactionType: "deposit",
        Amount:        300.00,
    }
    db.Create(&transaction1)

    transaction2 := Transaction{
        FromAccountID: account1.ID,
        ToAccountID:   0, // For deposits, the ToAccountID is 0
        TransactionType: "deposit",
        Amount:        200.00,
    }
    db.Create(&transaction2)

    transaction3 := Transaction{
        FromAccountID: account1.ID,
        ToAccountID:   0, // For deposits, the ToAccountID is 0
        TransactionType: "withdrawal",
        Amount:        50.00,
    }
    db.Create(&transaction3)

    transaction4 := Transaction{
        FromAccountID: account1.ID,
        ToAccountID:   0, // For deposits, the ToAccountID is 0
        TransactionType: "transfer",
        Amount:        100.00,
    }
    db.Create(&transaction4)
}

func GetAccountHistoryByAccountNumber(db *gorm.DB, accountNumber string) (string, error) {
    var account Account
    if err := db.Where("account_number = ?", accountNumber).First(&account).Error; err != nil {
        return "", err
    }

    var transactions []Transaction
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

func main() {
    // Connect to the PostgreSQL database
    dsn := "host=localhost user=postgres password=postgres dbname=bank port=5432 sslmode=disable"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }

    // Migrate all tables
    MigrateTables(db)

    // Seed sample data
    SeedData(db)

    // Get account history for account number "123456789"
    history, err := GetAccountHistoryByAccountNumber(db, "123456789")
    if err != nil {
        panic(err)
    }

    fmt.Println(history)
}
