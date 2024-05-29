package main

import (
	"fmt"
	db2 "github.com/manimovassagh/bank-go/db"
	"github.com/manimovassagh/bank-go/services"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Connect to the PostgreSQL database
	dsn := "host=localhost user=postgres password=postgres dbname=bank port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate all tables
	db2.MigrateTables(db)

	// Seed sample data
	db2.SeedData(db)

	// Get account history for account number "123456789"
	history, err := services.GetAccountHistoryByAccountNumber(db, "123456789")
	if err != nil {
		panic(err)
	}

	fmt.Println(history)
}
